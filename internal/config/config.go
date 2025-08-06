package config

import (
	"fmt"
	"os"
	"strings"

	"igc-tool/internal/units"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// General settings
	AltitudeUnit string `mapstructure:"altitude-unit"`
	TimeFormat   string `mapstructure:"time-format"`
	SpeedUnit    string `mapstructure:"speed-unit"`
	ClimbUnit    string `mapstructure:"climb-unit"`

	// Logbook command settings
	LogbookFormat             string  `mapstructure:"logbook-format"`
	SitesDatabaseFileLocation string  `mapstructure:"sites-database-location"`
	SpeedWindow               float64 `mapstructure:"speed-window"`

	// Internal fields (not loaded from config file)
	ConfigFile string `mapstructure:"-"`
}

// Load initializes and returns the application configuration
func Load() *Config {
	viper.SetConfigName("igc-tool")
	viper.SetConfigType("toml")

	// Look for config in various locations
	viper.AddConfigPath(".") // Current directory
	viper.AddConfigPath("$HOME/.config/igc-tool")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc/igc-tool")

	// Set environment variable prefix
	viper.SetEnvPrefix("IGC")
	viper.AutomaticEnv()

	// Replace hyphens with underscores in env vars
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Set defaults
	setDefaults()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	}

	// Unmarshal into config struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}

	// Set the config file path that was actually used
	cfg.ConfigFile = viper.ConfigFileUsed()

	return cfg
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("altitude-unit", units.AltitudeMeters)
	viper.SetDefault("time-format", units.TimeFormat24h)
	viper.SetDefault("speed-unit", units.SpeedKmh)
	viper.SetDefault("climb-unit", units.ClimbMs)
	defaultTemplate := "{{.Date}} {{.TakeoffSite}} {{.TakeoffAlt}}{{.AltitudeUnit}} {{.AltitudeDiff}}{{.AltitudeUnit}} {{.FlightDuration}} {{.MaxAltitude}}{{.AltitudeUnit}} {{.MaxGroundSpeed}}{{.SpeedUnit}} +{{.MaxClimbRate}}{{.VerticalSpeedUnit}} -{{.MaxDescentRate}}{{.VerticalSpeedUnit}}\n"
	viper.SetDefault("logbook-format", defaultTemplate)
	viper.SetDefault("sites-database-location", "")
	viper.SetDefault("speed-window", 5.0)
}
