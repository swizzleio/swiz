package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "envcreate",
		Usage:  "TBD",
		Action: envCreateCmd,
	})
}

func envCreateCmd(ctx *cli.Context) error {
	fmt.Printf("To be implemented\n")

	return nil
}
