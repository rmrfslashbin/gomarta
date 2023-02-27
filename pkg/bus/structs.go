package bus

import (
	"time"

	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
)

// FetchOutput is the output for the Fetch method
type FetchOutput struct {
	// Trips is a map of route "short names" (ie: bus line; ex: "37") to a map of Trip struct
	Trips    map[string]*Trip
	Vehicles []*Vehicle
}

type Vehicle struct {
	Id string

	//EpochTimestamp  uint64
	Timestamp       time.Time
	OccupancyStatus string

	Latitude  float32
	Longitude float32
	Bearing   float32
	Speed     float32
	Orometer  float64
	Geohash   string

	VehicleId    string
	VehicleLabel string
	LicensePlate string

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
	Raw     *gtfsrt.FeedEntity
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

	Trip  *gtfspec.Trip
	Route *gtfspec.Route

	// Vehicle info isn't currently used
	//Vehicle *Vehicle
}

type StopTimeUpdate struct {
	StopSequence uint32
	StopId       int
	Arrival      *Arrival
	Departure    *Departure
	Stop         *gtfspec.Stop
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
