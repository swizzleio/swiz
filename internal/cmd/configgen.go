package cmd

import "github.com/urfave/cli/v2"

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "generate",
		Usage:  "Generate the app config and stack config",
		Action: configGenCmd,
		Flags:  []cli.Flag{},
	})
}

func configGenCmd(ctx *cli.Context) error {
	return nil
}
