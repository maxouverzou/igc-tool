package cmd

import (
	"fmt"
	"os"

	"igc-tool/internal/config"
	"igc-tool/internal/display"
	"igc-tool/internal/flags"
	"igc-tool/internal/parser"

	"github.com/spf13/cobra"
)

// NewParseCmd creates and returns the parse command
func NewParseCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var parseCmd = &cobra.Command{
		Use:   "parse [IGC file]",
		Short: "Parse and display detailed IGC flight data",
		Long:  `Parse an IGC file and display all flight information including fixes, waypoints, and metadata.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			parseFlags := flagConfig.GetParseFromFlags(cmd)
			commonFlags := flagConfig.GetCommonFromConfig(cmd, cfg)

			flight, err := parser.ParseIGCFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			display.PrintFlightData(flight, parseFlags.Summary, commonFlags.AltitudeUnit, commonFlags.TimeFormat)
		},
	}

	// Set up flags
	flagConfig.AddParseFlags(parseCmd)
	flagConfig.AddCommonFlags(parseCmd)

	return parseCmd
}
