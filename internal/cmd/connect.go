package cmd

import (
	"fmt"
	"getswizzle.io/swiz/pkg/infra/aws"
	"getswizzle.io/swiz/pkg/network/ssh"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
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
				Name:     "host",
				Aliases:  []string{"h"},
				Usage:    "remote host to tunnel to",
				Required: true,
			},
		},
	})
}

// connectCmd runs the connect command
func connectCmd(ctx *cli.Context) error {
	fmt.Printf("Connecting to host\n")
	bastion := ctx.String("bastion")
	key := ctx.String("key")
	host := ctx.String("host")

	// Dump info
	aws.InitService()

	// Connect
	keyAuth := ssh.NewPrivateKeyAuth()
	err := keyAuth.InitFromFile(key)
	if err != nil {
		log.Fatalf("loading key %v", err)
	}
	tun := ssh.NewSshTunnel(bastion, keyAuth.GetAuthMethod(), host, 0)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		err = tun.Start()
		if err != nil {
			log.Fatalf("starting tunnel %v", err)
		}

		// Run Cleanup
		tun.Close()
		os.Exit(1)
	}()

	return nil
}
