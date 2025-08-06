package cmd

import (
	"fmt"

	"igc-tool/internal/config"
	"igc-tool/internal/flags"
	"igc-tool/internal/version"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates and returns the version command
func NewVersionCmd(cfg *config.Config, flagConfig *flags.FlagConfig) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display version, build information, and platform details.`,
		Run: func(cmd *cobra.Command, args []string) {
			buildInfo := version.GetBuildInfo()
			versionFlags := flagConfig.GetVersionFromFlags(cmd)
			if versionFlags.Detailed {
				fmt.Println(buildInfo.DetailedString())
			} else {
				fmt.Println(buildInfo.String())
			}
		},
	}

	// Set up flags
	flagConfig.AddVersionFlags(versionCmd)

	return versionCmd
}
