package cmd

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	Version    = "dev"
	CommitHash = "n/a"
)

// App config
var appConfig *appconfig.AppConfig

// Internal list of commands
var commands = []*cli.Command{}

// CLI app
var appCli *cli.App

// addCommand adds commands to the list
func addCommand(cmd *cli.Command) {
	commands = append(commands, cmd)
}

// addSubCommand adds subcommands to the list
func addSubCommand(cmdName string, subCmd *cli.Command) {
	for _, cmd := range commands {
		if cmd.Name == cmdName {
			cmd.Subcommands = append(cmd.Subcommands, subCmd)
		}
	}
}

// Execute adds all child commands to the root command. This is called by main and is considered the main entry point.
func Execute() {

	appCli = &cli.App{
		Name:     "swiz",
		Usage:    "swiz [CMD]",
		HelpName: `Swizzle your stacks together`,
		Version:  Version,
		Commands: commands,
		Before: func(ctx *cli.Context) error {
			err := genCommandPreflight(true)(ctx)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("Error: %v\n", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "appconfig",
				Usage: "specify the location of the appconfig file",
			},
		},
	}

	err := appCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
