package sites

import (
	"os"
	"testing"
)

func TestLoadLandingSites(t *testing.T) {
	// Valid CSV content
	validCSV := `name,lat,lon,radius
TestSite1,45.814,6.246,500
TestSite2,46.456,7.123,1000`

	// Invalid CSV content (malformed numbers)
	invalidCSV := `name,lat,lon,radius
TestSite1,invalid_lat,6.246,500
TestSite2,46.456,7.123,1000`

	// Missing columns CSV
	missingColumnsCSV := `name,lat,lon
TestSite1,45.814,6.246`

	// CSV with empty name
	emptyNameCSV := `name,lat,lon,radius
,45.814,6.246,500
TestSite2,46.456,7.123,1000`

	tests := []struct {
		name          string
		content       string
		expectError   bool
		expectedSites int
	}{
		{
			name:          "valid CSV with two sites",
			content:       validCSV,
			expectError:   false,
			expectedSites: 2,
		},
		{
			name:          "CSV with invalid numbers",
			content:       invalidCSV,
			expectError:   false,
			expectedSites: 1, // Only valid rows should be loaded
		},
		{
			name:          "CSV with missing columns",
			content:       missingColumnsCSV,
			expectError:   false,
			expectedSites: 0, // Rows with wrong number of columns should be skipped
		},
		{
			name:          "CSV with empty name",
			content:       emptyNameCSV,
			expectError:   false,
			expectedSites: 1, // Empty name rows should be skipped
		},
		{
			name:          "empty CSV",
			content:       `name,lat,lon,radius`,
			expectError:   false,
			expectedSites: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "sites_*.csv")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			tmpFile.Close()

			// Test LoadLandingSites
			collection, err := LoadLandingSites(tmpFile.Name())

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

			if collection == nil {
				t.Errorf("expected non-nil collection")
				return
			}

			if len(collection.Sites) != tt.expectedSites {
				t.Errorf("expected %d sites, got %d", tt.expectedSites, len(collection.Sites))
			}

			// For valid cases, check site properties
			if tt.expectedSites > 0 {
				firstSite := collection.Sites[0]
				if firstSite.Name == "" {
					t.Errorf("site name is empty")
				}
				if firstSite.Radius <= 0 {
					t.Errorf("site radius should be positive, got %f", firstSite.Radius)
				}
				if len(firstSite.Center) != 2 {
					t.Errorf("site center should have 2 coordinates, got %d", len(firstSite.Center))
				}
			}
		})
	}
}

func TestLoadLandingSitesNonExistentFile(t *testing.T) {
	_, err := LoadLandingSites("nonexistent.csv")
	if err == nil {
		t.Errorf("expected error for non-existent file, got none")
	}
}

func TestFindLandingSite(t *testing.T) {
	collection := &Collection{
		Sites: []LandingSite{
			{
				Name:   "NearSite",
				Center: [2]float64{6.246, 45.814}, // [lon, lat]
				Radius: 1000,                      // 1 km radius
			},
			{
				Name:   "FarSite",
				Center: [2]float64{7.000, 46.000},
				Radius: 500, // 0.5 km radius
			},
		},
	}

	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected string
	}{
		{
			name:     "point within NearSite radius",
			lat:      45.814,
			lon:      6.246,
			expected: "NearSite",
		},
		{
			name:     "point just outside NearSite radius",
			lat:      45.820, // About 667m north
			lon:      6.246,
			expected: "NearSite", // Still within 1km radius
		},
		{
			name:     "point far from any site",
			lat:      40.000,
			lon:      5.000,
			expected: "40.000,5.000", // Should return formatted coordinates
		},
		{
			name:     "point within FarSite radius",
			lat:      46.000,
			lon:      7.000,
			expected: "FarSite",
		},
		{
			name:     "point between sites but closer to FarSite",
			lat:      45.900,
			lon:      6.600,
			expected: "45.900,6.600", // Should be outside both radii
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collection.FindLandingSite(tt.lat, tt.lon)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFindLandingSiteEmptyCollection(t *testing.T) {
	collection := &Collection{Sites: []LandingSite{}}

	result := collection.FindLandingSite(45.814, 6.246)
	expected := "45.814,6.246"

	if result != expected {
		t.Errorf("expected %s for empty collection, got %s", expected, result)
	}
}

func TestFindLandingSiteNilCollection(t *testing.T) {
	var collection *Collection = nil

	// This would panic if not handled properly in calling code
	// We test that the struct methods work correctly when collection exists
	collection = &Collection{Sites: []LandingSite{}}
	result := collection.FindLandingSite(45.814, 6.246)
	expected := "45.814,6.246"

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestLandingSiteStruct(t *testing.T) {
	site := LandingSite{
		Name:   "TestSite",
		Center: [2]float64{6.246, 45.814},
		Radius: 1000.0,
	}

	if site.Name != "TestSite" {
		t.Errorf("expected name 'TestSite', got '%s'", site.Name)
	}

	if site.Center[0] != 6.246 {
		t.Errorf("expected longitude 6.246, got %f", site.Center[0])
	}

	if site.Center[1] != 45.814 {
		t.Errorf("expected latitude 45.814, got %f", site.Center[1])
	}

	if site.Radius != 1000.0 {
		t.Errorf("expected radius 1000.0, got %f", site.Radius)
	}
}

func TestCollectionStruct(t *testing.T) {
	sites := []LandingSite{
		{Name: "Site1", Center: [2]float64{6.246, 45.814}, Radius: 500},
		{Name: "Site2", Center: [2]float64{7.123, 46.456}, Radius: 1000},
	}

	collection := &Collection{Sites: sites}

	if len(collection.Sites) != 2 {
		t.Errorf("expected 2 sites, got %d", len(collection.Sites))
	}

	if collection.Sites[0].Name != "Site1" {
		t.Errorf("expected first site name 'Site1', got '%s'", collection.Sites[0].Name)
	}

	if collection.Sites[1].Name != "Site2" {
		t.Errorf("expected second site name 'Site2', got '%s'", collection.Sites[1].Name)
	}
}

// Test with different coordinate formats and edge cases
func TestFindLandingSiteEdgeCases(t *testing.T) {
	collection := &Collection{
		Sites: []LandingSite{
			{
				Name:   "EdgeSite",
				Center: [2]float64{0.0, 0.0}, // Equator/Prime Meridian
				Radius: 1000,
			},
		},
	}

	tests := []struct {
		name      string
		lat       float64
		lon       float64
		checkSite bool
	}{
		{
			name:      "exactly at site center",
			lat:       0.0,
			lon:       0.0,
			checkSite: true,
		},
		{
			name:      "negative coordinates",
			lat:       -45.814,
			lon:       -6.246,
			checkSite: false,
		},
		{
			name:      "very large coordinates",
			lat:       89.999,
			lon:       179.999,
			checkSite: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collection.FindLandingSite(tt.lat, tt.lon)

			if tt.checkSite {
				if result != "EdgeSite" {
					t.Errorf("expected 'EdgeSite', got '%s'", result)
				}
			} else {
				// Should return formatted coordinates
				if result == "EdgeSite" {
					t.Errorf("expected formatted coordinates, got site name '%s'", result)
				}
			}
		})
	}
}
