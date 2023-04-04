package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "envinfo",
		Usage:  "TBD",
		Action: envInfoCmd,
	})
}

func envInfoCmd(ctx *cli.Context) error {
	fmt.Printf("To be implemented\n")

	return nil
}
