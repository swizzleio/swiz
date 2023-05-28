package cmd

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/environment"
	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "delete",
		Usage:  "Delete an environment",
		Action: envDeleteCmd,
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
				Name:  "no-orphan-delete",
				Usage: "Do not attempt to delete orphan stacks. Can be overriden in config",
			},
			&cli.BoolFlag{
				Name:  "fast-delete",
				Usage: "Delete everything in parallel. Can be overridden in config",
			},
		},
	})
}

func envDeleteCmd(ctx *cli.Context) error {
	enclave := ctx.String("enclave")
	envDef := ctx.String("env-def")
	envName := ctx.String("name")
	dryRun := ctx.Bool("dry-run")
	noOrphanDelete := ctx.Bool("no-orphan-delete")
	fastDelete := ctx.Bool("fast-delete")

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	stackInfo, err := svc.DeleteEnvironment(enclave, envDef, envName, dryRun, noOrphanDelete, fastDelete)
	if err != nil {
		return err
	}

	for _, stack := range stackInfo {
		fmt.Printf("Stack: %v [%v] - %v\n", stack.Name, stack.DeployStatus.State, stack.NextAction)
	}

	return nil
}
