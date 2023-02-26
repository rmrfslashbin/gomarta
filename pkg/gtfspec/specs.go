package gtfspec

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

// Specs is a container for all the GTFS spec data
type Specs struct {
	// Agencies is a map of agency IDs to agency data
	Agencies map[string]*Agency `json:"agencies"`

	// Calendars is a map of service IDs to calendar data
	Calendars map[int]*Calendar `json:"calendars"`

	// CalendarDates is a map of service IDs to calendar date data
	CalendarDates []*CalendarDate `json:"calendar_dates"`

	// Routes is a map of route IDs to route data
	Routes map[int]*Route `json:"routes"`

	// Shapes is a map of shape IDs to Sequence IDs to shape data
	Shapes map[int]map[int]*Shape `json:"shapes"`

	// Stops is a map of stop IDs to stop data
	Stops map[int]*Stop `json:"stops"`

	// StopTimes is a map of trip IDs to stop IDs to stop time data
	StopTimes map[int]map[int]*StopTime `json:"stop_times"`

	// Trips is a map of trip IDs to route IDs to trip data
	Trips map[int]map[int]*Trip `json:"trips"`
}

// AddAgency adds an agency to the Specs
func (s *Specs) AddAgency(agencyID string, agency *Agency) {
	if s.Agencies == nil {
		s.Agencies = make(map[string]*Agency, 0)
	}
	s.Agencies[agencyID] = agency
}

// AddCalendar adds a calendar to the Specs
func (s *Specs) AddCalendar(serviceId int, calendar *Calendar) {
	if s.Calendars == nil {
		s.Calendars = make(map[int]*Calendar, 0)
	}
	s.Calendars[serviceId] = calendar
}

// AddCalendarDate adds a calendar date to the Specs
func (s *Specs) AddCalendarDate(calendarDate *CalendarDate) {
	if s.CalendarDates == nil {
		s.CalendarDates = make([]*CalendarDate, 0)
	}
	s.CalendarDates = append(s.CalendarDates, calendarDate)
}

// AddRoute adds a route to the Specs
func (s *Specs) AddRoute(routeId int, route *Route) {
	if s.Routes == nil {
		s.Routes = make(map[int]*Route, 0)
	}
	s.Routes[routeId] = route
}

// AddShape adds a shape to the Specs
func (s *Specs) AddShape(shapeData *ShapeData, shape *Shape) {
	if s.Shapes == nil {
		s.Shapes = make(map[int]map[int]*Shape, 0)
	}
	if s.Shapes[shapeData.ShapeId] == nil {
		s.Shapes[shapeData.ShapeId] = make(map[int]*Shape, 0)
	}
	s.Shapes[shapeData.ShapeId][shapeData.Sequence] = shape
}

// AddStop adds a stop to the Specs
func (s *Specs) AddStop(stopId int, stop *Stop) {
	if s.Stops == nil {
		s.Stops = make(map[int]*Stop, 0)
	}
	s.Stops[stopId] = stop
}

// AddStopTime adds a stop time to the Specs
func (s *Specs) AddStopTime(stopTimeData *StopTimeData, stopTime *StopTime) {
	if s.StopTimes == nil {
		s.StopTimes = make(map[int]map[int]*StopTime, 0)
	}
	if s.StopTimes[stopTimeData.TripId] == nil {
		s.StopTimes[stopTimeData.TripId] = make(map[int]*StopTime, 0)
	}
	s.StopTimes[stopTimeData.TripId][stopTimeData.StopId] = stopTime
}

// AddTrip adds a trip to the Specs
func (s *Specs) AddTrip(tripData *TripData, trip *Trip) {
	if s.Trips == nil {
		s.Trips = make(map[int]map[int]*Trip, 0)
	}
	if s.Trips[tripData.TripId] == nil {
		s.Trips[tripData.TripId] = make(map[int]*Trip, 0)
	}
	if _, ok := s.Trips[tripData.TripId][tripData.RouteId]; ok {
		panic("Duplicate trip")
	}
	s.Trips[tripData.TripId][tripData.RouteId] = trip
}

// ToGobFile writes the Specs to a Gob file
func (s *Specs) ToGOBFile(fpqn string) error {
	fp, err := os.Create(fpqn)
	if err != nil {
		return &ErrCreatingFile{Err: err, Filename: fpqn}
	}
	defer fp.Close()
	if err := gob.NewEncoder(fp).Encode(s); err != nil {
		return &ErrMarshallingGOB{Err: err}
	}

	return nil
}

// ToJSONFile writes the specs to a JSON file
func (s *Specs) ToJSONFile(fpqn string) error {
	fp, err := os.Create(fpqn)
	if err != nil {
		return &ErrCreatingFile{Err: err, Filename: fpqn}
	}
	defer fp.Close()
	if data, err := json.Marshal(s); err != nil {
		return &ErrMarshallingJSON{Err: err}
	} else {
		if _, err := fp.Write(data); err != nil {
			return &ErrWritingFile{Err: err, Filename: fpqn}
		}
	}
	return nil
}

// FromGOBFile reads a Gob file and returns a Specs
func FromGOBFile(fpqn string) (*Specs, error) {
	fp, err := os.Open(fpqn)
	if err != nil {
		return nil, &ErrOpeningFile{Err: err, Filename: fpqn}
	}
	defer fp.Close()
	var s Specs
	if err := gob.NewDecoder(fp).Decode(&s); err != nil {
		return nil, &ErrUnmarshallingGOB{Err: err}
	}
	return &s, nil
}

/*
id:"1602"

	vehicle:{
		trip:{
			trip_id:"7118032"
			route_id:"17051"
			start_date:"20220820"
			}
		vehicle:{
			id:"2302"
			label:"1602"
		}
		position:{
			latitude:33.74004
			longitude:-84.33842
			bearing:90
		}
		timestamp:1661039730
		occupancy_status:MANY_SEATS_AVAILABLE
	}
*/

/*
type Position struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Bearing   float32 `json:"bearing"`
}
*/

/*
type VehiclePosition struct {
	ID              string   `json:"id" `
	VehicleID       string   `json:"vehicle"`
	VehicleLabel    string   `json:"vehicle_label"`
	TripID          int      `json:"trip_id"`
	Trip            Trip     `json:"trip"`
	RouteID         int      `json:"route_id"`
	Route           Route    `json:"route"`
	TripStartDate   string   `json:"start_date"`
	Position        Position `json:"position"`
	Timestamp       uint64   `json:"timestamp"`
	OccupancyStatus string   `json:"occupancy_status"`
}
*/
