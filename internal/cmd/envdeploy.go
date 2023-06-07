package cmd

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/environment"
	"strings"

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
			&cli.StringSliceFlag{
				Name:    "stack",
				Aliases: []string{"s"},
				Usage:   "Stacks to deploy. Can be specified multiple times or be a comma seperated list",
			},
			&cli.BoolFlag{
				Name:  "deploy-all",
				Usage: "Deploy all stacks",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "If this is a dry run (also similar to plan)",
			},
			&cli.BoolFlag{
				Name:  "no-update-deploy",
				Usage: "Fail if a stack or environment already exists. Can be overridden in config",
			},
		},
	})
}

func envCreateCmd(ctx *cli.Context) error {
	enclave := ctx.String("enclave")
	envDef := ctx.String("env-def")
	envName := ctx.String("name")
	stacks := ctx.StringSlice("stack")
	deployAll := ctx.Bool("deploy-all")
	dryRun := ctx.Bool("dry-run")
	noUpdate := ctx.Bool("no-update-deploy")

	stackList := []string{}
	for _, stack := range stacks {
		commaSeperated := strings.Split(stack, ",")
		stackList = append(stackList, commaSeperated...)
	}

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	stackInfo, err := svc.DeployEnvironment(ctx.Context, enclave, envDef, envName, deployAll, stackList, dryRun, noUpdate)
	if err != nil {
		return err
	}

	for _, stack := range stackInfo {
		fmt.Printf("Stack: %v [%v] - %v\n", stack.Name, stack.DeployStatus.State, stack.NextAction)
	}

	return nil
}
