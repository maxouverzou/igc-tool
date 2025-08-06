package utils

import (
	"fmt"
	"time"
)

// FormatTime formats time according to the specified format
func FormatTime(t time.Time, format string) string {
	switch format {
	case "ampm":
		return t.Format("3:04:05 PM")
	default: // 24h
		return t.Format("15:04:05")
	}
}

// FormatDuration formats a duration as "XhYm"
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// FormatCoordinates formats lat/lon as a string
func FormatCoordinates(lat, lon float64) string {
	return fmt.Sprintf("%.3f,%.3f", lat, lon)
}
