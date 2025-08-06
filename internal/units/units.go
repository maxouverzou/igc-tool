package units

// Unit constants
const (
	// Altitude units
	AltitudeMeters = "m"
	AltitudeFeet   = "ft"

	// Speed units
	SpeedKmh   = "kmh"
	SpeedMph   = "mph"
	SpeedKnots = "kts"
	SpeedMs    = "ms"

	// Climb units
	ClimbMs  = "ms"  // meters per second
	ClimbFpm = "fpm" // feet per minute

	// Time formats
	TimeFormat24h  = "24h"
	TimeFormatAMPM = "ampm"
)

// Unit conversion constants
const (
	MetersToFeet = 3.28084
	KmhToMph     = 0.621371
	KmhToKnots   = 0.539957
	MsToKmh      = 3.6 // meters per second to kilometers per hour
)

// Altitude converts altitude from meters to the specified unit
func Altitude(meters float64, unit string) float64 {
	switch unit {
	case AltitudeFeet:
		return meters * MetersToFeet
	default: // meters
		return meters
	}
}

// Speed converts speed from km/h to the specified unit
func Speed(kmh float64, unit string) float64 {
	switch unit {
	case SpeedMph:
		return kmh * KmhToMph
	case SpeedKnots:
		return kmh * KmhToKnots
	case SpeedMs:
		return kmh / MsToKmh
	default: // kmh
		return kmh
	}
}

// Climb converts vertical speed from m/s to the specified unit
func Climb(ms float64, unit string) float64 {
	switch unit {
	case ClimbFpm:
		return ms * MetersToFeet * 60 // m/s to ft/min
	default: // ms (meters per second)
		return ms
	}
}

// AltitudeSymbol returns the symbol for the altitude unit
func AltitudeSymbol(unit string) string {
	switch unit {
	case AltitudeFeet:
		return "ft"
	default:
		return "m"
	}
}

// SpeedSymbol returns the symbol for the speed unit
func SpeedSymbol(unit string) string {
	switch unit {
	case SpeedMph:
		return "mph"
	case SpeedKnots:
		return "kts"
	case SpeedMs:
		return "m/s"
	default:
		return "km/h"
	}
}

// ClimbSymbol returns the symbol for climb rate unit
func ClimbSymbol(unit string) string {
	switch unit {
	case ClimbFpm:
		return "ft/min"
	default:
		return "m/s"
	}
}

// ValidateAltitudeUnit checks if the given altitude unit is valid
func ValidateAltitudeUnit(unit string) bool {
	switch unit {
	case AltitudeMeters, AltitudeFeet:
		return true
	default:
		return false
	}
}

// ValidateSpeedUnit checks if the given speed unit is valid
func ValidateSpeedUnit(unit string) bool {
	switch unit {
	case SpeedKmh, SpeedMph, SpeedKnots, SpeedMs:
		return true
	default:
		return false
	}
}

// ValidateClimbUnit checks if the given climb unit is valid
func ValidateClimbUnit(unit string) bool {
	switch unit {
	case ClimbMs, ClimbFpm:
		return true
	default:
		return false
	}
}

// ValidateTimeFormat checks if the given time format is valid
func ValidateTimeFormat(format string) bool {
	switch format {
	case TimeFormat24h, TimeFormatAMPM:
		return true
	default:
		return false
	}
}
