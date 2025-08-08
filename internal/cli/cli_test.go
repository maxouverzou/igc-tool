package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"igc-tool/internal/logbook"
)

func TestFindIGCFiles(t *testing.T) {
	// Create temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "igc_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"flight1.igc",
		"flight2.IGC", // Test case insensitive
		"not_igc.txt",
		"subdir/flight3.igc",
		"subdir/flight4.igc",
		"subdir/nested/flight5.igc",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name          string
		paths         []string
		recursive     bool
		expectedCount int
		expectError   bool
		expectedFiles []string
	}{
		{
			name:          "single IGC file",
			paths:         []string{filepath.Join(tmpDir, "flight1.igc")},
			recursive:     false,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "non-IGC file",
			paths:         []string{filepath.Join(tmpDir, "not_igc.txt")},
			recursive:     false,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "directory non-recursive",
			paths:         []string{tmpDir},
			recursive:     false,
			expectedCount: 2, // flight1.igc and flight2.IGC
			expectError:   false,
		},
		{
			name:          "directory recursive",
			paths:         []string{tmpDir},
			recursive:     true,
			expectedCount: 5, // all .igc files including subdirectories
			expectError:   false,
		},
		{
			name:          "non-existent file",
			paths:         []string{filepath.Join(tmpDir, "nonexistent.igc")},
			recursive:     false,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "mixed paths",
			paths:         []string{filepath.Join(tmpDir, "flight1.igc"), filepath.Join(tmpDir, "subdir")},
			recursive:     false,
			expectedCount: 3, // flight1.igc + 2 files in subdir
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindIGCFiles(tt.paths, tt.recursive)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d files, got %d: %v", tt.expectedCount, len(result), result)
			}

			// Check that all returned files have .igc extension
			for _, file := range result {
				ext := strings.ToLower(filepath.Ext(file))
				if ext != ".igc" {
					t.Errorf("expected .igc extension, got %s for file %s", ext, file)
				}
			}
		})
	}
}

func TestLoadLandingSitesIfSpecified(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectSites bool
		expectError bool
		createFile  bool
		fileContent string
	}{
		{
			name:        "empty filename",
			filename:    "",
			expectSites: false,
			expectError: false,
		},
		{
			name:        "non-existent file",
			filename:    "nonexistent.csv",
			expectSites: false,
			expectError: false, // Should not error, just return nil and print warning
		},
		{
			name:        "valid file",
			filename:    "test.csv",
			expectSites: true,
			expectError: false,
			createFile:  true,
			fileContent: "name,lat,lon,radius\nforclaz,45.814,6.246,200",
		},
		{
			name:        "invalid file",
			filename:    "invalid.csv",
			expectSites: false,
			expectError: false, // Should not error, just return nil and print warning
			createFile:  true,
			fileContent: "\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFile *os.File
			var err error

			if tt.createFile {
				tmpFile, err = os.CreateTemp("", "sites_*.geojson")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.WriteString(tt.fileContent); err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}
				tmpFile.Close()

				tt.filename = tmpFile.Name()
			}

			sites, err := LoadLandingSitesIfSpecified(tt.filename)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectSites && sites == nil {
				t.Errorf("expected sites collection but got nil")
			}

			if !tt.expectSites && sites != nil {
				t.Errorf("expected nil sites collection but got %v", sites)
			}
		})
	}
}

