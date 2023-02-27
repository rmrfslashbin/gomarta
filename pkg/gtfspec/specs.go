package gtfspec

// Specs is a container for all the GTFS spec data
type Specs struct {
	// Agencies is a map of agency IDs to agency data
	Agencies []*Agency `json:"agencies"`

	// Calendars is a map of service IDs to calendar data
	Calendars []*Calendar `json:"calendars"`

	// CalendarDates is a map of service IDs to calendar date data
	CalendarDates []*CalendarDate `json:"calendar_dates"`

	// Routes is a map of route IDs to route data
	Routes []*Route `json:"routes"`

	// Shapes is a map of shape IDs to Sequence IDs to shape data
	Shapes []*Shape `json:"shapes"`

	// Stops is a map of stop IDs to stop data
	Stops []*Stop `json:"stops"`

	// StopTimes is a map of trip IDs to stop IDs to stop time data
	StopTimes []*StopTime `json:"stop_times"`

	// Trips is a map of trip IDs to route IDs to trip data
	Trips []*Trip `json:"trips"`
}

/*
func (s *Specs) GetAgency(agencyId string) *Agency {
	if _, ok := s.Agencies[agencyId]; !ok {
		return nil
	}
	return s.Agencies[agencyId]
}

func (s *Specs) GetRoute(routeId int) *Route {
	if _, ok := s.Routes[routeId]; !ok {
		return nil
	}
	return s.Routes[routeId]
}

func (s *Specs) GetStop(stopId int) *Stop {
	if _, ok := s.Stops[stopId]; !ok {
		return nil
	}
	return s.Stops[stopId]
}

func (s *Specs) GetTrip(tripId int, routeId int) *Trip {
	if _, ok := s.Trips[tripId]; !ok {
		return nil
	}
	if _, ok := s.Trips[tripId][routeId]; !ok {
		return nil
	}
	return s.Trips[tripId][routeId]
}
*/
