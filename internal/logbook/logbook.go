package logbook

import (
	"math"
	"reflect"

	"igc-tool/internal/config"
	"igc-tool/internal/flight"
	"igc-tool/internal/sites"
	"igc-tool/internal/units"
	"igc-tool/internal/utils"
)

// Data represents the data structure used for logbook template rendering
type Data struct {
	Date               string
	TakeoffLat         float64
	TakeoffLon         float64
	TakeoffPosition    string
	TakeoffSite        string
	LandingLat         float64
	LandingLon         float64
	LandingPosition    string
	LandingSite        string
	TakeoffAlt         int
	LandingAlt         int
	AltitudeDiff       int
	MaxAltitude        int
	MinAltitude        int
	MaxGroundSpeed     int
	MaxClimbRate       float64
	MaxDescentRate     float64
	FlightDuration     string
	TakeoffTime        string
	LandingTime        string
	Pilot              string
	Crew               string
	GliderType         string
	GliderID           string
	CompetitionID      string
	FlightRecorderType string
	Filename           string
	// Unit symbols for formatting
	AltitudeUnit      string
	SpeedUnit         string
	VerticalSpeedUnit string // Unit for climb/descent rates
}

// Options holds configuration for creating logbook data
type Options struct {
	LandingSites *sites.Collection
	Filename     string
	SpeedWindow  float64
	AltitudeUnit string
	SpeedUnit    string
	ClimbUnit    string
	TimeFormat   string
}

// CreateData creates logbook data from a flight using the provided options
func CreateData(f *flight.Flight, opts Options) *Data {
	if len(f.Fixes) == 0 {
		return nil
	}

	takeoffFix := f.Fixes[0]
	landingFix := f.Fixes[len(f.Fixes)-1]
	duration := landingFix.Time.Sub(takeoffFix.Time)
	altitudeDiff := int(landingFix.AltWGS84) - int(takeoffFix.AltWGS84)

	// Calculate flight statistics
	stats := f.GetStatistics(opts.SpeedWindow)

	// Determine takeoff and landing sites
	takeoffSite := utils.FormatCoordinates(takeoffFix.Lat, takeoffFix.Lon)
	landingSite := utils.FormatCoordinates(landingFix.Lat, landingFix.Lon)

	if opts.LandingSites != nil {
		takeoffSite = opts.LandingSites.FindLandingSite(takeoffFix.Lat, takeoffFix.Lon)
		landingSite = opts.LandingSites.FindLandingSite(landingFix.Lat, landingFix.Lon)
	}

	// Apply unit conversions
	takeoffAltConverted := int(units.Altitude(float64(takeoffFix.AltWGS84), opts.AltitudeUnit))
	landingAltConverted := int(units.Altitude(float64(landingFix.AltWGS84), opts.AltitudeUnit))
	altitudeDiffConverted := int(units.Altitude(float64(altitudeDiff), opts.AltitudeUnit))
	maxAltitudeConverted := int(units.Altitude(float64(stats.MaxAltitude), opts.AltitudeUnit))
	minAltitudeConverted := int(units.Altitude(float64(stats.MinAltitude), opts.AltitudeUnit))
	maxGroundSpeedConverted := int(math.Round(units.Speed(stats.MaxGroundSpeed, opts.SpeedUnit)))
	maxClimbRateConverted := math.Round(units.Climb(stats.MaxClimbRate, opts.ClimbUnit))
	maxDescentRateConverted := math.Round(units.Climb(stats.MaxDescentRate, opts.ClimbUnit))

	return &Data{
		Date:               f.Date.Format("2006-01-02"),
		TakeoffLat:         takeoffFix.Lat,
		TakeoffLon:         takeoffFix.Lon,
		TakeoffPosition:    utils.FormatCoordinates(takeoffFix.Lat, takeoffFix.Lon),
		TakeoffSite:        takeoffSite,
		LandingLat:         landingFix.Lat,
		LandingLon:         landingFix.Lon,
		LandingPosition:    utils.FormatCoordinates(landingFix.Lat, landingFix.Lon),
		LandingSite:        landingSite,
		TakeoffAlt:         takeoffAltConverted,
		LandingAlt:         landingAltConverted,
		AltitudeDiff:       altitudeDiffConverted,
		MaxAltitude:        maxAltitudeConverted,
		MinAltitude:        minAltitudeConverted,
		MaxGroundSpeed:     maxGroundSpeedConverted,
		MaxClimbRate:       maxClimbRateConverted,
		MaxDescentRate:     maxDescentRateConverted,
		FlightDuration:     utils.FormatDuration(duration),
		TakeoffTime:        utils.FormatTime(takeoffFix.Time, opts.TimeFormat),
		LandingTime:        utils.FormatTime(landingFix.Time, opts.TimeFormat),
		Pilot:              f.Pilot,
		Crew:               f.Crew,
		GliderType:         f.GliderType,
		GliderID:           f.GliderID,
		CompetitionID:      f.CompetitionID,
		FlightRecorderType: f.FlightRecorderType,
		Filename:           opts.Filename,
		AltitudeUnit:       units.AltitudeSymbol(opts.AltitudeUnit),
		SpeedUnit:          units.SpeedSymbol(opts.SpeedUnit),
		VerticalSpeedUnit:  units.ClimbSymbol(opts.ClimbUnit),
	}
}

// GetDataFields returns a list of available template fields for Data
func GetDataFields() []string {
	var fields []string
	t := reflect.TypeOf(Data{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Only include exported fields (fields that start with uppercase)
		if field.PkgPath == "" {
			fields = append(fields, field.Name)
		}
	}

	return fields
}

// CreateOptions creates Options from config
func CreateOptions(cfg *config.Config, landingSites *sites.Collection, filename string) Options {
	return Options{
		LandingSites: landingSites,
		Filename:     filename,
		SpeedWindow:  cfg.SpeedWindow,
		AltitudeUnit: cfg.AltitudeUnit,
		SpeedUnit:    cfg.SpeedUnit,
		ClimbUnit:    cfg.ClimbUnit,
		TimeFormat:   cfg.TimeFormat,
	}
}
