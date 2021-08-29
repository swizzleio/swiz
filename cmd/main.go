package main

import (
	"getswizzle.io/swiz/internal/cmd"
	"getswizzle.io/swiz/internal/config"
)

// main is the main entry point
func main() {
	// Load config
	cfg := config.LoadConfigMust()

	// Execute CLI
	cmd.Execute(cfg)
}
