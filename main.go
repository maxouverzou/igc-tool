package main

import (
	"fmt"
	"os"

	"igc-tool/cmd"
	"igc-tool/internal/config"
	"igc-tool/internal/flags"
)

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Initialize centralized flag configuration
	flagConfig := flags.NewFlagConfig(cfg)

	// Create root command with all subcommands
	rootCmd := cmd.NewRootCmd(cfg, flagConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
