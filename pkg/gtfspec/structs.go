package gtfspec

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Agencies      map[string]*Agency      `json:"agencies"`
	Calendars     map[int]*Calendar       `json:"calendars"`
	CalendarDates map[int64]*CalendarDate `json:"calendar_dates"`
	Routes        map[int]*Route          `json:"routes"`
	Shapes        map[int64]*Shape        `json:"shapes"`
	StopTimes     map[int64]*StopTime     `json:"stop_times"`
	Stops         map[int]*Stop           `json:"stops"`
	Trips         map[int]*Trip           `json:"trips"`
}

// agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url
// MARTA,Metropolitan Atlanta Rapid Transit Authority,https://www.itsmarta.com,America/New_York,en,404-848-5000,https://www.itsmarta.com/fare-programs.aspx
type Agency struct {
	Name     string   `json:"agency_name"`
	Url      *url.URL `json:"agency_url"`
	Timezone string   `json:"agency_timezone"`
	Lang     string   `json:"agency_lang"`
	Phone    string   `json:"agency_phone"`
	FareUrl  string   `json:"agency_fare_url"`
}

// service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date
// 20,0,0,0,0,0,0,0,20220423,20220812
type Calendar struct {
	Monday    int       `json:"monday"`
	Tuesday   int       `json:"tuesday"`
	Wednesday int       `json:"wednesday"`
	Thursday  int       `json:"thursday"`
	Friday    int       `json:"friday"`
	Saturday  int       `json:"saturday"`
	Sunday    int       `json:"sunday"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// service_id,date,exception_type
// 34,20220530,1
type CalendarDate struct {
	ServiceId     int       `json:"service_id"`
	Date          time.Time `json:"date"`
	ExceptionType int       `json:"exception_type"`
}

// route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color
// 16883,MARTA,1,Marietta Blvd/Joseph E Lowery Blvd,,3,https://itsmarta.com/1.aspx,FF00FF,000000
type Route struct {
	AgencyId  string   `json:"agency_id"`
	ShortName string   `json:"route_short_name"`
	LongName  string   `json:"route_long_name"`
	Desc      string   `json:"route_desc"`
	Type      int      `json:"route_type"`
	Url       *url.URL `json:"route_url"`
	Color     []uint8  `json:"route_color"`
	TextColor []uint8  `json:"route_text_color"`
}

// shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
// 100095,33.818860,-84.450519,1,0.0000
type Shape struct {
	Id       int     `json:"shape_id"`
	Lat      float64 `json:"shape_pt_lat"`
	Lon      float64 `json:"shape_pt_lon"`
	Sequence int     `json:"shape_pt_sequence"`
	Distance float64 `json:"shape_dist_traveled"`
}

// trip_id,arrival_time,departure_time,stop_id,stop_sequence,stop_headsign,pickup_type,drop_off_type,shape_dist_traveled,timepoint
// 7142673, 6:43:00, 6:43:00,27,1,,0,0,,1
type StopTime struct {
	TripId            int     `json:"trip_id"`
	ArrivalTime       string  `json:"arrival_time"`
	DepartureTime     string  `json:"departure_time"`
	StopId            int     `json:"stop_id"`
	StopSequence      int     `json:"stop_sequence"`
	StopHeadsign      string  `json:"stop_headsign"`
	PickupType        int     `json:"pickup_type"`
	DropOffType       int     `json:"drop_off_type"`
	ShapeDistTraveled float64 `json:"shape_dist_traveled"`
	Timepoint         int     `json:"timepoint"`
}

// stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station,stop_timezone,wheelchair_boarding
// 27,907933,HAMILTON E HOLMES STATION,70 HAMILTON E HOLMES DR NW & CSX TRANSPORTATION,33.754553,-84.469302,,,,,,1
type Stop struct {
	Code               int      `json:"stop_code"`
	Name               string   `json:"stop_name"`
	Desc               string   `json:"stop_desc"`
	Lat                float64  `json:"stop_lat"`
	Lon                float64  `json:"stop_lon"`
	ZoneId             int      `json:"zone_id"`
	Url                *url.URL `json:"stop_url"`
	LocationType       int      `json:"location_type"`
	ParentStation      string   `json:"parent_station"`
	Timezone           string   `json:"stop_timezone"`
	WheelchairBoarding bool     `json:"wheelchair_boarding"`
}

// route_id,service_id,trip_id,trip_headsign,trip_short_name,direction_id,block_id,shape_id,wheelchair_accessible,bikes_allowed
// 17114,2,7142675,BLUE EASTBOUND TO INDIAN CREEK STATION,,0,1075016,100750,0,0
type Trip struct {
	RouteId      int    `json:"route_id"`
	ServiceId    int    `json:"service_id"`
	Headsign     string `json:"trip_headsign"`
	ShortName    string `json:"trip_short_name"`
	DirectionId  int    `json:"direction_id"`
	BlockId      int    `json:"block_id"`
	ShapeId      int    `json:"shape_id"`
	Wheelchair   bool   `json:"wheelchair_accessible"`
	BikesAllowed bool   `json:"bikes_allowed"`
}

func MakePair(k1, k2 int64) int64 {
	return (k1*k1+k1)/2 + k2
}

func DecodePair(pair int64) (int64, int64) {
	k1 := (pair - (pair & 1)) / 2
	k2 := pair - k1*(k1+1)/2 - 1
	return k1, k2
}

func (a *Agency) Add(record []string) (*string, error) {
	var err error

	a.Name = record[1]
	a.Timezone = record[3]
	a.Lang = record[4]
	a.Phone = record[5]
	a.FareUrl = record[6]

	if a.Url, err = url.Parse(record[2]); err != nil {
		return nil, err
	}

	return &record[0], nil
}

func (c *Calendar) Add(record []string) (*int, error) {
	var err error
	if c.Monday, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("monday: %v", err)
	}
	if c.Tuesday, err = strconv.Atoi(record[2]); err != nil {
		return nil, fmt.Errorf("tuesday: %v", err)
	}
	if c.Wednesday, err = strconv.Atoi(record[3]); err != nil {
		return nil, fmt.Errorf("wednesday: %v", err)
	}
	if c.Thursday, err = strconv.Atoi(record[4]); err != nil {
		return nil, fmt.Errorf("thursday: %v", err)
	}
	if c.Friday, err = strconv.Atoi(record[5]); err != nil {
		return nil, fmt.Errorf("friday: %v", err)
	}
	if c.Saturday, err = strconv.Atoi(record[6]); err != nil {
		return nil, fmt.Errorf("saturday: %v", err)
	}
	if c.Sunday, err = strconv.Atoi(record[7]); err != nil {
		return nil, fmt.Errorf("sunday: %v", err)
	}
	if c.StartDate, err = time.Parse("20060102", record[8]); err != nil {
		return nil, fmt.Errorf("start_date: %v", err)
	}
	if c.EndDate, err = time.Parse("20060102", record[9]); err != nil {
		return nil, fmt.Errorf("end_date: %v", err)
	}
	if serviceId, err := strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("service_id: %v", err)
	} else {
		return &serviceId, nil
	}
}

func (c *CalendarDate) Add(record []string) (*int64, error) {
	var err error
	if c.ServiceId, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("service_id: %v", err)
	}
	if c.Date, err = time.Parse("20060102", record[1]); err != nil {
		return nil, fmt.Errorf("date: %v", err)
	}
	if c.ExceptionType, err = strconv.Atoi(record[2]); err != nil {
		return nil, fmt.Errorf("exception_type: %v", err)
	}

	pair := MakePair(int64(c.ServiceId), int64(c.Date.Unix()))
	return &pair, nil
}

func (r *Route) Add(record []string) (*int, error) {
	var err error

	r.AgencyId = record[1]
	r.ShortName = record[2]
	r.LongName = record[3]
	r.Desc = record[4]

	if r.Type, err = strconv.Atoi(record[5]); err != nil {
		return nil, fmt.Errorf("route type: %v", err)
	}
	if color, err := hex.DecodeString(record[7]); err != nil {
		return nil, fmt.Errorf("route color: %v", err)
	} else {
		r.Color = color
	}
	if textColor, err := hex.DecodeString(record[8]); err != nil {
		return nil, fmt.Errorf("route text color: %v", err)
	} else {
		r.TextColor = textColor
	}
	if r.Url, err = url.Parse(record[6]); err != nil {
		return nil, fmt.Errorf("route url: %v", err)
	}

	if routeId, err := strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("route id: %v", err)
	} else {
		return &routeId, nil
	}
}

func (s *Shape) Add(record []string) (*int64, error) {
	var err error

	if s.Lat, err = strconv.ParseFloat(record[1], 64); err != nil {
		return nil, fmt.Errorf("shape lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[2], 64); err != nil {
		return nil, fmt.Errorf("shape lon: %v", err)
	}
	if s.Sequence, err = strconv.Atoi(record[3]); err != nil {
		return nil, fmt.Errorf("shape sequence: %v", err)
	}
	if s.Distance, err = strconv.ParseFloat(record[4], 64); err != nil {
		return nil, fmt.Errorf("shape distance: %v", err)
	}
	if s.Id, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("shape id: %v", err)
	}

	pair := MakePair(int64(s.Id), int64(s.Sequence))
	return &pair, nil

	// Encode
	/*
		pair := s.Id + s.Sequence
		pair = pair * (pair + 1)
		pair = pair / 2
		pair = pair + s.Sequence
		return &pair, nil
	*/
}

func (s *StopTime) Add(record []string) (*int64, error) {
	var err error
	s.ArrivalTime = record[1]
	s.DepartureTime = record[2]
	s.StopHeadsign = record[5]

	if s.TripId, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("trip id: %v", err)
	}
	if s.StopId, err = strconv.Atoi(record[3]); err != nil {
		return nil, fmt.Errorf("stop id: %v", err)
	}
	if s.StopSequence, err = strconv.Atoi(record[4]); err != nil {
		return nil, fmt.Errorf("stop sequence: %v", err)
	}
	if s.PickupType, err = strconv.Atoi(record[6]); err != nil {
		return nil, fmt.Errorf("pickup type: %v", err)
	}
	if s.DropOffType, err = strconv.Atoi(record[7]); err != nil {
		return nil, fmt.Errorf("drop off type: %v", err)
	}
	if s.ShapeDistTraveled, err = strconv.ParseFloat(strings.TrimSpace(record[8]), 64); err != nil {
		if strings.TrimSpace(record[8]) == "" {
			s.ShapeDistTraveled = 0
		} else {
			return nil, fmt.Errorf("shape dist traveled: %v", err)
		}
	}
	if s.Timepoint, err = strconv.Atoi(record[9]); err != nil {
		return nil, fmt.Errorf("timepoint: %v", err)
	}

	pair := MakePair(int64(s.TripId), int64(s.StopSequence))
	return &pair, nil
}

func (s *Stop) Add(record []string) (*int, error) {
	var err error
	s.Name = record[2]
	s.Desc = record[3]
	s.ParentStation = record[9]

	if s.Code, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("stop code: %v", err)
	}

	if s.Lat, err = strconv.ParseFloat(record[4], 64); err != nil {
		return nil, fmt.Errorf("stop lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[5], 64); err != nil {
		return nil, fmt.Errorf("stop lon: %v", err)
	}
	if s.ZoneId, err = strconv.Atoi(record[6]); err != nil {
		if strings.TrimSpace(record[6]) == "" {
			s.ZoneId = 0
		} else {
			return nil, fmt.Errorf("stop zone id: %v", err)
		}
	}
	if s.Url, err = url.Parse(record[7]); err != nil {
		return nil, fmt.Errorf("stop url: %v", err)
	}
	if s.LocationType, err = strconv.Atoi(record[8]); err != nil {
		if strings.TrimSpace(record[8]) == "" {
			s.LocationType = 0
		} else {
			return nil, fmt.Errorf("stop location type: %v", err)
		}
	}
	if strings.TrimSpace(record[10]) == "1" {
		s.WheelchairBoarding = true
	} else {
		s.WheelchairBoarding = false
	}
	if stopId, err := strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("stop id: %v", err)
	} else {
		return &stopId, nil
	}
}

func (t *Trip) Add(record []string) (*int, error) {
	var err error
	t.Headsign = strings.TrimSpace(record[3])
	t.ShortName = strings.TrimSpace(record[4])

	if t.RouteId, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("route id: %v", err)
	}
	if t.ServiceId, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("service id: %v", err)
	}
	if t.BlockId, err = strconv.Atoi(record[5]); err != nil {
		return nil, fmt.Errorf("block id: %v", err)
	}
	if t.DirectionId, err = strconv.Atoi(record[6]); err != nil {
		return nil, fmt.Errorf("direction id: %v", err)
	}
	if t.ShapeId, err = strconv.Atoi(record[7]); err != nil {
		return nil, fmt.Errorf("shape id: %v", err)
	}
	if strings.TrimSpace(record[8]) == "1" {
		t.Wheelchair = true
	} else {
		t.Wheelchair = false
	}
	if strings.TrimSpace(record[9]) == "1" {
		t.BikesAllowed = true
	} else {
		t.BikesAllowed = false
	}
	if tripId, err := strconv.Atoi(record[2]); err != nil {
		return nil, fmt.Errorf("trip id: %v", err)
	} else {
		return &tripId, nil
	}
}

func (d *Data) Write(fqpn string) error {
	fh, err := os.Create(fqpn)
	if err != nil {
		return err
	}
	defer fh.Close()

	//var buf bytes.Buffer
	zw := gzip.NewWriter(fh)
	defer zw.Close()

	enc := gob.NewEncoder(zw)
	if err := enc.Encode(d); err != nil {
		return err
	}

	return nil
}

func (d *Data) Read(fqpn string) error {
	fh, err := os.Open(fqpn)
	if err != nil {
		return err
	}
	defer fh.Close()

	zr, err := gzip.NewReader(fh)
	if err != nil {
		return err
	}
	defer zr.Close()

	dec := gob.NewDecoder(zr)
	if err := dec.Decode(d); err != nil {
		return err
	}

	return nil
}
