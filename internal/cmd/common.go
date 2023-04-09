package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandPreflight(ctx *cli.Context) error {
	fmt.Println("checking preconditions")
	return nil
}
