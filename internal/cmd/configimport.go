package cmd

import (
	"fmt"
	appcli "github.com/swizzleio/swiz/pkg/cli"
	"github.com/swizzleio/swiz/pkg/security"
	"github.com/urfave/cli/v2"
)

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "import",
		Usage:  "Import the app config",
		Action: configImportCmd,
		Flags:  []cli.Flag{},
	})
}

func configImportCmd(ctx *cli.Context) error {
	const sigMatchErr = "signature does not match, let your security team know"

	prompts := []appcli.AskManyOpts{
		{
			Key:           "AppConfig",
			Message:       "Paste the app config here",
			Required:      true,
			TransformMode: appcli.TransformModeTrimSpace,
		},
		{
			Key:           "Signature",
			Message:       "Paste the signature here",
			TransformMode: appcli.TransformModeTrimSpace,
		},
	}

	answers, err := cl.AskMany(prompts)
	if err != nil {
		return err
	}

	sig, wordList := security.GetSha256AndWordList(answers["AppConfig"])
	matches := sig == answers["Signature"]
	if answers["Signature"] == "" {
		// If the signature was not provided, ask the user
		cl.Info("\nVerify one or the other matches the signature provided with the app config\n")
		cl.Info("Signature: %v\n", sig)
		cl.Info("Word list: %v\n", wordList)

		matches, err = cl.AskConfirm("Does one or the other match?")
		if err != nil {
			return err
		}
	}

	if !matches {
		return fmt.Errorf(sigMatchErr)
	}

	return appConfigMgr.GenFromB64(answers["AppConfig"], true)
}
