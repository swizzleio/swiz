package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "delete",
		Usage:  "TBD",
		Action: envDeleteCmd,
	})
}

func envDeleteCmd(ctx *cli.Context) error {
	fmt.Printf("To be implemented\n")

	return nil
}
