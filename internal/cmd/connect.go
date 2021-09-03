package cmd

import (
	"context"
	"fmt"
	"getswizzle.io/swiz/pkg/clihelper"
	"getswizzle.io/swiz/pkg/infra"
	"getswizzle.io/swiz/pkg/network/ssh"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	addCommand(&cli.Command{
		Name:   "connect",
		Usage:  "Connect to a cloud resource",
		Action: connectCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bastion",
				Aliases:  []string{"b"},
				Usage:    "bastion host address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Aliases:  []string{"k"},
				Usage:    "ssh key file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "connect",
				Aliases:  []string{"c"},
				Usage:    "remote endpoint to tunnel to",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "service",
				Aliases:  []string{"s"},
				Usage:    "service to use",
				Required: false,
			},
		},
	})
}

func launchTunnel(tun *ssh.Tunnel) {
	cCtx, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})
	go func(ctx context.Context) {
		err := tun.Start()
		if err != nil {
			log.Fatalf("starting tunnel %v", err)
		}

		select {
		case <-ctx.Done():
			log.Printf("exiting...")
			tun.Close()
			exitCh <- struct{}{}
		default: // to make this non blocking
		}
	}(cCtx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)

	go func() {
		select {
		case <-sigCh:
			cancel()
			return
		}
	}()

	time.Sleep(500 * time.Millisecond) // TODO: Poll for channel

	<-exitCh
}

// connectCmd runs the connect command
func connectCmd(ctx *cli.Context) error {
	fmt.Printf("Connecting to host\n")
	bastion := ctx.String("bastion")
	key := ctx.String("key")
	host := ctx.String("connect")

	// Create the infra services and determine which service to use
	svc, err := infra.NewInfraService()
	if err != nil {
		log.Fatalf("creating infrastructure service. %v", err)
	}

	services := svc.ListServices()

	service := clihelper.GetOrPromptOptions(ctx, "service", "Select the service that you would like to connect to",
		services, "Quit")
	if service == "" {
		// Quit
		return nil
	}
	log.Printf("Fetching instances from %v", service)

	hosts, err := svc.GetInstances(service)
	if err != nil {
		log.Fatalf("error fetching instances: %v", err)
	}
	hostMap := map[string]string{}
	for k, v := range hosts {
		hostMap[v.String()] = k
	}
	host = clihelper.GetOrPromptOptions(ctx, "connect", "Select the host that you want to connect to", hostMap, "Quit")
	if host == "" {
		// Quit
		return nil
	}
	log.Printf("%v", host)

	// Connect
	keyAuth := ssh.NewPrivateKeyAuth()
	err = keyAuth.InitFromFile(key)
	if err != nil {
		log.Fatalf("loading key %v", err)
	}
	tun := ssh.NewSshTunnel(bastion, keyAuth.GetAuthMethod(), host, 0)

	launchTunnel(tun)

	return nil
}
