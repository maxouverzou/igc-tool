package flight

import (
	"math"
	"testing"
	"time"

	"github.com/twpayne/go-igc"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name           string
		lat1, lon1     float64
		lat2, lon2     float64
		expectedMeters float64
		tolerance      float64
	}{
		{
			name:           "same point",
			lat1:           45.814,
			lon1:           6.246,
			lat2:           45.814,
			lon2:           6.246,
			expectedMeters: 0,
			tolerance:      1,
		},
		{
			name:           "short distance",
			lat1:           45.814,
			lon1:           6.246,
			lat2:           45.815,
			lon2:           6.247,
			expectedMeters: 141, // approximately 141 meters
			tolerance:      10,
		},
		{
			name:           "longer distance",
			lat1:           45.814,
			lon1:           6.246,
			lat2:           46.814,
			lon2:           7.246,
			expectedMeters: 134829, // approximately 134.8 km
			tolerance:      1000,
		},
		{
			name:           "across equator",
			lat1:           1.0,
			lon1:           0.0,
			lat2:           -1.0,
			lon2:           0.0,
			expectedMeters: 222390, // approximately 222 km
			tolerance:      1000,
		},
		{
			name:           "antipodal points (approximately)",
			lat1:           0.0,
			lon1:           0.0,
			lat2:           0.0,
			lon2:           180.0,
			expectedMeters: 20015087, // More accurate expected value
			tolerance:      10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(result-tt.expectedMeters) > tt.tolerance {
				t.Errorf("expected distance %f ± %f, got %f", tt.expectedMeters, tt.tolerance, result)
			}
		})
	}
}

func TestFlightCalculateMaxAltitude(t *testing.T) {
	tests := []struct {
		name     string
		fixes    []*igc.BRecord
		expected int
	}{
		{
			name:     "empty fixes",
			fixes:    []*igc.BRecord{},
			expected: 0,
		},
		{
			name: "single fix",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500},
			},
			expected: 1500,
		},
		{
			name: "multiple fixes",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500},
				{AltWGS84: 2000},
				{AltWGS84: 1800},
				{AltWGS84: 2200},
				{AltWGS84: 1900},
			},
			expected: 2200,
		},
		{
			name: "negative altitudes",
			fixes: []*igc.BRecord{
				{AltWGS84: -100},
				{AltWGS84: -50},
				{AltWGS84: -200},
			},
			expected: -50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flight := &Flight{Fixes: tt.fixes}
			result := flight.CalculateMaxAltitude()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestFlightCalculateMinAltitude(t *testing.T) {
	tests := []struct {
		name     string
		fixes    []*igc.BRecord
		expected int
	}{
		{
			name:     "empty fixes",
			fixes:    []*igc.BRecord{},
			expected: 0,
		},
		{
			name: "single fix",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500},
			},
			expected: 1500,
		},
		{
			name: "multiple fixes",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500},
				{AltWGS84: 2000},
				{AltWGS84: 1800},
				{AltWGS84: 1200},
				{AltWGS84: 1900},
			},
			expected: 1200,
		},
		{
			name: "negative altitudes",
			fixes: []*igc.BRecord{
				{AltWGS84: -100},
				{AltWGS84: -50},
				{AltWGS84: -200},
			},
			expected: -200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flight := &Flight{Fixes: tt.fixes}
			result := flight.CalculateMinAltitude()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestFlightCalculateMaxGroundSpeed(t *testing.T) {
	baseTime := time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		fixes     []*igc.BRecord
		window    float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "empty fixes",
			fixes:     []*igc.BRecord{},
			window:    5.0,
			expected:  0,
			tolerance: 0,
		},
		{
			name: "single fix",
			fixes: []*igc.BRecord{
				{Lat: 45.814, Lon: 6.246, Time: baseTime},
			},
			window:    5.0,
			expected:  0,
			tolerance: 0,
		},
		{
			name: "two fixes with movement",
			fixes: []*igc.BRecord{
				{Lat: 45.814, Lon: 6.246, Time: baseTime},
				{Lat: 45.815, Lon: 6.247, Time: baseTime.Add(10 * time.Second)},
			},
			window:    5.0,
			expected:  50.76, // approximately 50 km/h
			tolerance: 5.0,
		},
		{
			name: "multiple fixes varying speeds",
			fixes: []*igc.BRecord{
				{Lat: 45.814, Lon: 6.246, Time: baseTime},
				{Lat: 45.815, Lon: 6.247, Time: baseTime.Add(10 * time.Second)}, // ~50 km/h
				{Lat: 45.816, Lon: 6.248, Time: baseTime.Add(15 * time.Second)}, // faster
				{Lat: 45.817, Lon: 6.249, Time: baseTime.Add(25 * time.Second)}, // slower
			},
			window:    5.0,
			expected:  100.0, // should find the maximum
			tolerance: 20.0,
		},
		{
			name: "fixes too close in time",
			fixes: []*igc.BRecord{
				{Lat: 45.814, Lon: 6.246, Time: baseTime},
				{Lat: 45.820, Lon: 6.250, Time: baseTime.Add(500 * time.Millisecond)}, // < 1 second
			},
			window:    5.0,
			expected:  0,
			tolerance: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flight := &Flight{Fixes: tt.fixes}
			result := flight.CalculateMaxGroundSpeed(tt.window)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected speed %f ± %f km/h, got %f km/h", tt.expected, tt.tolerance, result)
			}
		})
	}
}

