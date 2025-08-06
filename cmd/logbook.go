package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"igc-tool/internal/cli"
	"igc-tool/internal/config"
	"igc-tool/internal/flags"
	"igc-tool/internal/logbook"
	"igc-tool/internal/parser"
	"igc-tool/internal/utils"

	"github.com/spf13/cobra"
)

// NewLogbookCmd creates and returns the logbook command
func NewLogbookCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var logbookCmd = &cobra.Command{
		Use:   "logbook [IGC files or directories...]",
		Short: "Generate logbook entries for flights",
		Long:  fmt.Sprintf("Generate logbook entries for flights. Accepts multiple IGC files and/or directories containing IGC files. Use --format to customize output with Go templates.\n\nAvailable template fields: %s", strings.Join(logbook.GetDataFields(), ", ")),
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logbookFlags := flagConfig.GetLogbookFromConfig(cmd, cfg)
			commonFlags := flagConfig.GetCommonFromConfig(cmd, cfg)

			// Load landing sites if specified
			landingSites, err := cli.LoadLandingSitesIfSpecified(logbookFlags.Sites)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading landing sites: %v\n", err)
				os.Exit(1)
			}

			// Find all IGC files from the provided arguments
			igcFiles, err := cli.FindIGCFiles(args, logbookFlags.Recursive)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding IGC files: %v\n", err)
				os.Exit(1)
			}

			if len(igcFiles) == 0 {
				fmt.Fprintf(os.Stderr, "No IGC files found\n")
				os.Exit(1)
			}

			// Track total flight time for multiple files
			var totalFlightTime time.Duration
			processedCount := 0

			// Process each IGC file
			for _, filename := range igcFiles {
				flight, err := parser.ParseIGCFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", filename, err)
					continue
				}

				// Create options using flag values
				opts := logbook.Options{
					LandingSites: landingSites,
					Filename:     filename,
					SpeedWindow:  logbookFlags.SpeedWindow,
					AltitudeUnit: commonFlags.AltitudeUnit,
					SpeedUnit:    logbookFlags.SpeedUnit,
					ClimbUnit:    logbookFlags.ClimbUnit,
					TimeFormat:   commonFlags.TimeFormat,
				}
				data := logbook.CreateData(flight, opts)

				err = cli.PrintTemplatedLogbook(data, logbookFlags.Format)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
					continue
				}

				// Add flight duration to total if we have fixes
				if len(flight.Fixes) > 0 {
					flightDuration := flight.Fixes[len(flight.Fixes)-1].Time.Sub(flight.Fixes[0].Time)
					totalFlightTime += flightDuration
					processedCount++
				}
			}

			// Print total flight time if more than one file was processed
			if processedCount > 1 {
				fmt.Printf("# total flight time: %s\n", utils.FormatDuration(totalFlightTime))
			}
		},
	}

	// Set up flags
	flagConfig.AddLogbookFlags(logbookCmd)
	flagConfig.AddCommonFlags(logbookCmd)

	return logbookCmd
}
