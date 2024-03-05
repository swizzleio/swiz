package cmds

import (
	"github.com/swizzleio/swiz/internal/environment"
	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "info",
		Usage:  "Fetch info on an environment",
		Action: envInfoCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the environment",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "enclave",
				Aliases:     []string{"e"},
				Usage:       "Enclave to use",
				DefaultText: "",
			},
		},
	})
}

func envInfoCmd(ctx *cli.Context) error {
	enclave := ctx.String("enclave")
	envDef := ctx.String("env-def")
	envName := ctx.String("name")

	svc, err := environment.NewEnvService(appConfigMgr.Get())
	if err != nil {
		return err
	}

	envInfo, err := svc.GetEnvironmentInfo(ctx.Context, enclave, envDef, envName)
	if err != nil {
		return err
	}

	cl.Info("Name: %v\n", envInfo.EnvironmentName)
	cl.Info("Status: %v\n", envInfo.DeployStatus)
	cl.Info("Stacks [Status]:\n")
	for _, stack := range envInfo.StackInfo {
		cl.Info("  %v [%v]\n", stack.Name, stack.DeployStatus.State)
	}

	return nil
}
