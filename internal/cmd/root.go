package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	Version    = "dev"
	CommitHash = "n/a"
)

// Internal list of commands
var commands = []*cli.Command{}

// addCommand adds commands to the list
func addCommand(cmd *cli.Command) {
	commands = append(commands, cmd)
}

// Execute adds all child commands to the root command. This is called by main and is considered the main entry point.
func Execute() {

	app := &cli.App{
		Name:     "swiz",
		Usage:    "swiz [CMD]",
		HelpName: `Swizzle your stacks together`,
		Version:  Version,
		Commands: commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
