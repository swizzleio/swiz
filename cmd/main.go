package main

import (
	"github.com/swizzleio/swiz/internal/cmd"
)

var (
	Version    = "dev"
	CommitHash = "n/a"
)

// main is the main entry point
func main() {
	// Execute CLI
	cmd.Execute()
}
