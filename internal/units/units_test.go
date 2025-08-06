package units

import (
	"math"
	"testing"
)

func TestAltitude(t *testing.T) {
	tests := []struct {
		name      string
		meters    float64
		unit      string
		expected  float64
		tolerance float64
	}{
		{
			name:      "meters to meters",
			meters:    1000,
			unit:      "m",
			expected:  1000,
			tolerance: 0.01,
		},
		{
			name:      "meters to feet",
			meters:    1000,
			unit:      "ft",
			expected:  3280.84,
			tolerance: 0.01,
		},
		{
			name:      "zero altitude",
			meters:    0,
			unit:      "ft",
			expected:  0,
			tolerance: 0.01,
		},
		{
			name:      "negative altitude",
			meters:    -100,
			unit:      "ft",
			expected:  -328.084,
			tolerance: 0.01,
		},
		{
			name:      "unknown unit defaults to meters",
			meters:    500,
			unit:      "unknown",
			expected:  500,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Altitude(tt.meters, tt.unit)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestSpeed(t *testing.T) {
	tests := []struct {
		name      string
		kmh       float64
		unit      string
		expected  float64
		tolerance float64
	}{
		{
			name:      "kmh to kmh",
			kmh:       100,
			unit:      "kmh",
			expected:  100,
			tolerance: 0.01,
		},
		{
			name:      "kmh to mph",
			kmh:       100,
			unit:      "mph",
			expected:  62.1371,
			tolerance: 0.01,
		},
		{
			name:      "kmh to knots",
			kmh:       100,
			unit:      "kts",
			expected:  53.9957,
			tolerance: 0.01,
		},
		{
			name:      "kmh to m/s",
			kmh:       36,
			unit:      "ms",
			expected:  10,
			tolerance: 0.01,
		},
		{
			name:      "zero speed",
			kmh:       0,
			unit:      "mph",
			expected:  0,
			tolerance: 0.01,
		},
		{
			name:      "unknown unit defaults to kmh",
			kmh:       50,
			unit:      "unknown",
			expected:  50,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Speed(tt.kmh, tt.unit)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestClimb(t *testing.T) {
	tests := []struct {
		name      string
		ms        float64
		unit      string
		expected  float64
		tolerance float64
	}{
		{
			name:      "m/s to m/s",
			ms:        5,
			unit:      "ms",
			expected:  5,
			tolerance: 0.01,
		},
		{
			name:      "m/s to ft/min",
			ms:        1,
			unit:      "fpm",
			expected:  196.85, // 1 * 3.28084 * 60
			tolerance: 0.01,
		},
		{
			name:      "zero climb rate",
			ms:        0,
			unit:      "fpm",
			expected:  0,
			tolerance: 0.01,
		},
		{
			name:      "negative climb rate (descent)",
			ms:        -2,
			unit:      "fpm",
			expected:  -393.7, // -2 * 3.28084 * 60
			tolerance: 0.1,
		},
		{
			name:      "unknown unit defaults to m/s",
			ms:        3,
			unit:      "unknown",
			expected:  3,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Climb(tt.ms, tt.unit)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestAltitudeSymbol(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "meters unit",
			unit:     "m",
			expected: "m",
		},
		{
			name:     "feet unit",
			unit:     "ft",
			expected: "ft",
		},
		{
			name:     "unknown unit defaults to meters",
			unit:     "unknown",
			expected: "m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AltitudeSymbol(tt.unit)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSpeedSymbol(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "kmh unit",
			unit:     "kmh",
			expected: "km/h",
		},
		{
			name:     "mph unit",
			unit:     "mph",
			expected: "mph",
		},
		{
			name:     "knots unit",
			unit:     "kts",
			expected: "kts",
		},
		{
			name:     "m/s unit",
			unit:     "ms",
			expected: "m/s",
		},
		{
			name:     "unknown unit defaults to km/h",
			unit:     "unknown",
			expected: "km/h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SpeedSymbol(tt.unit)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestClimbSymbol(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "m/s unit",
			unit:     "ms",
			expected: "m/s",
		},
		{
			name:     "ft/min unit",
			unit:     "fpm",
			expected: "ft/min",
		},
		{
			name:     "unknown unit defaults to m/s",
			unit:     "unknown",
			expected: "m/s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClimbSymbol(tt.unit)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
