package cmds

import "github.com/urfave/cli/v2"

func init() {
	addCommand(&cli.Command{
		Name:   "config",
		Usage:  "configure swiz",
		Before: genCommandPreflight(true),
	})
}
