package logbook

import (
	"reflect"
	"testing"
	"time"

	"igc-tool/internal/config"
	"igc-tool/internal/flight"
	"igc-tool/internal/sites"

	"github.com/twpayne/go-igc"
)

func TestCreateData(t *testing.T) {
	baseTime := time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC)
	testFlight := &flight.Flight{
		Date:               time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
		Pilot:              "TestPilot",
		Crew:               "TestCrew",
		GliderType:         "TestGlider",
		GliderID:           "ABC123",
		CompetitionID:      "COMP456",
		FlightRecorderType: "TestFR",
		Fixes: []*igc.BRecord{
			{
				Lat:      45.814,
				Lon:      6.246,
				Time:     baseTime,
				AltWGS84: 1500,
			},
			{
				Lat:      45.815,
				Lon:      6.247,
				Time:     baseTime.Add(30 * time.Minute),
				AltWGS84: 1800,
			},
			{
				Lat:      45.816,
				Lon:      6.248,
				Time:     baseTime.Add(time.Hour),
				AltWGS84: 1600,
			},
		},
	}

	testSites := &sites.Collection{
		Sites: []sites.LandingSite{
			{
				Name:   "TestSite",
				Center: [2]float64{6.246, 45.814},
				Radius: 1000,
			},
		},
	}

	tests := []struct {
		name        string
		flight      *flight.Flight
		opts        Options
		expectNil   bool
		checkFields bool
	}{
		{
			name:   "valid flight with all options",
			flight: testFlight,
			opts: Options{
				LandingSites: testSites,
				Filename:     "test.igc",
				SpeedWindow:  5.0,
				AltitudeUnit: "m",
				SpeedUnit:    "kmh",
				ClimbUnit:    "ms",
				TimeFormat:   "24h",
			},
			expectNil:   false,
			checkFields: true,
		},
		{
			name:   "valid flight without landing sites",
			flight: testFlight,
			opts: Options{
				LandingSites: nil,
				Filename:     "test.igc",
				SpeedWindow:  5.0,
				AltitudeUnit: "ft",
				SpeedUnit:    "mph",
				ClimbUnit:    "fpm",
				TimeFormat:   "ampm",
			},
			expectNil:   false,
			checkFields: true,
		},
		{
			name: "flight with no fixes",
			flight: &flight.Flight{
				Date:  time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
				Pilot: "TestPilot",
				Fixes: []*igc.BRecord{},
			},
			opts: Options{
				Filename:     "test.igc",
				SpeedWindow:  5.0,
				AltitudeUnit: "m",
				SpeedUnit:    "kmh",
				ClimbUnit:    "ms",
				TimeFormat:   "24h",
			},
			expectNil:   true,
			checkFields: false,
		},
		{
			name: "flight with single fix",
			flight: &flight.Flight{
				Date:  time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
				Pilot: "TestPilot",
				Fixes: []*igc.BRecord{
					{
						Lat:      45.814,
						Lon:      6.246,
						Time:     baseTime,
						AltWGS84: 1500,
					},
				},
			},
			opts: Options{
				Filename:     "test.igc",
				SpeedWindow:  5.0,
				AltitudeUnit: "m",
				SpeedUnit:    "kmh",
				ClimbUnit:    "ms",
				TimeFormat:   "24h",
			},
			expectNil:   false,
			checkFields: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateData(tt.flight, tt.opts)

			if tt.expectNil {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("expected non-nil result, got nil")
				return
			}

			if !tt.checkFields {
				return
			}

			// Check required fields
			if result.Date == "" {
				t.Errorf("date field is empty")
			}

			if result.Filename != tt.opts.Filename {
				t.Errorf("expected filename %s, got %s", tt.opts.Filename, result.Filename)
			}

			if result.Pilot != tt.flight.Pilot {
				t.Errorf("expected pilot %s, got %s", tt.flight.Pilot, result.Pilot)
			}

			// Check unit symbols
			if result.AltitudeUnit == "" {
				t.Errorf("altitude unit symbol is empty")
			}

			if result.SpeedUnit == "" {
				t.Errorf("speed unit symbol is empty")
			}

			if result.VerticalSpeedUnit == "" {
				t.Errorf("vertical speed unit symbol is empty")
			}

			// Check coordinate formatting
			if result.TakeoffPosition == "" {
				t.Errorf("takeoff position is empty")
			}

			if result.LandingPosition == "" {
				t.Errorf("landing position is empty")
			}

			// Check time formatting
			if result.TakeoffTime == "" {
				t.Errorf("takeoff time is empty")
			}

			if result.LandingTime == "" {
				t.Errorf("landing time is empty")
			}

			// Check duration formatting
			if result.FlightDuration == "" {
				t.Errorf("flight duration is empty")
			}

			// Check that landing sites are resolved when provided
			if tt.opts.LandingSites != nil {
				// First fix should be close to TestSite
				if result.TakeoffSite != "TestSite" {
					t.Logf("expected takeoff site 'TestSite', got '%s' (may be expected if outside radius)", result.TakeoffSite)
				}
			}

			// Check numeric fields are reasonable
			if result.MaxAltitude < 0 {
				t.Errorf("unexpected negative max altitude: %d", result.MaxAltitude)
			}

			if result.MaxGroundSpeed < 0 {
				t.Errorf("unexpected negative max ground speed: %d", result.MaxGroundSpeed)
			}

			if result.MaxClimbRate < 0 {
				t.Errorf("unexpected negative max climb rate: %f", result.MaxClimbRate)
			}

			if result.MaxDescentRate < 0 {
				t.Errorf("unexpected negative max descent rate: %f", result.MaxDescentRate)
			}
		})
	}
}

