package logbook

import (
	"fmt"
	"math"
	"reflect"
	"time"

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

// TemplateData represents the complete data structure for template rendering
// including individual flights and aggregated statistics
type TemplateData struct {
	Flights        []*Data
	TotalTime      string
	TotalFlights   int
	FirstDate      string
	LastDate       string
	TotalDistance  float64
	AvgFlightTime  string
	MaxFlightTime  string
	MinFlightTime  string
	MaxAltitude    int
	AvgMaxAltitude int
	UniquePilots   []string
	UniqueGliders  []string
	UniqueSites    []string
	// Unit symbols for formatting
	AltitudeUnit      string
	SpeedUnit         string
	VerticalSpeedUnit string
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

// GetTemplateDataFields returns a list of available template fields for TemplateData
func GetTemplateDataFields() []string {
	var fields []string
	t := reflect.TypeOf(TemplateData{})

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

// CreateTemplateData creates aggregated template data from multiple flight entries
func CreateTemplateData(flights []*Data, opts Options) *TemplateData {
	if len(flights) == 0 {
		return &TemplateData{
			Flights:           []*Data{},
			TotalFlights:      0,
			AltitudeUnit:      units.AltitudeSymbol(opts.AltitudeUnit),
			SpeedUnit:         units.SpeedSymbol(opts.SpeedUnit),
			VerticalSpeedUnit: units.ClimbSymbol(opts.ClimbUnit),
		}
	}

	// Calculate aggregated statistics
	var totalDuration time.Duration
	var totalAltitude int
	var maxAltitude int
	var maxDuration time.Duration
	var minDuration time.Duration = time.Hour * 24 // Start with a large value

	pilots := make(map[string]bool)
	gliders := make(map[string]bool)
	sites := make(map[string]bool)

	var firstDate, lastDate time.Time

	for i, flight := range flights {
		// Parse flight duration
		duration, err := parseDuration(flight.FlightDuration)
		if err == nil {
			totalDuration += duration
			if duration > maxDuration {
				maxDuration = duration
			}
			if duration < minDuration {
				minDuration = duration
			}
		}

		// Track altitude statistics
		totalAltitude += flight.MaxAltitude
		if flight.MaxAltitude > maxAltitude {
			maxAltitude = flight.MaxAltitude
		}

		// Track unique values
		if flight.Pilot != "" {
			pilots[flight.Pilot] = true
		}
		if flight.GliderType != "" {
			gliders[flight.GliderType] = true
		}
		if flight.TakeoffSite != "" {
			sites[flight.TakeoffSite] = true
		}

		// Track date range
		if date, err := time.Parse("2006-01-02", flight.Date); err == nil {
			if i == 0 || date.Before(firstDate) {
				firstDate = date
			}
			if i == 0 || date.After(lastDate) {
				lastDate = date
			}
		}
	}

	// Convert maps to slices
	uniquePilots := make([]string, 0, len(pilots))
	for pilot := range pilots {
		uniquePilots = append(uniquePilots, pilot)
	}

	uniqueGliders := make([]string, 0, len(gliders))
	for glider := range gliders {
		uniqueGliders = append(uniqueGliders, glider)
	}

	uniqueSites := make([]string, 0, len(sites))
	for site := range sites {
		uniqueSites = append(uniqueSites, site)
	}

	// Calculate averages
	avgFlightTime := totalDuration / time.Duration(len(flights))
	avgMaxAltitude := totalAltitude / len(flights)

	// Handle edge cases for min duration
	if minDuration == time.Hour*24 {
		minDuration = 0
	}

	return &TemplateData{
		Flights:           flights,
		TotalTime:         utils.FormatDuration(totalDuration),
		TotalFlights:      len(flights),
		FirstDate:         firstDate.Format("2006-01-02"),
		LastDate:          lastDate.Format("2006-01-02"),
		AvgFlightTime:     utils.FormatDuration(avgFlightTime),
		MaxFlightTime:     utils.FormatDuration(maxDuration),
		MinFlightTime:     utils.FormatDuration(minDuration),
		MaxAltitude:       maxAltitude,
		AvgMaxAltitude:    avgMaxAltitude,
		UniquePilots:      uniquePilots,
		UniqueGliders:     uniqueGliders,
		UniqueSites:       uniqueSites,
		AltitudeUnit:      units.AltitudeSymbol(opts.AltitudeUnit),
		SpeedUnit:         units.SpeedSymbol(opts.SpeedUnit),
		VerticalSpeedUnit: units.ClimbSymbol(opts.ClimbUnit),
	}
}

// parseDuration parses a duration string in the format used by utils.FormatDuration
func parseDuration(durationStr string) (time.Duration, error) {
	// Handle the custom format "XhYm" used by utils.FormatDuration
	var hours, minutes int
	n, err := fmt.Sscanf(durationStr, "%dh%dm", &hours, &minutes)
	if err != nil || n != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute, nil
}
