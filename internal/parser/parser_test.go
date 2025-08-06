package parser

import (
	"os"
	"testing"
	"time"

	"github.com/twpayne/go-igc"
)

func TestGetHRecordValue(t *testing.T) {
	tests := []struct {
		name     string
		records  map[string]*igc.HRecord
		key      string
		expected string
	}{
		{
			name: "existing record",
			records: map[string]*igc.HRecord{
				"PLT": {Value: "TestPilot"},
			},
			key:      "PLT",
			expected: "TestPilot",
		},
		{
			name:     "non-existing record",
			records:  map[string]*igc.HRecord{},
			key:      "PLT",
			expected: "",
		},
		{
			name: "nil record",
			records: map[string]*igc.HRecord{
				"PLT": nil,
			},
			key:      "PLT",
			expected: "",
		},
		{
			name: "empty value",
			records: map[string]*igc.HRecord{
				"PLT": {Value: ""},
			},
			key:      "PLT",
			expected: "",
		},
		{
			name: "multiple records",
			records: map[string]*igc.HRecord{
				"PLT": {Value: "TestPilot"},
				"GTY": {Value: "TestGlider"},
				"GID": {Value: "ABC123"},
			},
			key:      "GTY",
			expected: "TestGlider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHRecordValue(tt.records, tt.key)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestParseIGCFileWithMockData tests the parsing logic with a temporary IGC file
func TestParseIGCFileWithMockData(t *testing.T) {
	// Create a temporary IGC file for testing
	igcContent := `AXSDUB54EB
HFDTE300723
HFPLTPILOTINCHARGE:TestPilot
HFCM2CREW2:NIL
HFGTYGLIDERTYPE:ACME Glider
HFGIDGLIDERID:ABC123
HFCIDCOMPETITIONID:COM123
HFDTMGPSDATUM:WGS84
HFRFWFIRMWAREVERSION:2023-06-30:1fe35e9c
HFRHWHARDWAREVERSION:ULTRABIP 1.0
HFFTYFRTYPE:STODEUS,ULTRABIP
HFGPSRECEIVER:GOTOP,GT1110SN,22,18000
HFTZNTIMEZONE:1
HFPRSPRESSALTSENSOR:INFINEON,DPS310,7000
HFALGALTGPS:GEO
HFALPALTPRESSURE:ISA
LMMMGPSPERIOD1000MSEC
I023638FXA3940SIU
B1152214548857N00614809EA012230150000308
B1152224548857N00614807EA012220150000308
B1152234548857N00614806EA012220150000308
`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.igc")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(igcContent); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Parse the file
	flight, err := ParseIGCFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to parse IGC file: %v", err)
	}

	// Verify parsed data
	if flight.Pilot != "TestPilot" {
		t.Errorf("expected pilot 'TestPilot', got '%s'", flight.Pilot)
	}

	if flight.GliderType != "ACME Glider" {
		t.Errorf("expected glider type 'ACME Glider', got '%s'", flight.GliderType)
	}

	if flight.GliderID != "ABC123" {
		t.Errorf("expected glider ID 'ABC123', got '%s'", flight.GliderID)
	}

	if flight.CompetitionID != "COM123" {
		t.Errorf("expected competition ID 'COM123', got '%s'", flight.CompetitionID)
	}

	// Check date parsing
	expectedDate := time.Date(2023, 7, 30, 0, 0, 0, 0, time.UTC)
	if !flight.Date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, flight.Date)
	}

	// Check fixes
	if len(flight.Fixes) != 3 {
		t.Errorf("expected 3 fixes, got %d", len(flight.Fixes))
	}

	if len(flight.Fixes) > 0 {
		firstFix := flight.Fixes[0]
		if firstFix.AltWGS84 != 1500 {
			t.Errorf("expected first fix altitude 1500, got %f", firstFix.AltWGS84)
		}
	}
}
