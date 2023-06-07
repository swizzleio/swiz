package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(&cli.Command{
		Name:   "env",
		Usage:  "perform actions against environments",
		Before: genCommandPreflight(false),
	})
}
