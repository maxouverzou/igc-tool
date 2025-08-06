package flags

import (
	"testing"

	"igc-tool/internal/config"

	"github.com/spf13/cobra"
)

func TestNewFlagConfig(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit:              "ft",
		TimeFormat:                "ampm",
		SpeedUnit:                 "mph",
		ClimbUnit:                 "fpm",
		LogbookFormat:             "test-format",
		SitesDatabaseFileLocation: "test-sites.json",
		SpeedWindow:               10.0,
	}

	fc := NewFlagConfig(cfg)

	// Test that the FlagConfig was created successfully
	if fc == nil {
		t.Error("expected FlagConfig to be created, got nil")
	}

	// Create a mock command to test flag resolution
	cmd := &cobra.Command{}
	fc.AddCommonFlags(cmd)
	fc.AddLogbookFlags(cmd)

	// Test common flags from config
	common := fc.GetCommonFromConfig(cmd, cfg)
	if common.AltitudeUnit != "ft" {
		t.Errorf("expected AltitudeUnit 'ft', got '%s'", common.AltitudeUnit)
	}

	if common.TimeFormat != "ampm" {
		t.Errorf("expected TimeFormat 'ampm', got '%s'", common.TimeFormat)
	}

	// Test logbook flags from config
	logbook := fc.GetLogbookFromConfig(cmd, cfg)
	if logbook.SpeedUnit != "mph" {
		t.Errorf("expected SpeedUnit 'mph', got '%s'", logbook.SpeedUnit)
	}

	if logbook.SpeedWindow != 10.0 {
		t.Errorf("expected SpeedWindow 10.0, got %f", logbook.SpeedWindow)
	}
}

func TestAddCommonFlags(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit: "m",
		TimeFormat:   "24h",
	}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddCommonFlags(cmd)

	altitudeFlag := cmd.Flags().Lookup("altitude-unit")
	if altitudeFlag == nil {
		t.Error("altitude-unit flag not found")
	} else if altitudeFlag.DefValue != "m" {
		t.Errorf("expected altitude-unit default 'm', got '%s'", altitudeFlag.DefValue)
	}

	timeFlag := cmd.Flags().Lookup("time-format")
	if timeFlag == nil {
		t.Error("time-format flag not found")
	} else if timeFlag.DefValue != "24h" {
		t.Errorf("expected time-format default '24h', got '%s'", timeFlag.DefValue)
	}
}

func TestAddLogbookFlags(t *testing.T) {
	cfg := &config.Config{
		LogbookFormat:             "test-format",
		SitesDatabaseFileLocation: "test-sites.json",
		SpeedWindow:               5.0,
		SpeedUnit:                 "kmh",
		ClimbUnit:                 "ms",
	}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddLogbookFlags(cmd)

	formatFlag := cmd.Flags().Lookup("format")
	if formatFlag == nil {
		t.Error("format flag not found")
	} else if formatFlag.DefValue != "test-format" {
		t.Errorf("expected format default 'test-format', got '%s'", formatFlag.DefValue)
	}

	speedWindowFlag := cmd.Flags().Lookup("speed-window")
	if speedWindowFlag == nil {
		t.Error("speed-window flag not found")
	} else if speedWindowFlag.DefValue != "5" {
		t.Errorf("expected speed-window default '5', got '%s'", speedWindowFlag.DefValue)
	}

	recursiveFlag := cmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("recursive flag not found")
	}
}

func TestGetCommonFromConfig(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit: "m",
		TimeFormat:   "24h",
	}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddCommonFlags(cmd)

	// Set flag values
	cmd.Flags().Set("altitude-unit", "ft")
	cmd.Flags().Set("time-format", "ampm")

	common := fc.GetCommonFromConfig(cmd, cfg)

	if common.AltitudeUnit != "ft" {
		t.Errorf("expected AltitudeUnit 'ft', got '%s'", common.AltitudeUnit)
	}

	if common.TimeFormat != "ampm" {
		t.Errorf("expected TimeFormat 'ampm', got '%s'", common.TimeFormat)
	}
}

func TestGetParseFromFlags(t *testing.T) {
	cfg := &config.Config{}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddParseFlags(cmd)

	// Set summary flag to true
	cmd.Flags().Set("summary", "true")

	parse := fc.GetParseFromFlags(cmd)

	if !parse.Summary {
		t.Error("expected Summary to be true")
	}
}

func TestGetLogbookFromConfig(t *testing.T) {
	cfg := &config.Config{
		LogbookFormat:             "test-format",
		SitesDatabaseFileLocation: "test-sites.json",
		SpeedWindow:               5.0,
		SpeedUnit:                 "kmh",
		ClimbUnit:                 "ms",
	}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddLogbookFlags(cmd)

	// Test with no flags set - should use config defaults
	logbook := fc.GetLogbookFromConfig(cmd, cfg)
	if logbook.Format != "test-format" {
		t.Errorf("expected Format 'test-format', got '%s'", logbook.Format)
	}
	if logbook.SpeedWindow != 5.0 {
		t.Errorf("expected SpeedWindow 5.0, got %f", logbook.SpeedWindow)
	}
	if logbook.Recursive {
		t.Error("expected Recursive to be false by default")
	}

	// Test with flags set - should use flag values
	cmd.Flags().Set("format", "custom-format")
	cmd.Flags().Set("speed-window", "10.0")
	cmd.Flags().Set("recursive", "true")

	logbook = fc.GetLogbookFromConfig(cmd, cfg)
	if logbook.Format != "custom-format" {
		t.Errorf("expected Format 'custom-format', got '%s'", logbook.Format)
	}
	if logbook.SpeedWindow != 10.0 {
		t.Errorf("expected SpeedWindow 10.0, got %f", logbook.SpeedWindow)
	}
	if !logbook.Recursive {
		t.Error("expected Recursive to be true")
	}
}

func TestGetAllFlags(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit:              "m",
		TimeFormat:                "24h",
		LogbookFormat:             "test-format",
		SitesDatabaseFileLocation: "test-sites.json",
		SpeedWindow:               5.0,
		SpeedUnit:                 "kmh",
		ClimbUnit:                 "ms",
	}
	fc := NewFlagConfig(cfg)

	cmd := &cobra.Command{}
	fc.AddCommonFlags(cmd)
	fc.AddLogbookFlags(cmd)
	fc.AddParseFlags(cmd)
	fc.AddVersionFlags(cmd)

	cmd.Flags().Set("altitude-unit", "ft")
	cmd.Flags().Set("summary", "true")

	common, logbook, parse, version := fc.GetAllFlags(cmd, cfg)

	if common.AltitudeUnit != "ft" {
		t.Errorf("expected AltitudeUnit 'ft', got '%s'", common.AltitudeUnit)
	}
	if !parse.Summary {
		t.Error("expected Summary to be true")
	}
	if logbook.Format != "test-format" {
		t.Errorf("expected Format 'test-format', got '%s'", logbook.Format)
	}
	if version.Detailed {
		t.Error("expected Detailed to be false by default")
	}
}
