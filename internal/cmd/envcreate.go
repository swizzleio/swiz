package cmd

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/fileutil"
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
	stack, err := fileutil.YamlFromLocation[model.StackConfig]("file://test/data/cloudformation/sleepstack-cfg.yaml")
	if err != nil {
		return err
	}

	fmt.Printf("Stack: %v\n", stack)

	return nil
}