func TestPrintTemplatedLogbookData(t *testing.T) {
	tests := []struct {
		name        string
		data        *logbook.TemplateData
		templateStr string
		expectError bool
		expectedOut string
	}{
		{
			name:        "nil data",
			data:        nil,
			templateStr: "{{.TotalFlights}}",
			expectError: false,
			expectedOut: "No flight data available for logbook entry\n",
		},
		{
			name: "valid template with single flight",
			data: &logbook.TemplateData{
				TotalFlights: 1,
				TotalTime:    "2h30m",
				FirstDate:    "2023-06-15",
				LastDate:     "2023-06-15",
				Flights: []*logbook.Data{
					{
						Date:           "2023-06-15",
						Pilot:          "Test Pilot",
						GliderType:     "Test Glider",
						FlightDuration: "2h30m",
						MaxAltitude:    1500,
						AltitudeUnit:   "m",
					},
				},
			},
			templateStr: "Total: {{.TotalFlights}} flights, {{.TotalTime}}",
			expectError: false,
			expectedOut: "Total: 1 flights, 2h30m",
		},
		{
			name: "template with flight details",
			data: &logbook.TemplateData{
				TotalFlights: 1,
				Flights: []*logbook.Data{
					{
						Date:           "2023-06-15",
						Pilot:          "John Doe",
						GliderType:     "Paraglider Alpha",
						FlightDuration: "1h45m",
						MaxAltitude:    1200,
						TakeoffSite:    "Mountain Peak",
						LandingSite:    "Valley Floor",
						AltitudeUnit:   "m",
						SpeedUnit:      "km/h",
					},
				},
			},
			templateStr: "{{range .Flights}}Date: {{.Date}}, Pilot: {{.Pilot}}, Duration: {{.FlightDuration}}{{end}}",
			expectError: false,
			expectedOut: "Date: 2023-06-15, Pilot: John Doe, Duration: 1h45m",
		},
		{
			name: "template with aggregated statistics",
			data: &logbook.TemplateData{
				TotalFlights:   3,
				TotalTime:      "6h15m",
				FirstDate:      "2023-06-01",
				LastDate:       "2023-06-15",
				MaxAltitude:    1800,
				AvgMaxAltitude: 1500,
				UniquePilots:   []string{"John Doe", "Jane Smith"},
				UniqueGliders:  []string{"Paraglider Alpha", "Paraglider Beta"},
				UniqueSites:    []string{"Mountain Peak", "Hill Top"},
				AltitudeUnit:   "m",
			},
			templateStr: "Flights: {{.TotalFlights}}, Time: {{.TotalTime}}, Pilots: {{len .UniquePilots}}, Max Alt: {{.MaxAltitude}}{{.AltitudeUnit}}",
			expectError: false,
			expectedOut: "Flights: 3, Time: 6h15m, Pilots: 2, Max Alt: 1800m",
		},
		{
			name: "empty template data",
			data: &logbook.TemplateData{
				TotalFlights: 0,
				Flights:      []*logbook.Data{},
			},
			templateStr: "{{.TotalFlights}} flights found",
			expectError: false,
			expectedOut: "0 flights found",
		},
		{
			name:        "invalid template syntax",
			data:        &logbook.TemplateData{TotalFlights: 1},
			templateStr: "{{.TotalFlights",
			expectError: true,
			expectedOut: "",
		},
		{
			name:        "template with non-existent field",
			data:        &logbook.TemplateData{TotalFlights: 1},
			templateStr: "{{.NonExistentField}}",
			expectError: true,
			expectedOut: "",
		},
		{
			name: "complex template with conditionals",
			data: &logbook.TemplateData{
				TotalFlights: 2,
				Flights: []*logbook.Data{
					{Date: "2023-06-01", Pilot: "John"},
					{Date: "2023-06-02", Pilot: "Jane"},
				},
			},
			templateStr: "{{if gt .TotalFlights 1}}Multiple flights: {{.TotalFlights}}{{else}}Single flight{{end}}",
			expectError: false,
			expectedOut: "Multiple flights: 2",
		},
		{
			name: "template with range and index",
			data: &logbook.TemplateData{
				Flights: []*logbook.Data{
					{Date: "2023-06-01", Pilot: "Alice"},
					{Date: "2023-06-02", Pilot: "Bob"},
				},
			},
			templateStr: "{{range $i, $flight := .Flights}}{{$i}}: {{$flight.Pilot}} {{end}}",
			expectError: false,
			expectedOut: "0: Alice 1: Bob ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute the function
			err := PrintTemplatedLogbookData(tt.data, tt.templateStr)

			// Restore stdout and get output
			w.Close()
			os.Stdout = oldStdout

			output := make([]byte, 1024)
			n, _ := r.Read(output)
			actualOut := string(output[:n])

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check output
			if actualOut != tt.expectedOut {
				t.Errorf("expected output %q, got %q", tt.expectedOut, actualOut)
			}
		})
	}
}

// Helper function to test template execution without capturing output
func testTemplateExecution(data *logbook.Data, templateStr string) error {
	if data == nil {
		return nil
	}

	tmpl, err := template.New("test").Parse(templateStr)
	if err != nil {
		return err
	}

	// Execute to a strings.Builder instead of stdout for testing
	var builder strings.Builder
	return tmpl.Execute(&builder, data)
}
