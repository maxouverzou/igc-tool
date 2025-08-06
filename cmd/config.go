package cmd

import (
	"fmt"

	"igc-tool/internal/config"
	"igc-tool/internal/flags"

	"github.com/spf13/cobra"
)

// NewConfigCmd creates and returns the config command
func NewConfigCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Long:  `Display the current configuration values from config files, environment variables, and defaults.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Current configuration:")
			configFile := cfg.ConfigFile
			if configFile == "" {
				configFile = "No config file found (using defaults)"
			}
			fmt.Printf("Config file used: %s\n\n", configFile)

			commonFlags := flagConfig.GetCommonFromConfig(cmd, cfg)
			logbookFlags := flagConfig.GetLogbookFromConfig(cmd, cfg)

			fmt.Printf("altitude-unit: %s\n", commonFlags.AltitudeUnit)
			fmt.Printf("time-format: %s\n", commonFlags.TimeFormat)
			fmt.Printf("speed-unit: %s\n", logbookFlags.SpeedUnit)
			fmt.Printf("climb-unit: %s\n", logbookFlags.ClimbUnit)
			fmt.Printf("logbook-format: %s\n", logbookFlags.Format)
			fmt.Printf("sites-database-location: %s\n", logbookFlags.Sites)
			fmt.Printf("speed-window: %g\n", logbookFlags.SpeedWindow)
		},
	}

	return configCmd
}
