package cmd

import (
	"fmt"
	"os"
	"strings"

	"igc-tool/internal/cli"
	"igc-tool/internal/config"
	"igc-tool/internal/flags"
	"igc-tool/internal/logbook"
	"igc-tool/internal/parser"

	"github.com/spf13/cobra"
)

// NewLogbookCmd creates and returns the logbook command
func NewLogbookCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var logbookCmd = &cobra.Command{
		Use:   "logbook [IGC files or directories...]",
		Short: "Generate logbook entries for flights",
		Long: fmt.Sprintf(`Generate logbook entries for flights. Accepts multiple IGC files and/or directories containing IGC files.

Template Variables (always available):
  Individual flight fields (access via .Flights array): %s
  
  Aggregated statistics: %s

Examples:
  # Basic usage (single flight)
  igc-tool logbook flight1.igc
  
  # Basic usage (multiple flights)
  igc-tool logbook flight1.igc flight2.igc
  
  # Custom format for individual flights
  igc-tool logbook --format "{{range .Flights}}{{.Date}}: {{.FlightDuration}} at {{.TakeoffSite}}\n{{end}}" *.igc
  
  # Access aggregated data in templates
  igc-tool logbook --format "{{range .Flights}}{{.Date}} {{.FlightDuration}}\n{{end}}Total: {{.TotalTime}}" *.igc
  
  # Show only summary statistics
  igc-tool logbook --format "Summary: {{.TotalFlights}} flights, {{.TotalTime}} total time\n" *.igc
  
  # Mix individual and aggregated data
  igc-tool logbook --format "Flights:\n{{range .Flights}}- {{.Date}}: {{.FlightDuration}}\n{{end}}Total time: {{.TotalTime}}\n" *.igc`,
			strings.Join(logbook.GetDataFields(), ", "),
			strings.Join(logbook.GetTemplateDataFields(), ", ")),
		Args: cobra.MinimumNArgs(1),
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

			// Collect all flight data
			var allFlights []*logbook.Data
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
				if data != nil {
					allFlights = append(allFlights, data)
					processedCount++
				}
			}

			if processedCount == 0 {
				fmt.Fprintf(os.Stderr, "No valid flights found\n")
				os.Exit(1)
			}

			// Always use TemplateData for consistent template variables
			templateData := logbook.CreateTemplateData(allFlights, logbook.Options{
				AltitudeUnit: commonFlags.AltitudeUnit,
				SpeedUnit:    logbookFlags.SpeedUnit,
				ClimbUnit:    logbookFlags.ClimbUnit,
			})

			// Use the template as-is - no automatic wrapping
			templateStr := logbookFlags.Format

			err = cli.PrintTemplatedLogbookData(templateData, templateStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Set up flags
	flagConfig.AddLogbookFlags(logbookCmd)
	flagConfig.AddCommonFlags(logbookCmd)

	return logbookCmd
}