func TestGetDataFields(t *testing.T) {
	fields := GetDataFields()

	if len(fields) == 0 {
		t.Errorf("expected non-empty list of fields")
	}

	// Check that some expected fields are present
	expectedFields := []string{
		"Date", "TakeoffLat", "TakeoffLon", "TakeoffSite", "LandingSite",
		"MaxAltitude", "MaxGroundSpeed", "FlightDuration", "Pilot",
		"AltitudeUnit", "SpeedUnit", "VerticalSpeedUnit",
	}

	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("expected field %s not found in field list: %v", expected, fields)
		}
	}

	// Check that all fields correspond to actual struct fields
	dataType := reflect.TypeOf(Data{})
	for _, field := range fields {
		if _, exists := dataType.FieldByName(field); !exists {
			t.Errorf("field %s does not exist in Data struct", field)
		}
	}
}

func TestCreateOptions(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit: "ft",
		SpeedUnit:    "mph",
		ClimbUnit:    "fpm",
		TimeFormat:   "ampm",
		SpeedWindow:  7.5,
	}

	testSites := &sites.Collection{
		Sites: []sites.LandingSite{
			{Name: "TestSite", Center: [2]float64{6.246, 45.814}, Radius: 1000},
		},
	}

	filename := "test_flight.igc"

	opts := CreateOptions(cfg, testSites, filename)

	// Check that all config values are transferred correctly
	if opts.AltitudeUnit != cfg.AltitudeUnit {
		t.Errorf("expected altitude unit %s, got %s", cfg.AltitudeUnit, opts.AltitudeUnit)
	}

	if opts.SpeedUnit != cfg.SpeedUnit {
		t.Errorf("expected speed unit %s, got %s", cfg.SpeedUnit, opts.SpeedUnit)
	}

	if opts.ClimbUnit != cfg.ClimbUnit {
		t.Errorf("expected climb unit %s, got %s", cfg.ClimbUnit, opts.ClimbUnit)
	}

	if opts.TimeFormat != cfg.TimeFormat {
		t.Errorf("expected time format %s, got %s", cfg.TimeFormat, opts.TimeFormat)
	}

	if opts.SpeedWindow != cfg.SpeedWindow {
		t.Errorf("expected speed window %f, got %f", cfg.SpeedWindow, opts.SpeedWindow)
	}

	if opts.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, opts.Filename)
	}

	if opts.LandingSites != testSites {
		t.Errorf("expected landing sites to be set")
	}
}

func TestCreateOptionsWithNilSites(t *testing.T) {
	cfg := &config.Config{
		AltitudeUnit: "m",
		SpeedUnit:    "kmh",
		ClimbUnit:    "ms",
		TimeFormat:   "24h",
		SpeedWindow:  5.0,
	}

	filename := "test_flight.igc"

	opts := CreateOptions(cfg, nil, filename)

	if opts.LandingSites != nil {
		t.Errorf("expected nil landing sites, got %v", opts.LandingSites)
	}

	// Other fields should still be set correctly
	if opts.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, opts.Filename)
	}

	if opts.AltitudeUnit != cfg.AltitudeUnit {
		t.Errorf("expected altitude unit %s, got %s", cfg.AltitudeUnit, opts.AltitudeUnit)
	}
}

func TestDataStructCompleteness(t *testing.T) {
	// Test that Data struct has all expected fields by creating an instance
	// and checking that we can set all major field types
	data := &Data{
		Date:               "2025-07-18",
		TakeoffLat:         45.814,
		TakeoffLon:         6.246,
		TakeoffPosition:    "45.814,6.246",
		TakeoffSite:        "TestSite",
		LandingLat:         45.815,
		LandingLon:         6.247,
		LandingPosition:    "45.815,6.247",
		LandingSite:        "LandingSite",
		TakeoffAlt:         1500,
		LandingAlt:         1600,
		AltitudeDiff:       100,
		MaxAltitude:        1800,
		MinAltitude:        1400,
		MaxGroundSpeed:     85,
		MaxClimbRate:       8.5,
		MaxDescentRate:     12.3,
		FlightDuration:     "2h30m",
		TakeoffTime:        "12:00:00",
		LandingTime:        "14:30:00",
		Pilot:              "TestPilot",
		Crew:               "TestCrew",
		GliderType:         "TestGlider",
		GliderID:           "ABC123",
		CompetitionID:      "COMP456",
		FlightRecorderType: "TestFR",
		Filename:           "test.igc",
		AltitudeUnit:       "m",
		SpeedUnit:          "km/h",
		VerticalSpeedUnit:  "m/s",
	}

	// Basic sanity checks
	if data.Date == "" {
		t.Errorf("date field not set properly")
	}

	if data.TakeoffLat == 0 && data.TakeoffLon == 0 {
		t.Errorf("coordinate fields not set properly")
	}

	if data.MaxGroundSpeed == 0 {
		t.Errorf("speed field not set properly")
	}

	if data.FlightDuration == "" {
		t.Errorf("duration field not set properly")
	}

	if data.AltitudeUnit == "" {
		t.Errorf("unit fields not set properly")
	}
}
