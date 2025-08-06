package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time using ldflags
var (
	Version   = "dev"     // The version of the application
	GitCommit = "unknown" // The git commit hash
	BuildDate = "unknown" // The build date
	GoVersion = runtime.Version()
)

// BuildInfo returns detailed build information
type BuildInfo struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Compiler  string
	Platform  string
}

// GetBuildInfo returns the current build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (b BuildInfo) String() string {
	return fmt.Sprintf("igc-tool version %s", b.Version)
}

// DetailedString returns a detailed version string with all build info
func (b BuildInfo) DetailedString() string {
	return fmt.Sprintf(`igc-tool version %s
Git commit: %s
Build date: %s
Go version: %s
Compiler: %s
Platform: %s`,
		b.Version, b.GitCommit, b.BuildDate, b.GoVersion, b.Compiler, b.Platform)
}
