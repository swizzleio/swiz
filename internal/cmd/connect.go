package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "connect",
		Usage:  "Connect to a cloud resource",
		Action: connectCmd,
	})
}

// connectCmd runs the connect command
func connectCmd(*cli.Context) error {
	fmt.Printf("Placeholder for connect command\n")
	return nil
}
