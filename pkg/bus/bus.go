package bus

import (
	"io"
	"net/http"
	"os"

	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

// Options for the bus instance
type Option func(c *Bus)

// Bus for the app instance
type Bus struct {
	log         *zerolog.Logger
	VehiclesUrl string
	TripsUrl    string
}

// New creates a new mastoclinet instance
func New(opts ...Option) (*Bus, error) {
	cfg := &Bus{}

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

	if cfg.TripsUrl == "" {
		cfg.TripsUrl = "https://gtfs-rt.itsmarta.com/TMGTFSRealTimeWebService/tripupdate/tripupdates.pb"
	}

	if cfg.VehiclesUrl == "" {
		cfg.VehiclesUrl = "https://gtfs-rt.itsmarta.com/TMGTFSRealTimeWebService/vehicle/vehiclepositions.pb"
	}

	return cfg, nil
}

// WithLogger sets the logger for the bus instance
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Bus) {
		c.log = log
	}
}

// WithTripsUrl sets the trips url for the app instance
func WithTripsUrl(url string) Option {
	return func(c *Bus) {
		c.TripsUrl = url
	}
}

func WithVehiclesUrl(url string) Option {
	return func(c *Bus) {
		c.VehiclesUrl = url
	}
}

// FetchInput is the input for the Fetch method
type FetchInput struct {
	Trips    bool
	Vehicles bool
}

// FetchOutput is the output for the Fetch method
type FetchOutput struct {
	Trips    []*gtfsrt.FeedEntity
	Vehicles []*gtfsrt.FeedEntity
}

// Fetch gets the current requested data from the API.
func (c *Bus) Fetch(input *FetchInput) (*FetchOutput, error) {
	output := &FetchOutput{}

	if input.Trips {
		c.log.Debug().
			Str("url", c.TripsUrl).
			Str("function", "pkg/bus.Fetch()").
			Msg("getting trip data")
		if trips, err := c.getData(c.TripsUrl); err != nil {
			return nil, err
		} else {
			output.Trips = trips
		}
	}

	if input.Vehicles {
		c.log.Debug().
			Str("url", c.VehiclesUrl).
			Str("function", "pkg/bus.Fetch()").
			Msg("getting vehicle data")
		if vehicles, err := c.getData(c.VehiclesUrl); err != nil {
			return nil, err
		} else {
			output.Vehicles = vehicles
		}
	}
	return output, nil
}

// getData gets the current bus data from the API.
func (c *Bus) getData(url string) ([]*gtfsrt.FeedEntity, error) {
	c.log.Debug().
		Str("url", url).
		Str("function", "pkg/bus.GetData()").
		Msg("getting requested data")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.log.Trace().
		Str("url", url).
		Str("function", "pkg/bus.GetData()").
		Str("body", string(body)).
		Msg("fetched data from API")

	feed := &gtfsrt.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	c.log.Debug().
		Str("url", url).
		Str("function", "pkg/bus.GetData()").
		Int("itmes", len(feed.GetEntity())).
		Msg("unmarshalled protobuf data")

	return feed.GetEntity(), nil
}
