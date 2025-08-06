package display

import (
	"fmt"

	"igc-tool/internal/flight"
	"igc-tool/internal/units"
	"igc-tool/internal/utils"

	"github.com/twpayne/go-igc"
)

// PrintFlightHeaders prints the flight header information
func PrintFlightHeaders(f *flight.Flight) {
	// Print parsed header data
	fmt.Printf("Date: %s\n", f.Date.Format("2006-01-02"))
	fmt.Printf("Pilot: %s\n", f.Pilot)
	if f.Crew != "" && f.Crew != "NIL" {
		fmt.Printf("Crew: %s\n", f.Crew)
	}
	fmt.Printf("Glider Type: %s\n", f.GliderType)
	if f.GliderID != "" && f.GliderID != "NKN" {
		fmt.Printf("Glider ID: %s\n", f.GliderID)
	}
	if f.CompetitionID != "" && f.CompetitionID != "NKN" {
		fmt.Printf("Competition ID: %s\n", f.CompetitionID)
	}
	if f.GPSDatum != "" {
		fmt.Printf("GPS Datum: %s\n", f.GPSDatum)
	}
	if f.FirmwareVersion != "" {
		fmt.Printf("Firmware Version: %s\n", f.FirmwareVersion)
	}
	if f.HardwareVersion != "" {
		fmt.Printf("Hardware Version: %s\n", f.HardwareVersion)
	}
	if f.FlightRecorderType != "" {
		fmt.Printf("Flight Recorder Type: %s\n", f.FlightRecorderType)
	}
	if f.GPSReceiver != "" {
		fmt.Printf("GPS Receiver: %s\n", f.GPSReceiver)
	}
	if f.TimeZone != "" {
		fmt.Printf("Time Zone: %s\n", f.TimeZone)
	}
	if f.PressureAltSensor != "" {
		fmt.Printf("Pressure Altitude Sensor: %s\n", f.PressureAltSensor)
	}
	if f.AltGPSRef != "" {
		fmt.Printf("GPS Altitude Reference: %s\n", f.AltGPSRef)
	}
	if f.AltPressureRef != "" {
		fmt.Printf("Pressure Altitude Reference: %s\n", f.AltPressureRef)
	}
}

// PrintFix prints a single fix with formatting
func PrintFix(fix *igc.BRecord, prefix string, altitudeUnit string, timeFormat string) {
	altitudeSymbol := units.AltitudeSymbol(altitudeUnit)
	altGPS := int(units.Altitude(float64(fix.AltWGS84), altitudeUnit))
	altBaro := int(units.Altitude(float64(fix.AltBarometric), altitudeUnit))
	timeStr := utils.FormatTime(fix.Time, timeFormat)

	fmt.Printf("  %s%s: (%.5f, %.5f), Alt(GPS): %d%s, Alt(Baro): %d%s\n",
		prefix,
		timeStr,
		fix.Lat, fix.Lon,
		altGPS, altitudeSymbol,
		altBaro, altitudeSymbol,
	)
}

// PrintFlightData prints complete flight data with optional summary mode
func PrintFlightData(f *flight.Flight, summary bool, altitudeUnit string, timeFormat string) {
	PrintFlightHeaders(f)

	fmt.Printf("\nFixes (%d total):\n", len(f.Fixes))

	if summary {
		// Show only first and last fix in summary mode
		if len(f.Fixes) > 0 {
			PrintFix(f.Fixes[0], "First: ", altitudeUnit, timeFormat)

			if len(f.Fixes) > 1 {
				PrintFix(f.Fixes[len(f.Fixes)-1], "Last:  ", altitudeUnit, timeFormat)
			}
		}
	} else {
		// Show all fixes in full mode
		for _, fix := range f.Fixes {
			PrintFix(fix, "", altitudeUnit, timeFormat)
		}
	}
}
