package utils

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	testTime := time.Date(2025, 7, 18, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		time     time.Time
		format   string
		expected string
	}{
		{
			name:     "24h format",
			time:     testTime,
			format:   "24h",
			expected: "14:30:45",
		},
		{
			name:     "ampm format",
			time:     testTime,
			format:   "ampm",
			expected: "2:30:45 PM",
		},
		{
			name:     "morning time ampm",
			time:     time.Date(2025, 7, 18, 9, 15, 30, 0, time.UTC),
			format:   "ampm",
			expected: "9:15:30 AM",
		},
		{
			name:     "midnight 24h",
			time:     time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
			format:   "24h",
			expected: "00:00:00",
		},
		{
			name:     "midnight ampm",
			time:     time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
			format:   "ampm",
			expected: "12:00:00 AM",
		},
		{
			name:     "noon 24h",
			time:     time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC),
			format:   "24h",
			expected: "12:00:00",
		},
		{
			name:     "noon ampm",
			time:     time.Date(2025, 7, 18, 12, 0, 0, 0, time.UTC),
			format:   "ampm",
			expected: "12:00:00 PM",
		},
		{
			name:     "unknown format defaults to 24h",
			time:     testTime,
			format:   "unknown",
			expected: "14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.time, tt.format)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "0h0m",
		},
		{
			name:     "one hour exactly",
			duration: time.Hour,
			expected: "1h0m",
		},
		{
			name:     "one hour thirty minutes",
			duration: time.Hour + 30*time.Minute,
			expected: "1h30m",
		},
		{
			name:     "only minutes",
			duration: 45 * time.Minute,
			expected: "0h45m",
		},
		{
			name:     "complex duration with seconds",
			duration: 2*time.Hour + 47*time.Minute + 30*time.Second,
			expected: "2h47m",
		},
		{
			name:     "multiple hours",
			duration: 5*time.Hour + 15*time.Minute,
			expected: "5h15m",
		},
		{
			name:     "duration with rounding",
			duration: 1*time.Hour + 59*time.Minute + 59*time.Second,
			expected: "1h59m",
		},
		{
			name:     "very long duration",
			duration: 25*time.Hour + 30*time.Minute,
			expected: "25h30m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatCoordinates(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected string
	}{
		{
			name:     "positive coordinates",
			lat:      45.814,
			lon:      6.246,
			expected: "45.814,6.246",
		},
		{
			name:     "negative coordinates",
			lat:      -45.814,
			lon:      -6.246,
			expected: "-45.814,-6.246",
		},
		{
			name:     "mixed coordinates",
			lat:      45.814,
			lon:      -6.246,
			expected: "45.814,-6.246",
		},
		{
			name:     "zero coordinates",
			lat:      0.0,
			lon:      0.0,
			expected: "0.000,0.000",
		},
		{
			name:     "high precision coordinates",
			lat:      45.8141592653,
			lon:      6.2467890123,
			expected: "45.814,6.247",
		},
		{
			name:     "very small coordinates",
			lat:      0.001,
			lon:      0.002,
			expected: "0.001,0.002",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCoordinates(tt.lat, tt.lon)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
