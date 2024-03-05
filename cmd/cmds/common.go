package cmds

import (
	"github.com/urfave/cli/v2"
	"strings"
)

func genCommandPreflight(appConfigSoftFail bool) cli.BeforeFunc {

	return func(ctx *cli.Context) error {
		// Validate app config
		appConfigLoc := ctx.String("appconfig")
		_, err := appConfigMgr.Load(appConfigLoc)
		if err != nil {
			if appConfigSoftFail {
				return nil
			}

			return err
		}

		for _, cmd := range appConfigMgr.Get().DisabledCommands {

			// Split command into parts in case there is a subcommand
			cmdList := strings.Split(cmd, ".")
			if len(cmdList) == 1 {
				// Disable command
				for _, c := range appCli.Commands {
					if c.Name == cmd {
						c.Hidden = true
					}
				}
			} else {
				// Disable subcommand
				for _, c := range appCli.Commands {
					if c.Name == cmdList[0] {
						for _, subCmd := range c.Subcommands {
							if subCmd.Name == cmdList[1] {
								subCmd.Hidden = true
							}
						}
					}
				}
			}
		}

		return nil
	}
}
