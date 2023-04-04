package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "envlist",
		Usage:  "TBD",
		Action: envListCmd,
	})
}

func envListCmd(ctx *cli.Context) error {
	fmt.Printf("To be implemented\n")

	return nil
}
