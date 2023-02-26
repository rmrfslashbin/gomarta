package bus

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mmcloughlin/geohash"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type Vehicle struct {
	Id string

	//EpochTimestamp  uint64
	Timestamp       time.Time
	OccupancyStatus string

	Latitude  float32
	Longitude float32
	Bearing   float32
	Speed     float32
	Geohash   string

	VehicleId    string
	VehicleLabel string

	TripId        int
	RouteId       int
	DirectionId   uint32
	TripStartDate time.Time

	CongestionLevel string
	CurrentStatus   string

	Agency *gtfspec.Agency
	Route  *gtfspec.Route
	Trip   *gtfspec.Trip
}

type Trip struct {
	Id      string
	Deleted bool

	Delay     int32
	Timestamp time.Time

	StopTimeUpdate []*StopTimeUpdate

	DirectionId uint32
	RouteId     int
	TripId      int
	StartTime   string
	StartDate   string

	Vehicle *Vehicle
}

type StopTimeUpdate struct {
	StopSequence uint32
	StopId       int
	Arrival      *Arrival
	Departure    *Departure
}

type Arrival struct {
	Delay       int32
	Time        time.Time
	Uncertainty int32
}

type Departure struct {
	Delay       int32
	Time        time.Time
	Uncertainty int32
}

// Options for the bus instance
type Option func(c *Bus)

// Bus for the app instance
type Bus struct {
	log         *zerolog.Logger
	VehiclesUrl string
	TripsUrl    string
	Specs       *gtfspec.Specs
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

	if cfg.Specs == nil {
		return nil, &ErrSpecsNotSet{}
	}

	return cfg, nil
}

// WithLogger sets the logger for the bus instance
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Bus) {
		c.log = log
	}
}

