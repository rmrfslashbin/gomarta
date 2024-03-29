package bus

import (
	"time"

	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
)

// Arrival is a struct for arrival data
type Arrival struct {
	Delay       int32
	Time        time.Time
	Uncertainty int32
}

// Departure is a struct for departure data
type Departure struct {
	Delay       int32
	Time        time.Time
	Uncertainty int32
}

// FetchOutput is the output for the Fetch method
type FetchOutput struct {
	// Trips is a map of route "short names" (ie: bus line; ex: "37") to a map of Trip struct
	Trips    map[string]*Trip
	Vehicles map[string]map[string]*Vehicle
}

// StopTimeUpdate is a struct for stop time update data
type StopTimeUpdate struct {
	StopSequence uint32
	StopId       int
	Arrival      *Arrival
	Departure    *Departure
	Stop         *gtfspec.Stop
}

// Trip is a struct for trip data
type Trip struct {
	// Raw is the raw GTFS-RT data
	Raw *gtfsrt.FeedEntity

	// Id is the gtfs feed entity id
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

// Vehicle is a struct for vehicle data
type Vehicle struct {
	// Raw is the raw GTFS-RT data
	Raw *gtfsrt.FeedEntity

	// Id is the gtfs feed entity id
	Id      string
	Deleted bool

	//EpochTimestamp  uint64
	Timestamp       time.Time
	OccupancyStatus string
	Delay           int32

	Latitude  float32
	Longitude float32
	Bearing   float32
	Speed     float32
	Orometer  float64
	Geohash   string

	VehicleId    string
	VehicleLabel string
	LicensePlate string
	Odometer     float64

	TripId        int
	RouteId       int
	DirectionId   uint32
	TripStartDate time.Time

	CongestionLevel      string
	CurrentStatus        string
	StopStatus           string
	CurrentStopSequence  uint32
	ScheduleRelationship string
	StartDate            string
	StartTime            string

	Agency *gtfspec.Agency
	Route  *gtfspec.Route
	Trip   *gtfspec.Trip
}
