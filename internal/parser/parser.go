package parser

import (
	"fmt"
	"os"
	"time"

	"igc-tool/internal/flight"

	"github.com/twpayne/go-igc"
)

// getHRecordValue extracts the value from an H record if it exists
func getHRecordValue(records map[string]*igc.HRecord, key string) string {
	if record, exists := records[key]; exists && record != nil {
		return record.Value
	}
	return ""
}

// ParseIGCFile parses an IGC file and returns a Flight struct
func ParseIGCFile(filename string) (*flight.Flight, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	igcData, err := igc.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IGC file: %w", err)
	}

	// Check if the file has any valid IGC data
	if len(igcData.HRecordsByTLC) == 0 && len(igcData.BRecords) == 0 {
		return nil, fmt.Errorf("file does not contain valid IGC data")
	}

	// Convert from go-igc format to our internal format
	var f flight.Flight

	// Extract date from HFDTE record
	if hfdteRecord, exists := igcData.HRecordsByTLC["DTE"]; exists && hfdteRecord != nil {
		if hfdteRecord.Value != "" && len(hfdteRecord.Value) >= 6 {
			dateStr := hfdteRecord.Value[:6] // DDMMYY format
			// Use time.Parse with Go's reference time format for DDMMYY (020106)
			if parsedDate, parseErr := time.Parse("020106", dateStr); parseErr == nil {
				f.Date = parsedDate
			}
		}
	}

	// Extract pilot information from H records
	f.Pilot = getHRecordValue(igcData.HRecordsByTLC, "PLT")
	f.Crew = getHRecordValue(igcData.HRecordsByTLC, "CM2")
	f.GliderType = getHRecordValue(igcData.HRecordsByTLC, "GTY")
	f.GliderID = getHRecordValue(igcData.HRecordsByTLC, "GID")
	f.CompetitionID = getHRecordValue(igcData.HRecordsByTLC, "CID")
	f.GPSDatum = getHRecordValue(igcData.HRecordsByTLC, "DTM")
	f.FirmwareVersion = getHRecordValue(igcData.HRecordsByTLC, "RFW")
	f.HardwareVersion = getHRecordValue(igcData.HRecordsByTLC, "RHW")
	f.FlightRecorderType = getHRecordValue(igcData.HRecordsByTLC, "FTY")
	f.GPSReceiver = getHRecordValue(igcData.HRecordsByTLC, "GPS")
	f.TimeZone = getHRecordValue(igcData.HRecordsByTLC, "TZN")
	f.PressureAltSensor = getHRecordValue(igcData.HRecordsByTLC, "PRS")
	f.AltGPSRef = getHRecordValue(igcData.HRecordsByTLC, "ALG")
	f.AltPressureRef = getHRecordValue(igcData.HRecordsByTLC, "ALP")

	// Convert B records to our Fix format
	f.Fixes = igcData.BRecords

	return &f, nil
}
