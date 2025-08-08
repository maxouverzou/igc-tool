package renderer

import (
	"encoding/json"
	"fmt"

	"igc-tool/internal/flight"
)

// GeoJSONFeature represents a GeoJSON feature
type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONGeometry represents a GeoJSON geometry
type GeoJSONGeometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// GeoJSONFeatureCollection represents a GeoJSON feature collection
type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

// RenderToGeoJSON converts a flight track to GeoJSON format
func RenderToGeoJSON(flight *flight.Flight, pretty bool, includeMetadata bool) ([]byte, error) {
	if len(flight.Fixes) == 0 {
		return nil, fmt.Errorf("no GPS fixes found in flight data")
	}

	// Extract coordinates from B records
	var coordinates [][]float64
	for _, fix := range flight.Fixes {
		if fix.Valid() {
			// GeoJSON coordinates are [longitude, latitude, altitude]
			coord := []float64{fix.Lon, fix.Lat}
			if fix.AltWGS84 != 0 {
				coord = append(coord, fix.AltWGS84)
			}
			coordinates = append(coordinates, coord)
		}
	}

	if len(coordinates) == 0 {
		return nil, fmt.Errorf("no valid GPS fixes found in flight data")
	}

	// Create LineString geometry
	geometry := GeoJSONGeometry{
		Type:        "LineString",
		Coordinates: coordinates,
	}

	// Create properties
	properties := make(map[string]interface{})

	if includeMetadata {
		if !flight.Date.IsZero() {
			properties["date"] = flight.Date.Format("2006-01-02")
		}
		if flight.Pilot != "" {
			properties["pilot"] = flight.Pilot
		}
		if flight.GliderType != "" {
			properties["glider_type"] = flight.GliderType
		}
		if flight.GliderID != "" {
			properties["glider_id"] = flight.GliderID
		}
		if flight.CompetitionID != "" {
			properties["competition_id"] = flight.CompetitionID
		}

		// Add flight statistics
		stats := flight.GetStatistics(3.0) // Use 3 second speed window as default
		properties["max_altitude"] = stats.MaxAltitude
		properties["min_altitude"] = stats.MinAltitude
		properties["max_ground_speed"] = stats.MaxGroundSpeed
		properties["max_climb_rate"] = stats.MaxClimbRate
		properties["max_descent_rate"] = stats.MaxDescentRate
		properties["flight_duration_seconds"] = stats.FlightDuration.Seconds()
		properties["total_fixes"] = len(coordinates)
	}

	// Create feature
	feature := GeoJSONFeature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: properties,
	}

	// Marshal to JSON
	var result []byte
	var err error

	if pretty {
		result, err = json.MarshalIndent(feature, "", "  ")
	} else {
		result, err = json.Marshal(feature)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal GeoJSON: %w", err)
	}

	return result, nil
}
