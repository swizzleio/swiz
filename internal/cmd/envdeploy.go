package cmd

import (
	"github.com/swizzleio/swiz/internal/environment"

	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "deploy",
		Usage:  "Deploy an environment. If the env already exists, it will be updated",
		Action: envCreateCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the environment",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "env-def",
				Aliases:     []string{"d"},
				Usage:       "Environment definition to use",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:        "enclave",
				Aliases:     []string{"e"},
				Usage:       "Enclave to use",
				DefaultText: "",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "If this is a dry run (also similar to plan)",
			},
			&cli.BoolFlag{
				Name:  "no-update-deploy",
				Usage: "Fail if a stack or environment already exists. Can be overriden in config",
			},
		},
	})
}

func envCreateCmd(ctx *cli.Context) error {
	enclave := ctx.String("enclave")
	envDef := ctx.String("env-def")
	envName := ctx.String("name")
	dryRun := ctx.Bool("dry-run")
	noUpdate := ctx.Bool("no-update-deploy")

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	return svc.DeployEnvironment(enclave, envDef, envName, dryRun, noUpdate)
}
