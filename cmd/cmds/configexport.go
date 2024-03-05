package cmds

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
	if !appConfigMgr.IsLoaded() {
		return fmt.Errorf("missing app config")
	}

	out, err := appConfigMgr.GetBase64()
	if err != nil {
		return err
	}

	cl.Info("This output can be used to share the app config with developers. The signature and word list are\n")
	cl.Info("used to verify the integrity of the app config and crytographically the same.\n\n")
	cl.Info("App config: %v\n", out.Encoded)
	cl.Info("Signature: %v\n", out.Signature)
	cl.Info("Word list: %v\n", out.WordList)

	return nil
}
