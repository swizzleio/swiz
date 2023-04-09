package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "info",
		Usage:  "TBD",
		Action: envInfoCmd,
	})
}

func envInfoCmd(ctx *cli.Context) error {
	fmt.Printf("To be implemented\n")

	return nil
}
