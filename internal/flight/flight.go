package flight

import (
	"math"
	"time"

	"github.com/twpayne/go-igc"
)

// Constants for calculations
const (
	EarthRadiusMeters  = 6371000 // Earth radius in meters
	DegreesToRadians   = math.Pi / 180
	MinTimeDiffSeconds = 1 // minimum time difference for speed calculations
)

// Flight represents parsed IGC flight data
type Flight struct {
	Date               time.Time
	Pilot              string
	Crew               string
	GliderType         string
	GliderID           string
	CompetitionID      string
	GPSDatum           string
	FirmwareVersion    string
	HardwareVersion    string
	FlightRecorderType string
	GPSReceiver        string
	TimeZone           string
	PressureAltSensor  string
	AltGPSRef          string
	AltPressureRef     string
	Fixes              []*igc.BRecord
}

// Statistics holds calculated flight statistics
type Statistics struct {
	MaxAltitude    int
	MinAltitude    int
	MaxGroundSpeed float64
	MaxClimbRate   float64
	MaxDescentRate float64
	FlightDuration time.Duration
}

// CalculateMaxAltitude finds the maximum GPS altitude in the flight
func (f *Flight) CalculateMaxAltitude() int {
	if len(f.Fixes) == 0 {
		return 0
	}

	maxAlt := int(f.Fixes[0].AltWGS84)
	for _, fix := range f.Fixes {
		if int(fix.AltWGS84) > maxAlt {
			maxAlt = int(fix.AltWGS84)
		}
	}
	return maxAlt
}

// CalculateMinAltitude finds the minimum GPS altitude in the flight
func (f *Flight) CalculateMinAltitude() int {
	if len(f.Fixes) == 0 {
		return 0
	}

	minAlt := int(f.Fixes[0].AltWGS84)
	for _, fix := range f.Fixes {
		if int(fix.AltWGS84) < minAlt {
			minAlt = int(fix.AltWGS84)
		}
	}
	return minAlt
}

// CalculateMaxGroundSpeed finds the maximum ground speed in km/h during the flight
func (f *Flight) CalculateMaxGroundSpeed(minTimeWindowSeconds float64) float64 {
	if len(f.Fixes) < 2 {
		return 0
	}

	maxSpeed := 0.0

	for i := 1; i < len(f.Fixes); i++ {
		prev := f.Fixes[i-1]
		curr := f.Fixes[i]

		distance := HaversineDistance(prev.Lat, prev.Lon, curr.Lat, curr.Lon)
		timeDiff := curr.Time.Sub(prev.Time).Seconds()

		if timeDiff < MinTimeDiffSeconds {
			continue
		}

		speedMS := distance / timeDiff
		speedKMH := speedMS * 3.6 // Convert m/s to km/h

		// Apply windowing for GPS noise reduction
		if timeDiff < minTimeWindowSeconds && i >= 5 {
			for j := i - 1; j >= 0; j-- {
				prevWindow := f.Fixes[j]
				windowTimeDiff := curr.Time.Sub(prevWindow.Time).Seconds()

				if windowTimeDiff >= minTimeWindowSeconds {
					windowDistance := HaversineDistance(prevWindow.Lat, prevWindow.Lon, curr.Lat, curr.Lon)
					windowSpeedMS := windowDistance / windowTimeDiff
					windowSpeedKMH := windowSpeedMS * 3.6

					if windowSpeedKMH < speedKMH {
						speedKMH = windowSpeedKMH
					}
					break
				}
			}
		}

		if speedKMH > maxSpeed {
			maxSpeed = speedKMH
		}
	}
	return maxSpeed
}

// CalculateVerticalSpeeds finds the maximum and minimum vertical speeds in m/s
func (f *Flight) CalculateVerticalSpeeds() (float64, float64) {
	if len(f.Fixes) < 2 {
		return 0, 0
	}

	maxVerticalSpeed := 0.0
	minVerticalSpeed := 0.0

	for i := 1; i < len(f.Fixes); i++ {
		prev := f.Fixes[i-1]
		curr := f.Fixes[i]

		altDiff := float64(curr.AltWGS84 - prev.AltWGS84)
		timeDiff := curr.Time.Sub(prev.Time).Seconds()

		if timeDiff < MinTimeDiffSeconds {
			continue
		}

		verticalSpeed := altDiff / timeDiff

		if verticalSpeed > maxVerticalSpeed {
			maxVerticalSpeed = verticalSpeed
		}
		if verticalSpeed < minVerticalSpeed {
			minVerticalSpeed = verticalSpeed
		}
	}

	return maxVerticalSpeed, minVerticalSpeed
}

// GetStatistics calculates all flight statistics
func (f *Flight) GetStatistics(speedWindow float64) *Statistics {
	maxClimbRate, minVerticalSpeed := f.CalculateVerticalSpeeds()

	var duration time.Duration
	if len(f.Fixes) >= 2 {
		duration = f.Fixes[len(f.Fixes)-1].Time.Sub(f.Fixes[0].Time)
	}

	return &Statistics{
		MaxAltitude:    f.CalculateMaxAltitude(),
		MinAltitude:    f.CalculateMinAltitude(),
		MaxGroundSpeed: f.CalculateMaxGroundSpeed(speedWindow),
		MaxClimbRate:   maxClimbRate,
		MaxDescentRate: math.Abs(minVerticalSpeed),
		FlightDuration: duration,
	}
}

// HaversineDistance calculates the distance between two points in meters
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * DegreesToRadians
	lon1Rad := lon1 * DegreesToRadians
	lat2Rad := lat2 * DegreesToRadians
	lon2Rad := lon2 * DegreesToRadians

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusMeters * c
}
