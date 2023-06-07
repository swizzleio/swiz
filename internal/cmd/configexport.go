package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "export",
		Usage:  "Export the app config",
		Action: configExportCmd,
		Flags:  []cli.Flag{},
	})
}

func configExportCmd(ctx *cli.Context) error {
	if appConfig == nil {
		return fmt.Errorf("missing app config")
	}

	out, err := appConfig.GetBase64()
	if err != nil {
		return err
	}

	fmt.Printf("This output can be used to share the app config with developers. The signature and word list are\n")
	fmt.Printf("used to verify the integrity of the app config and crytographically the same.\n\n")
	fmt.Printf("App config: %v\n", out.Encoded)
	fmt.Printf("Signature: %v\n", out.Signature)
	fmt.Printf("Word list: %v\n", out.WordList)

	return nil
}
