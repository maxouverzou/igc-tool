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

func TestPrintTemplatedLogbook(t *testing.T) {
	// Test data
	testData := &logbook.Data{
		Date:              "2025-07-18",
		TakeoffSite:       "TestSite",
		TakeoffAlt:        1500,
		AltitudeDiff:      300,
		FlightDuration:    "2h30m",
		MaxAltitude:       1800,
		MaxGroundSpeed:    85,
		MaxClimbRate:      8.5,
		MaxDescentRate:    12.3,
		Pilot:             "TestPilot",
		AltitudeUnit:      "m",
		SpeedUnit:         "km/h",
		VerticalSpeedUnit: "m/s",
	}

	tests := []struct {
		name         string
		data         *logbook.Data
		template     string
		expectError  bool
		expectOutput bool
	}{
		{
			name:         "valid template",
			data:         testData,
			template:     "{{.Date}} {{.Pilot}} {{.TakeoffSite}}",
			expectError:  false,
			expectOutput: true,
		},
		{
			name:         "template with units",
			data:         testData,
			template:     "{{.MaxAltitude}}{{.AltitudeUnit}} {{.MaxGroundSpeed}}{{.SpeedUnit}}",
			expectError:  false,
			expectOutput: true,
		},
		{
			name:         "invalid template syntax",
			data:         testData,
			template:     "{{.InvalidField",
			expectError:  true,
			expectOutput: false,
		},
		{
			name:         "template with invalid field",
			data:         testData,
			template:     "{{.NonExistentField}}",
			expectError:  true,
			expectOutput: false,
		},
		{
			name:         "nil data",
			data:         nil,
			template:     "{{.Date}}",
			expectError:  false,
			expectOutput: true, // Should print message about no data
		},
		{
			name:         "empty template",
			data:         testData,
			template:     "",
			expectError:  false,
			expectOutput: false,
		},
		{
			name:         "template with newline",
			data:         testData,
			template:     "{{.Date}}\n",
			expectError:  false,
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PrintTemplatedLogbook(tt.data, tt.template)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
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