func TestFlightCalculateVerticalSpeeds(t *testing.T) {
	baseTime := time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		fixes              []*igc.BRecord
		expectedMaxClimb   float64
		expectedMaxDescent float64
		tolerance          float64
	}{
		{
			name:               "empty fixes",
			fixes:              []*igc.BRecord{},
			expectedMaxClimb:   0,
			expectedMaxDescent: 0,
			tolerance:          0,
		},
		{
			name: "single fix",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500, Time: baseTime},
			},
			expectedMaxClimb:   0,
			expectedMaxDescent: 0,
			tolerance:          0,
		},
		{
			name: "steady climb",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500, Time: baseTime},
				{AltWGS84: 1600, Time: baseTime.Add(10 * time.Second)}, // +10 m/s
				{AltWGS84: 1700, Time: baseTime.Add(20 * time.Second)}, // +10 m/s
			},
			expectedMaxClimb:   10.0,
			expectedMaxDescent: 0.0,
			tolerance:          0.1,
		},
		{
			name: "steady descent",
			fixes: []*igc.BRecord{
				{AltWGS84: 1700, Time: baseTime},
				{AltWGS84: 1600, Time: baseTime.Add(10 * time.Second)}, // -10 m/s
				{AltWGS84: 1500, Time: baseTime.Add(20 * time.Second)}, // -10 m/s
			},
			expectedMaxClimb:   0.0,
			expectedMaxDescent: 10.0, // Should be positive (absolute value)
			tolerance:          0.1,
		},
		{
			name: "mixed climb and descent",
			fixes: []*igc.BRecord{
				{AltWGS84: 1500, Time: baseTime},
				{AltWGS84: 1650, Time: baseTime.Add(10 * time.Second)}, // +15 m/s
				{AltWGS84: 1600, Time: baseTime.Add(15 * time.Second)}, // -10 m/s
				{AltWGS84: 1500, Time: baseTime.Add(20 * time.Second)}, // -20 m/s
			},
			expectedMaxClimb:   15.0,
			expectedMaxDescent: 20.0,
			tolerance:          0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flight := &Flight{Fixes: tt.fixes}
			maxClimb, minVerticalSpeed := flight.CalculateVerticalSpeeds()

			// maxDescent should be the absolute value of the minimum vertical speed
			maxDescent := math.Abs(minVerticalSpeed)

			if math.Abs(maxClimb-tt.expectedMaxClimb) > tt.tolerance {
				t.Errorf("expected max climb %f ± %f, got %f", tt.expectedMaxClimb, tt.tolerance, maxClimb)
			}
			if math.Abs(maxDescent-tt.expectedMaxDescent) > tt.tolerance {
				t.Errorf("expected max descent %f ± %f, got %f", tt.expectedMaxDescent, tt.tolerance, maxDescent)
			}
		})
	}
}

func TestFlightGetStatistics(t *testing.T) {
	baseTime := time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC)

	flight := &Flight{
		Fixes: []*igc.BRecord{
			{AltWGS84: 1500, Time: baseTime, Lat: 45.814, Lon: 6.246},
			{AltWGS84: 1600, Time: baseTime.Add(10 * time.Second), Lat: 45.815, Lon: 6.247},
			{AltWGS84: 1800, Time: baseTime.Add(20 * time.Second), Lat: 45.816, Lon: 6.248},
			{AltWGS84: 1200, Time: baseTime.Add(30 * time.Second), Lat: 45.817, Lon: 6.249},
		},
	}

	stats := flight.GetStatistics(5.0)

	if stats == nil {
		t.Fatal("expected non-nil statistics")
	}

	// Check that all fields are populated
	if stats.MaxAltitude != 1800 {
		t.Errorf("expected max altitude 1800, got %d", stats.MaxAltitude)
	}

	if stats.MinAltitude != 1200 {
		t.Errorf("expected min altitude 1200, got %d", stats.MinAltitude)
	}

	if stats.FlightDuration != 30*time.Second {
		t.Errorf("expected flight duration 30s, got %v", stats.FlightDuration)
	}

	// Speed and vertical speeds should be positive
	if stats.MaxGroundSpeed < 0 {
		t.Errorf("expected positive ground speed, got %f", stats.MaxGroundSpeed)
	}

	if stats.MaxClimbRate < 0 {
		t.Errorf("expected positive climb rate, got %f", stats.MaxClimbRate)
	}

	if stats.MaxDescentRate < 0 {
		t.Errorf("expected positive descent rate, got %f", stats.MaxDescentRate)
	}
}

func TestFlightEmptyFixes(t *testing.T) {
	flight := &Flight{Fixes: []*igc.BRecord{}}

	// Test all methods with empty fixes
	if maxAlt := flight.CalculateMaxAltitude(); maxAlt != 0 {
		t.Errorf("expected 0 max altitude for empty fixes, got %d", maxAlt)
	}

	if minAlt := flight.CalculateMinAltitude(); minAlt != 0 {
		t.Errorf("expected 0 min altitude for empty fixes, got %d", minAlt)
	}

	if speed := flight.CalculateMaxGroundSpeed(5.0); speed != 0 {
		t.Errorf("expected 0 speed for empty fixes, got %f", speed)
	}

	maxClimb, maxDescent := flight.CalculateVerticalSpeeds()
	if maxClimb != 0 || maxDescent != 0 {
		t.Errorf("expected 0 vertical speeds for empty fixes, got climb=%f, descent=%f", maxClimb, maxDescent)
	}

	stats := flight.GetStatistics(5.0)
	if stats.FlightDuration != 0 {
		t.Errorf("expected 0 duration for empty fixes, got %v", stats.FlightDuration)
	}
}
