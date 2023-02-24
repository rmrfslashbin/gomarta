package gtfspec

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

// Input is the input for the Update function
type Input struct {
	// Url is the URL to the GTFS zip file
	Url string

	// Log is the logger to use
	Log *zerolog.Logger
}

func Update(input *Input) (*Specs, error) {
	resp, err := http.Get(input.Url)
	if err != nil {
		return nil, &ErrFetchingURL{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ErrReadingUrlBody{Err: err}
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, &ErrZipReader{Err: err}
	}

	specs := &Specs{}

	for _, file := range zipReader.File {
		zipData, err := readZipFile(file)
		if err != nil {
			return nil, &ErrZipFileReader{Err: err}
		}
		csvData, err := csv.NewReader(bytes.NewReader(zipData)).ReadAll()
		if err != nil {
			return nil, &ErrCSVReader{Err: err}
		}

		switch file.Name {
		case "agency.txt":
			input.Log.Info().Msg("parsing agency.txt")

			for _, row := range csvData[1:] {
				agency := &Agency{}
				if agencyID, err := agency.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddAgency(agencyID, agency)
				}

			}

		case "calendar.txt":
			input.Log.Info().Msg("parsing calendar.txt")

			for _, row := range csvData[1:] {
				calendar := &Calendar{}

				if serviceId, err := calendar.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddCalendar(serviceId, calendar)
				}
			}

		case "calendar_dates.txt":
			input.Log.Info().Msg("parsing calendar_dates.txt")

			for _, row := range csvData[1:] {
				calendarDate := &CalendarDate{}

				if err := calendarDate.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				}
				specs.AddCalendarDate(calendarDate)
			}

		case "routes.txt":
			input.Log.Info().Msg("parsing routes.txt")

			for _, row := range csvData[1:] {
				route := &Route{}
				if routeId, err := route.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddRoute(routeId, route)
				}
			}

		case "shapes.txt":
			input.Log.Info().Msg("parsing shapes.txt")

			for _, row := range csvData[1:] {
				shape := &Shape{}
				if shapeData, err := shape.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddShape(shapeData, shape)
				}
			}

		case "stop_times.txt":
			input.Log.Info().Msg("parsing stop_times.txt")

			for _, row := range csvData[1:] {
				stopTime := &StopTime{}
				if stopTimeData, err := stopTime.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddStopTime(stopTimeData, stopTime)
				}
			}

		case "stops.txt":
			input.Log.Info().Msg("parsing stops.txt")

			for _, row := range csvData[1:] {
				stop := &Stop{}
				if stopId, err := stop.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddStop(stopId, stop)
				}
			}

		case "trips.txt":
			input.Log.Info().Msg("parsing trips.txt")

			for _, row := range csvData[1:] {
				trip := &Trip{}
				if tripData, err := trip.Add(row); err != nil {
					return nil, &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.AddTrip(tripData, trip)
				}
			}
		}
	}
	return specs, nil
}

// readZipFile reads the contents of a zip file
func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
