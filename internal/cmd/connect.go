package cmd

import (
	"context"
	"fmt"
	"getswizzle.io/swiz/internal/config"
	"getswizzle.io/swiz/pkg/clihelper"
	"getswizzle.io/swiz/pkg/infra"
	"getswizzle.io/swiz/pkg/infra/model"
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
				Name:    "cfg",
				Aliases: []string{"c"},
				Usage:   "config file location",
			},
			&cli.StringFlag{
				Name:     "remote",
				Aliases:  []string{"r"},
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

func launchTunnel(tun *ssh.Tunnel) chan struct{} {
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

	return exitCh
}

// connectCmd runs the connect command
func connectCmd(ctx *cli.Context) error {
	fmt.Printf("Connecting to host\n")
	host := ctx.String("remote")
	filename := ctx.String("cfg")

	// Load the config
	cfgStore := config.NewConfigMustLoad(filename)

	// Create the infra services and determine which service to use
	svc, err := infra.NewInfraService()
	if err != nil {
		log.Fatalf("creating infrastructure service. %v", err)
	}

	var hostObj *model.TargetInstance
	if host == "" {
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
		hostObjMap := map[string]*model.TargetInstance{}
		for k, v := range hosts {
			hostMap[v.String()] = k
			h := hosts[k]
			hostObjMap[k] = &h
		}
		host = clihelper.GetOrPromptOptions(ctx, "remote", "Select the host that you want to connect to", hostMap, "Quit")
		if host == "" {
			// Quit
			return nil
		}

		hostObj = hostObjMap[host]
	}

	launchInfo, err := cfgStore.GetHostLaunchInfo(*hostObj)
	if err != nil {
		log.Fatalf("fetching host launch info %v", err)
	}

	// Connect
	keyAuth := ssh.NewPrivateKeyAuth()
	err = keyAuth.InitFromFile(launchInfo.BastionAuth.KeyFilename) // TODO: BastionAuth should come from a keyauth factory
	if err != nil {
		log.Fatalf("loading key %v", err)
	}
	tun := ssh.NewSshTunnel(launchInfo.BastionAddr, launchInfo.BastionSignature, keyAuth.GetAuthMethod(), launchInfo.HostString, 0)

	exitCh := launchTunnel(tun)

	/*
		// Launch app.
		clientSvc := client.NewService()
		err = clientSvc.Launch(launchInfo.Os, launchInfo.ClientConfig)
		if err != nil {
			log.Fatalf("launching client app %v", err)
		}
	*/
	// Wait for exit. TODO: Clean up these channels once end to end is working
	<-exitCh

	return nil
}
