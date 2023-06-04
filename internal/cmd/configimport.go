package cmd

import "github.com/urfave/cli/v2"

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "import",
		Usage:  "Import the app config",
		Action: configImportCmd,
		Flags:  []cli.Flag{},
	})
}

func configImportCmd(ctx *cli.Context) error {
	return nil
}
