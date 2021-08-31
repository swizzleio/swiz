package cmd

import (
	"context"
	"fmt"
	"getswizzle.io/swiz/pkg/infra/aws"
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
				Required: true,
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

	// Dump info
	aws.InitService()

	// Connect
	keyAuth := ssh.NewPrivateKeyAuth()
	err := keyAuth.InitFromFile(key)
	if err != nil {
		log.Fatalf("loading key %v", err)
	}
	tun := ssh.NewSshTunnel(bastion, keyAuth.GetAuthMethod(), host, 0)

	launchTunnel(tun)

	return nil
}
