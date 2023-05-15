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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the environment",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "envDef",
				Aliases:     []string{"d"},
				Usage:       "Environment definition to use",
				DefaultText: "",
			},
		},
	})
}

func envCreateCmd(ctx *cli.Context) error {
	envDef := ctx.String("envDef")
	envName := ctx.String("name")

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	return svc.CreateEnvironment("dev", envDef, envName)
}
