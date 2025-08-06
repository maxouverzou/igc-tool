package cmd

import (
	"fmt"

	"igc-tool/internal/config"
	"igc-tool/internal/flags"
	"igc-tool/internal/version"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command
func NewRootCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "igc-tool",
		Short: "Parse and display IGC flight data",
		Long:  `A tool to parse IGC (International Gliding Commission) flight files and display flight information including fixes, waypoints, and metadata.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Handle global version flag when no subcommand is provided
			globalFlags := flagConfig.GetGlobalFromFlags(cmd)
			if globalFlags.Version {
				fmt.Println(version.GetBuildInfo().String())
				return
			}
			// Show help if no args provided
			cmd.Help()
		},
	}

	// Add global version flag
	flagConfig.AddGlobalFlags(rootCmd)

	// Add subcommands
	rootCmd.AddCommand(NewParseCmd(cfg, flagConfig))
	rootCmd.AddCommand(NewLogbookCmd(cfg, flagConfig))
	rootCmd.AddCommand(NewConfigCmd(cfg, flagConfig))
	rootCmd.AddCommand(NewVersionCmd(cfg, flagConfig))

	return rootCmd
}