// WithSpecs sets the specs for the bus instance
func WithSpecs(specs *gtfspec.Specs) Option {
	return func(c *Bus) {
		c.Specs = specs
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
	Trips    []*Trip
	Vehicles []*Vehicle
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
			output.Trips = make([]*Trip, len(trips))
			for ndx, trip := range trips {
				t := &Trip{}

				t.Id = trip.GetId()
				t.Deleted = trip.GetIsDeleted()

				alert := trip.GetAlert()
				if alert != nil {
					spew.Dump(alert)
					/*
						alert.GetActivePeriod()
						alert.GetCause().String()
						alert.GetDescriptionText().String()
						alert.GetEffect().String()
						alert.GetHeaderText().String()
						alert.GetInformedEntity()
						alert.GetUrl().String()
					*/
				}

				tripUpdate := trip.GetTripUpdate()
				if tripUpdate != nil {
					t.Delay = tripUpdate.GetDelay()
					t.Timestamp = time.Unix(int64(tripUpdate.GetTimestamp()), 0)

					t.StopTimeUpdate = make([]*StopTimeUpdate, len(tripUpdate.GetStopTimeUpdate()))

					for ndx, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
						stu := &StopTimeUpdate{}

						stu.StopSequence = stopTimeUpdate.GetStopSequence()
						stu.StopId, _ = strconv.Atoi(stopTimeUpdate.GetStopId())

						arrival := stopTimeUpdate.GetArrival()
						if arrival != nil {
							stu.Arrival = &Arrival{
								Delay:       arrival.GetDelay(),
								Time:        time.Unix(arrival.GetTime(), 0),
								Uncertainty: arrival.GetUncertainty(),
							}
						}

						departure := stopTimeUpdate.GetDeparture()
						if departure != nil {
							stu.Departure = &Departure{
								Delay:       departure.GetDelay(),
								Time:        time.Unix(departure.GetTime(), 0),
								Uncertainty: departure.GetUncertainty(),
							}
						}
						t.StopTimeUpdate[ndx] = stu
					}

					tripDescriptor := tripUpdate.GetTrip()
					if tripDescriptor == nil {
						t.DirectionId = tripDescriptor.GetDirectionId()
						t.RouteId, _ = strconv.Atoi(tripDescriptor.GetRouteId())
						t.TripId, _ = strconv.Atoi(tripDescriptor.GetTripId())
						//tripDescriptor.GetScheduleRelationship()
						t.StartDate = tripDescriptor.GetStartDate()
						t.StartTime = tripDescriptor.GetStartTime()

					}

					vehicleDescriptor := tripUpdate.GetVehicle()
					if vehicleDescriptor != nil {
						vehicleDescriptor.GetId()
						vehicleDescriptor.GetLabel()
						vehicleDescriptor.GetLicensePlate()
					}
				}

				VehiclePosition := trip.GetVehicle()
				if VehiclePosition != nil {
					t.Vehicle = &Vehicle{
						CongestionLevel: VehiclePosition.GetCongestionLevel().String(),
						CurrentStatus:   VehiclePosition.GetCurrentStatus().String(),
					}

					VehiclePosition.GetCurrentStopSequence()
					occupancyStatus := VehiclePosition.GetOccupancyStatus().String()
					position := VehiclePosition.GetPosition()
					if position != nil {
						position.GetBearing()
						position.GetLatitude()
						position.GetLongitude()
						position.GetOdometer()
						position.GetSpeed()
					}

					VehiclePosition.GetStopId()
					VehiclePosition.GetTimestamp()

					vTripDescriptor := VehiclePosition.GetTrip()
					if vTripDescriptor != nil {
						vTripDescriptor.GetDirectionId()
						vTripDescriptor.GetRouteId()
						//vTripDescriptor.GetScheduleRelationship()
						vTripDescriptor.GetStartDate()
						vTripDescriptor.GetStartTime()
						vTripDescriptor.GetTripId()
					}
					vVehicleDescriptor := VehiclePosition.GetVehicle()
					if vVehicleDescriptor != nil {
						vVehicleDescriptor.GetId()
						vVehicleDescriptor.GetLabel()
						vVehicleDescriptor.GetLicensePlate()
					}

				}
			}
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
			output.Vehicles = make([]*Vehicle, len(vehicles))
			for ndx, vehicle := range vehicles {
				v := &Vehicle{}

				v.Id = *vehicle.Id

				//v.EpochTimestamp = vehicle.Vehicle.GetTimestamp()
				v.OccupancyStatus = vehicle.Vehicle.GetOccupancyStatus().String()

				if vehicle.Vehicle.GetPosition().Latitude != nil {
					v.Latitude = *vehicle.Vehicle.GetPosition().Latitude
				}
				if vehicle.Vehicle.GetPosition().Longitude != nil {
					v.Longitude = *vehicle.Vehicle.GetPosition().Longitude
				}
				if vehicle.Vehicle.GetPosition().Bearing != nil {

					v.Bearing = *vehicle.Vehicle.GetPosition().Bearing
				}
				if vehicle.Vehicle.GetPosition().Speed != nil {
					v.Speed = *vehicle.Vehicle.GetPosition().Speed
				}

				if v.Latitude != 0 && v.Longitude != 0 {
					v.Geohash = geohash.Encode(float64(v.Latitude), float64(v.Longitude))
				}

				v.VehicleId = vehicle.Vehicle.GetVehicle().GetId()
				v.VehicleLabel = vehicle.Vehicle.GetVehicle().GetLabel()

				//v.TripId = vehicle.Vehicle.GetTrip().GetTripId()
				//v.RouteId = vehicle.Vehicle.GetTrip().GetRouteId()

				v.TripId, err = strconv.Atoi(vehicle.Vehicle.GetTrip().GetTripId())
				if err != nil {
					return nil, err
				}
				v.RouteId, err = strconv.Atoi(vehicle.Vehicle.GetTrip().GetRouteId())
				if err != nil {
					return nil, err
				}

				v.DirectionId = vehicle.Vehicle.GetTrip().GetDirectionId()
				v.TripStartDate, err = time.Parse("20060102", vehicle.Vehicle.GetTrip().GetStartDate())
				if err != nil {
					return nil, err
				}

				v.Timestamp = time.Unix(int64(vehicle.Vehicle.GetTimestamp()), 0)

				if _, ok := c.Specs.Trips[v.TripId]; ok {
					if _, ok := c.Specs.Trips[v.TripId][v.RouteId]; ok {
						v.Trip = c.Specs.Trips[v.TripId][v.RouteId]
					}
				}

				if _, ok := c.Specs.Routes[v.RouteId]; ok {
					v.Route = c.Specs.Routes[v.RouteId]
				}

				if strings.TrimSpace(v.Route.AgencyId) != "" {
					v.Agency = c.Specs.Agencies[v.Route.AgencyId]
				}

				output.Vehicles[ndx] = v
			}
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
