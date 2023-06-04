package cmd

import "github.com/urfave/cli/v2"

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "export",
		Usage:  "Export the app config",
		Action: configExportCmd,
		Flags:  []cli.Flag{},
	})
}

func configExportCmd(ctx *cli.Context) error {
	return nil
}
