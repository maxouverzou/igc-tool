package sites

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"igc-tool/internal/flight"
	"igc-tool/internal/utils"

	"github.com/paulmach/orb"
)

// LandingSite represents a landing site with name, center point, and radius
type LandingSite struct {
	Name   string
	Center orb.Point
	Radius float64 // radius in meters
}

// Collection holds a collection of landing sites
type Collection struct {
	Sites []LandingSite
}

// LoadLandingSites loads landing sites from a CSV file
func LoadLandingSites(filename string) (*Collection, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read landing sites file %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return &Collection{Sites: []LandingSite{}}, nil
	}

	var sites []LandingSite

	// Skip header row if it exists (check if first row has "name" as first column)
	startRow := 0
	if len(records) > 0 && len(records[0]) > 0 && records[0][0] == "name" {
		startRow = 1
	}

	for i := startRow; i < len(records); i++ {
		record := records[i]
		if len(record) != 4 {
			continue // Skip rows that don't have exactly 4 columns
		}

		name := record[0]
		if name == "" {
			continue
		}

		lat, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		lon, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}

		radius, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			continue
		}

		sites = append(sites, LandingSite{
			Name:   name,
			Center: orb.Point{lon, lat}, // orb.Point is [longitude, latitude]
			Radius: radius,
		})
	}

	return &Collection{Sites: sites}, nil
}

// FindLandingSite finds the landing site name for given coordinates
func (c *Collection) FindLandingSite(lat, lon float64) string {
	for _, site := range c.Sites {
		siteLat := site.Center[1]
		siteLon := site.Center[0]
		distance := flight.HaversineDistance(lat, lon, siteLat, siteLon)

		if distance <= site.Radius {
			return site.Name
		}
	}
	return utils.FormatCoordinates(lat, lon)
}
