package flags

import (
	"igc-tool/internal/config"
	"igc-tool/internal/units"

	"github.com/spf13/cobra"
)

// CommonFlags defines flags shared across multiple commands
type CommonFlags struct {
	AltitudeUnit string
	TimeFormat   string
}

// ParseFlags defines flags specific to the parse command
type ParseFlags struct {
	Summary bool
}

// LogbookFlags defines flags specific to the logbook command
type LogbookFlags struct {
	Format      string
	Sites       string
	SpeedWindow float64
	SpeedUnit   string
	ClimbUnit   string
	Recursive   bool
}

// VersionFlags defines flags specific to the version command
type VersionFlags struct {
	Detailed bool
}

// GlobalFlags defines global flags
type GlobalFlags struct {
	Version bool
}

// FlagConfig holds all flag configurations and provides unified flag resolution
type FlagConfig struct {
	cfg *config.Config
}

// FlagResolver provides a unified way to resolve flag values with proper priority
type FlagResolver struct {
	cmd *cobra.Command
	cfg *config.Config
}

// NewFlagConfig creates a new flag configuration with config reference
// NewFlagConfig creates a new flag configuration with config reference
func NewFlagConfig(cfg *config.Config) *FlagConfig {
	return &FlagConfig{
		cfg: cfg,
	}
}

// NewResolver creates a new flag resolver for a command
func (fc *FlagConfig) NewResolver(cmd *cobra.Command) *FlagResolver {
	return &FlagResolver{
		cmd: cmd,
		cfg: fc.cfg,
	}
}

// getString resolves a string flag with priority: explicit flag > config value > default
func (r *FlagResolver) getString(flagName string, configValue string) string {
	if flag := r.cmd.Flags().Lookup(flagName); flag != nil && flag.Changed {
		return flag.Value.String()
	}
	return configValue
}

// getBool resolves a bool flag with priority: explicit flag > default
func (r *FlagResolver) getBool(flagName string, defaultValue bool) bool {
	if flag := r.cmd.Flags().Lookup(flagName); flag != nil && flag.Changed {
		if val, err := r.cmd.Flags().GetBool(flagName); err == nil {
			return val
		}
	}
	return defaultValue
}

// getFloat64 resolves a float64 flag with priority: explicit flag > config value > default
func (r *FlagResolver) getFloat64(flagName string, configValue float64) float64 {
	if flag := r.cmd.Flags().Lookup(flagName); flag != nil && flag.Changed {
		if val, err := r.cmd.Flags().GetFloat64(flagName); err == nil {
			return val
		}
	}
	return configValue
}

// AddCommonFlags adds common flags to a command
// AddCommonFlags adds common flags to a command
func (fc *FlagConfig) AddCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("altitude-unit", "a", fc.cfg.AltitudeUnit, "Unit for altitude display ("+units.AltitudeMeters+", "+units.AltitudeFeet+")")
	cmd.Flags().StringP("time-format", "t", fc.cfg.TimeFormat, "Time format ("+units.TimeFormat24h+", "+units.TimeFormatAMPM+")")
}

// AddParseFlags adds parse-specific flags to a command
func (fc *FlagConfig) AddParseFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("summary", false, "Show only headers and first/last fixes instead of all fixes")
}

// AddLogbookFlags adds logbook-specific flags to a command
func (fc *FlagConfig) AddLogbookFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", fc.cfg.LogbookFormat, "Go template string for formatting the output")
	cmd.Flags().StringP("sites", "s", fc.cfg.SitesDatabaseFileLocation, "Path to GeoJSON file containing landing site definitions")
	cmd.Flags().Float64P("speed-window", "w", fc.cfg.SpeedWindow, "Time window in seconds for ground speed calculations (larger values reduce GPS noise)")
	cmd.Flags().StringP("speed-unit", "u", fc.cfg.SpeedUnit, "Unit for speed display ("+units.SpeedKmh+", "+units.SpeedMph+", "+units.SpeedKnots+", "+units.SpeedMs+")")
	cmd.Flags().StringP("climb-unit", "c", fc.cfg.ClimbUnit, "Unit for climb rate display ("+units.ClimbMs+", "+units.ClimbFpm+")")
	cmd.Flags().BoolP("recursive", "r", false, "Recursively search for IGC files in directories")
}

// AddVersionFlags adds version-specific flags to a command
func (fc *FlagConfig) AddVersionFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("detailed", "d", false, "Show detailed version information including build details")
}

// AddGlobalFlags adds global flags to a command
func (fc *FlagConfig) AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP("version", "v", false, "Show version information")
}

// GetCommonFromConfig retrieves common flag values, preferring runtime flag values over config defaults
func (fc *FlagConfig) GetCommonFromConfig(cmd *cobra.Command, cfg *config.Config) CommonFlags {
	resolver := fc.NewResolver(cmd)
	return CommonFlags{
		AltitudeUnit: resolver.getString("altitude-unit", cfg.AltitudeUnit),
		TimeFormat:   resolver.getString("time-format", cfg.TimeFormat),
	}
}

// GetParseFromFlags retrieves parse flag values from cobra command
func (fc *FlagConfig) GetParseFromFlags(cmd *cobra.Command) ParseFlags {
	resolver := fc.NewResolver(cmd)
	return ParseFlags{
		Summary: resolver.getBool("summary", false),
	}
}

// GetLogbookFromConfig retrieves logbook flag values, preferring runtime flag values over config defaults
func (fc *FlagConfig) GetLogbookFromConfig(cmd *cobra.Command, cfg *config.Config) LogbookFlags {
	resolver := fc.NewResolver(cmd)
	return LogbookFlags{
		Format:      resolver.getString("format", cfg.LogbookFormat),
		Sites:       resolver.getString("sites", cfg.SitesDatabaseFileLocation),
		SpeedWindow: resolver.getFloat64("speed-window", cfg.SpeedWindow),
		SpeedUnit:   resolver.getString("speed-unit", cfg.SpeedUnit),
		ClimbUnit:   resolver.getString("climb-unit", cfg.ClimbUnit),
		Recursive:   resolver.getBool("recursive", false),
	}
}

// GetVersionFromFlags retrieves version flag values from cobra command
func (fc *FlagConfig) GetVersionFromFlags(cmd *cobra.Command) VersionFlags {
	resolver := fc.NewResolver(cmd)
	return VersionFlags{
		Detailed: resolver.getBool("detailed", false),
	}
}

// GetGlobalFromFlags retrieves global flag values from cobra command
func (fc *FlagConfig) GetGlobalFromFlags(cmd *cobra.Command) GlobalFlags {
	resolver := fc.NewResolver(cmd)
	return GlobalFlags{
		Version: resolver.getBool("version", false),
	}
}

// GetAllFlags retrieves all flag values for a command, preferring runtime flags over config defaults
func (fc *FlagConfig) GetAllFlags(cmd *cobra.Command, cfg *config.Config) (CommonFlags, LogbookFlags, ParseFlags, VersionFlags) {
	common := fc.GetCommonFromConfig(cmd, cfg)
	logbook := fc.GetLogbookFromConfig(cmd, cfg)
	parse := fc.GetParseFromFlags(cmd)
	version := fc.GetVersionFromFlags(cmd)

	return common, logbook, parse, version
}
