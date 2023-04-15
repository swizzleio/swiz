package cmd

import (
	"github.com/swizzleio/swiz/internal/environment"

	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "create",
		Usage:  "TBD",
		Action: envCreateCmd,
	})
}

func envCreateCmd(ctx *cli.Context) error {
	envName := ctx.String("name")

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	return svc.CreateEnvironment("dev", envName)
}
