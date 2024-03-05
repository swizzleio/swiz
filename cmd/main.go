package main

import (
	"github.com/swizzleio/swiz/cmd/cmds"
	"os"
)

// main is the main entry point
func main() {
	// Execute CLI
	os.Exit(cmds.Execute())
}
