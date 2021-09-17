package cmd

import (
	"getswizzle.io/swiz/internal/config"
	"getswizzle.io/swiz/pkg/common"
	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
)

func init() {
	addCommand(&cli.Command{
		Name:   "config",
		Usage:  "Configure application",
		Action: configCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cfg",
				Aliases: []string{"c"},
				Usage:   "config file location",
			},
		},
	})
}

// configCmd runs the create config command
func configCmd(ctx *cli.Context) error {
	log.Printf("Configure application\n")
	filename := ctx.String("cfg")

	// Load the config
	cfgStore := config.NewConfig()

	cfg, err := cfgStore.Load(filename)
	if err != nil {
		log.Println("Config file does not exist, new file will be created")
	}

	// Ask the user about bastion hosts
	bastionHosts := strings.Join(cfg.BastionAddrs, ", ")
	bastionPem := cfg.BastionAuth.KeyFilename
	questions := []*survey.Question{
		{
			Name: "bastionhosts",
			Prompt: &survey.Input{
				Message: "Provide a comma seperated list of bastion hosts in the form of username@ip: ",
				Default: bastionHosts,
			},
			Validate: survey.Required,
		},
		{
			Name: "bastionhostkey",
			Prompt: &survey.Input{
				Message: "Provide the pem filename used for the bastion host: ",
				Default: bastionPem,
			},
			Validate: survey.Required,
		},
		{
			Name: "linuxhostkey",
			Prompt: &survey.Input{
				Message: "Provide the pem filename used for linux hosts: ",
				Default: bastionPem,
			},
		},
		{
			Name: "linuxflavor",
			Prompt: &survey.Select{
				Message: "Choose the flavor of linux you use: ",
				Options: []string{common.LinuxAmazon, common.LinuxCentos, common.LinuxUbuntu, common.LinuxOther},
				Default: "amazon",
			},
		},
		{
			Name: "write",
			Prompt: &survey.Select{
				Message: "Are you sure you want to overwrite the file?",
				Options: []string{"yes", "no"},
				Default: "no",
			},
		},
	}

	answers := struct {
		BastionHosts   string
		BastionHostKey string
		LinuxHostKey   string
		LinuxFlavor    string
		Write          string
	}{}
	err = survey.Ask(questions, &answers)
	if err != nil {
		log.Fatalf("prompting user %v", err)
	}

	if answers.Write == "yes" {
		cfgStore.GenerateDefaults(answers.LinuxFlavor, answers.LinuxHostKey, answers.BastionHosts, answers.BastionHostKey)
		err = cfgStore.Save(filename)
		if err != nil {
			log.Fatalf("saving config file %v", err)
		}
	}

	return nil
}
