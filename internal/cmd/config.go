package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "config",
		Usage:  "Configure application",
		Action: configCmd,
	})
}

// configCmd runs the create config command
func configCmd(*cli.Context) error {
	fmt.Printf("Placeholder for config command\n")
	return nil
}
