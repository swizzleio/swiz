package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "version",
		Usage:  "app version number",
		Action: versionCmd,
	})
}

func versionCmd(ctx *cli.Context) error {
	fmt.Printf("Version is %v(%v)\n", Version, CommitHash)

	return nil
}
