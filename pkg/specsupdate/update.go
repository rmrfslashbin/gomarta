package specsupdate

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/rmrfslashbin/gomarta/pkg/database"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rs/zerolog"
)

// Options for the bus instance
type Option func(c *SpecsConfig)

// Database for the app instance
type SpecsConfig struct {
	log *zerolog.Logger
	url *string
	db  *database.Database
}

// New creates a new mastoclinet instance
func New(opts ...Option) (*SpecsConfig, error) {
	cfg := &SpecsConfig{}

	// apply the list of options to Bus
	for _, opt := range opts {
		opt(cfg)
	}

	// set up logger if not provided
	if cfg.log == nil {
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		cfg.log = &log
	}

	if cfg.url == nil {
		return nil, &ErrNoURL{}
	}

	if cfg.db == nil {
		return nil, &ErrNoDatabase{}
	}

	return cfg, nil
}

// WithDatabase sets the database for the bus instance
func WithDatabase(db *database.Database) Option {
	return func(c *SpecsConfig) {
		c.db = db
	}
}

// WithLogger sets the logger for the bus instance
func WithLogger(log *zerolog.Logger) Option {
	return func(c *SpecsConfig) {
		c.log = log
	}
}

// WithUrl sets the URL for the bus instance
func WithUrl(url string) Option {
	return func(c *SpecsConfig) {
		c.url = &url
	}
}

func (c *SpecsConfig) Update() error {
	resp, err := http.Get(*c.url)
	if err != nil {
		return &ErrFetchingURL{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ErrReadingUrlBody{Err: err}
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return &ErrZipReader{Err: err}
	}

	specs := &gtfspec.Specs{}

	for _, file := range zipReader.File {
		zipData, err := readZipFile(file)
		if err != nil {
			return &ErrZipFileReader{Err: err}
		}
		csvData, err := csv.NewReader(bytes.NewReader(zipData)).ReadAll()
		if err != nil {
			return &ErrCSVReader{Err: err}
		}

		headers := makeHeaders(csvData[0])

		switch file.Name {
		case "agency.txt":
			c.log.Info().Msg("parsing agency.txt")
			specs.Agencies = make([]*gtfspec.Agency, 0)

			for _, row := range csvData[1:] {
				agency := &gtfspec.Agency{}

				if err := agency.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Agencies = append(specs.Agencies, agency)
				}

			}

		case "calendar.txt":
			c.log.Info().Msg("parsing calendar.txt")
			specs.Calendars = make([]*gtfspec.Calendar, 0)

			for _, row := range csvData[1:] {
				calendar := &gtfspec.Calendar{}

				if err := calendar.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Calendars = append(specs.Calendars, calendar)
				}
			}

		case "calendar_dates.txt":
			c.log.Info().Msg("parsing calendar_dates.txt")
			specs.CalendarDates = make([]*gtfspec.CalendarDate, 0)

			for _, row := range csvData[1:] {
				calendarDate := &gtfspec.CalendarDate{}

				if err := calendarDate.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				}
				specs.CalendarDates = append(specs.CalendarDates, calendarDate)
			}

		case "routes.txt":
			c.log.Info().Msg("parsing routes.txt")
			specs.Routes = make([]*gtfspec.Route, 0)

			for _, row := range csvData[1:] {
				route := &gtfspec.Route{}
				if err := route.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Routes = append(specs.Routes, route)
				}
			}

		case "shapes.txt":
			c.log.Info().Msg("parsing shapes.txt")
			specs.Shapes = make([]*gtfspec.Shape, 0)

			for _, row := range csvData[1:] {
				shape := &gtfspec.Shape{}
				if err := shape.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Shapes = append(specs.Shapes, shape)
				}
			}

		case "stop_times.txt":
			c.log.Info().Msg("parsing stop_times.txt")
			specs.StopTimes = make([]*gtfspec.StopTime, 0)

			for _, row := range csvData[1:] {
				stopTime := &gtfspec.StopTime{}
				if err := stopTime.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.StopTimes = append(specs.StopTimes, stopTime)
				}
			}

		case "stops.txt":
			c.log.Info().Msg("parsing stops.txt")
			specs.Stops = make([]*gtfspec.Stop, 0)

			for _, row := range csvData[1:] {
				stop := &gtfspec.Stop{}
				if err := stop.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Stops = append(specs.Stops, stop)
				}
			}

		case "trips.txt":
			c.log.Info().Msg("parsing trips.txt")
			specs.Trips = make([]*gtfspec.Trip, 0)

			for _, row := range csvData[1:] {
				trip := &gtfspec.Trip{}
				if err := trip.Add(headers, row); err != nil {
					return &ErrParsingFile{Err: err, File: file.Name}
				} else {
					specs.Trips = append(specs.Trips, trip)
				}
			}
		}
	}

	c.log.Info().Msg("adding Agencies to database")
	if _, err := c.db.Create(specs.Agencies); err != nil {
		return &ErrAddingData{Err: err, Structure: "Agencies"}
	}
	c.log.Info().Msg("adding Calendars to database")
	if _, err := c.db.Create(specs.Calendars); err != nil {
		return &ErrAddingData{Err: err, Structure: "Calendars"}
	}
	c.log.Info().Msg("adding CalendarDates to database")
	if _, err := c.db.Create(specs.CalendarDates); err != nil {
		return &ErrAddingData{Err: err, Structure: "CalendarDates"}
	}
	c.log.Info().Msg("adding Routes to database")
	if _, err := c.db.Create(specs.Routes); err != nil {
		return &ErrAddingData{Err: err, Structure: "Routes"}
	}
	c.log.Info().Msg("adding Shapes to database")
	if _, err := c.db.Create(specs.Shapes); err != nil {
		return &ErrAddingData{Err: err, Structure: "Shapes"}
	}
	c.log.Info().Msg("adding Stops to database")
	if _, err := c.db.Create(specs.Stops); err != nil {
		return &ErrAddingData{Err: err, Structure: "Stops"}
	}
	c.log.Info().Msg("adding StopTimes to database")
	if _, err := c.db.Create(specs.StopTimes); err != nil {
		return &ErrAddingData{Err: err, Structure: "StopTimes"}
	}
	c.log.Info().Msg("adding Trips to database")
	if _, err := c.db.Create(specs.Trips); err != nil {
		return &ErrAddingData{Err: err, Structure: "Trips"}
	}

	return nil
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

func makeHeaders(headerRow []string) map[string]int {
	headers := make(map[string]int)

	for i, header := range headerRow {
		headers[header] = i
	}

	return headers
}
