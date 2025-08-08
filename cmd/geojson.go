package cmd

import (
	"fmt"
	"os"

	"igc-tool/internal/config"
	"igc-tool/internal/flags"
	"igc-tool/internal/geojson"
	"igc-tool/internal/parser"

	"github.com/spf13/cobra"
)

// NewGeoJSONCmd creates and returns the geojson command
func NewGeoJSONCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var geojsonCmd = &cobra.Command{
		Use:   "geojson [IGC file]",
		Short: "Convert IGC flight track to GeoJSON",
		Long:  `Parse an IGC file and convert the flight track to a GeoJSON LineString feature.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			renderFlags := flagConfig.GetRenderFromFlags(cmd)

			flight, err := parser.ParseIGCFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			geojsonData, err := geojson.RenderToGeoJSON(flight, renderFlags.Pretty, renderFlags.IncludeMetadata)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering GeoJSON: %v\n", err)
				os.Exit(1)
			}

			if renderFlags.Output != "" {
				err := os.WriteFile(renderFlags.Output, geojsonData, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", renderFlags.Output, err)
					os.Exit(1)
				}
				fmt.Fprintf(os.Stderr, "GeoJSON written to %s\n", renderFlags.Output)
			} else {
				fmt.Print(string(geojsonData))
			}
		},
	}

	// Set up flags
	flagConfig.AddRenderFlags(geojsonCmd)

	return geojsonCmd
}
