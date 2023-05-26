package cmd

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/environment"
	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("env", &cli.Command{
		Name:   "list",
		Usage:  "List all environments",
		Action: envListCmd,
		Flags:  []cli.Flag{
			&cli.StringFlag{
				Name:        "enclave",
				Aliases:     []string{"e"},
				Usage:       "Enclave to use",
				DefaultText: "",
			},
		},
	})
}

func envListCmd(ctx *cli.Context) error {
	enclave := ctx.String("enclave")

	svc, err := environment.NewEnvService(appConfig)
	if err != nil {
		return err
	}

	envList, err := svc.ListEnvironments(enclave)
	if err != nil {
		return err
	}

	for _, env := range envList {
		fmt.Printf("%v\n", env)
	}

	return nil
}
