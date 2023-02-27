package bus

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mmcloughlin/geohash"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
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
			output.Trips = make(map[string]*Trip, len(trips))
			for _, trip := range trips {
				t := &Trip{}
				t.Raw = trip

				t.Id = trip.GetId()
				t.Deleted = trip.GetIsDeleted()

				alert := trip.GetAlert()
				if alert != nil {
					c.log.Error().Msg("alert provided in trip data")
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
						stu.Stop = c.Specs.Stops[stu.StopId]

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
					if tripDescriptor != nil {
						t.DirectionId = tripDescriptor.GetDirectionId()
						t.RouteId, _ = strconv.Atoi(tripDescriptor.GetRouteId())
						t.TripId, _ = strconv.Atoi(tripDescriptor.GetTripId())
						//tripDescriptor.GetScheduleRelationship()
						t.StartDate = tripDescriptor.GetStartDate()
						t.StartTime = tripDescriptor.GetStartTime()

						t.Route = c.Specs.GetRoute(t.RouteId)
						t.Trip = c.Specs.GetTrip(t.TripId, t.RouteId)

					}

					/* Vehicle info isn't provided
					vehicleDescriptor := tripUpdate.GetVehicle()
					if vehicleDescriptor != nil {
						vehicleDescriptor.GetId()
						vehicleDescriptor.GetLabel()
						vehicleDescriptor.GetLicensePlate()
					}
					*/
				}

				/* Vehicle info isn't provided
				VehiclePosition := trip.GetVehicle()
				t.Vehicle = &Vehicle{}
				if VehiclePosition != nil {
					t.Vehicle.CongestionLevel = VehiclePosition.GetCongestionLevel().String()
					t.Vehicle.CurrentStatus = VehiclePosition.GetCurrentStatus().String()

					VehiclePosition.GetCurrentStopSequence()

					t.Vehicle.OccupancyStatus = VehiclePosition.GetOccupancyStatus().String()

					position := VehiclePosition.GetPosition()
					if position != nil {
						t.Vehicle.Bearing = position.GetBearing()
						t.Vehicle.Latitude = position.GetLatitude()
						t.Vehicle.Longitude = position.GetLongitude()
						t.Vehicle.Orometer = position.GetOdometer()
						t.Vehicle.Speed = position.GetSpeed()
					}

					t.StopId, _ = strconv.Atoi(VehiclePosition.GetStopId())
					t.Vehicle.Timestamp = time.Unix(int64(VehiclePosition.GetTimestamp()), 0)

					vTripDescriptor := VehiclePosition.GetTrip()
					if vTripDescriptor != nil {
						t.Vehicle.DirectionId = vTripDescriptor.GetDirectionId()
						t.Vehicle.RouteId, _ = strconv.Atoi(vTripDescriptor.GetRouteId())
						//vTripDescriptor.GetScheduleRelationship()
						t.Vehicle.TripId, _ = strconv.Atoi(vTripDescriptor.GetTripId())
					}
					vVehicleDescriptor := VehiclePosition.GetVehicle()
					if vVehicleDescriptor != nil {
						t.Vehicle.VehicleId = vVehicleDescriptor.GetId()
						t.Vehicle.VehicleLabel = vVehicleDescriptor.GetLabel()
						t.Vehicle.LicensePlate = vVehicleDescriptor.GetLicensePlate()
					}

				}
				*/
				output.Trips[t.Route.ShortName] = t
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
			output.Vehicles = make(map[string]map[string]*Vehicle, len(vehicles))
			for _, vehicle := range vehicles {
				v := &Vehicle{}
				v.Raw = vehicle
				v.Id = vehicle.GetId()
				v.Deleted = vehicle.GetIsDeleted()

				vehicle.GetAlert()

				vehiclePosition := vehicle.GetVehicle()
				if vehiclePosition != nil {
					v.CongestionLevel = vehiclePosition.GetCongestionLevel().String()
					v.StopStatus = vehiclePosition.GetCurrentStatus().String()
					v.CurrentStopSequence = vehiclePosition.GetCurrentStopSequence()
					v.OccupancyStatus = vehiclePosition.GetOccupancyStatus().String()
					v.Timestamp = time.Unix(int64(vehiclePosition.GetTimestamp()), 0)

					position := vehiclePosition.GetPosition()
					if position != nil {
						v.Bearing = position.GetBearing()
						v.Latitude = position.GetLatitude()
						v.Longitude = position.GetLongitude()
						v.Odometer = position.GetOdometer()
						v.Speed = position.GetSpeed()
						if v.Latitude != 0 && v.Longitude != 0 {
							v.Geohash = geohash.Encode(float64(v.Latitude), float64(v.Longitude))
						}
					}

					trip := vehiclePosition.GetTrip()
					if trip != nil {
						v.DirectionId = trip.GetDirectionId()
						v.RouteId, _ = strconv.Atoi(trip.GetRouteId())
						v.Route = c.Specs.GetRoute(v.RouteId)
						v.Agency = c.Specs.GetAgency(v.Route.AgencyId) //c.Specs.Agencies[v.Route.AgencyId]

						v.ScheduleRelationship = trip.GetScheduleRelationship().String()
						v.StartDate = trip.GetStartDate()
						v.StartTime = trip.GetStartTime()
						v.TripId, _ = strconv.Atoi(trip.GetTripId())
						v.Trip = c.Specs.GetTrip(v.TripId, v.RouteId)

						v.TripStartDate, _ = time.Parse("20060102", vehicle.Vehicle.GetTrip().GetStartDate())
					}

					vehicleDescriptor := vehiclePosition.GetVehicle()
					if vehicleDescriptor != nil {
						v.VehicleId = vehicleDescriptor.GetId()
						v.VehicleLabel = vehicleDescriptor.GetLabel()
						v.LicensePlate = vehicleDescriptor.GetLicensePlate()
					}

				}

				/* Trip info isn't provided
				tripUpdate := vehicle.GetTripUpdate()
				if tripUpdate != nil {
					tripUpdate.GetDelay()
					time.Unix(int64(tripUpdate.GetTimestamp()), 0)

					tripUpdate.GetTrip()
					tripUpdate.GetVehicle()
					tripUpdate.GetStopTimeUpdate()
				}
				*/

				if _, ok := output.Vehicles[v.Route.ShortName]; !ok {
					output.Vehicles[v.Route.ShortName] = make(map[string]*Vehicle)
				}
				output.Vehicles[v.Route.ShortName][v.Id] = v
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
