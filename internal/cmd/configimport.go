package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/pkg/security"
	"github.com/urfave/cli/v2"
	"strings"
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
	qs := []*survey.Question{
		{
			Name:      "appConfig",
			Prompt:    &survey.Input{Message: "Paste the app config here"},
			Transform: survey.TransformString(strings.TrimSpace),
		},
		{
			Name:      "signature",
			Prompt:    &survey.Input{Message: "Paste the signature here"},
			Transform: survey.TransformString(strings.TrimSpace),
		},
	}

	answers := struct {
		AppConfig string
		Signature string
	}{}

	const sigMatchErr = "signature does not match, let your security team know"

	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	if answers.AppConfig == "" {
		return fmt.Errorf("missing app config")
	}

	sig, wordList := security.GetSha256AndWordList(answers.AppConfig)
	matches := sig == answers.Signature
	if answers.Signature == "" {
		// If the signature was not provided, ask the user
		fmt.Printf("\nVerify one or the other matches the signature provided with the app config\n")
		fmt.Printf("Signature: %v\n", sig)
		fmt.Printf("Word list: %v\n", wordList)

		prompt := &survey.Confirm{
			Message: "Does one or the other match?",
		}

		err = survey.AskOne(prompt, &matches)
		if err != nil {
			return err
		}
	}

	if !matches {
		return fmt.Errorf(sigMatchErr)
	}

	return appconfig.Fetch(answers.AppConfig)
}
