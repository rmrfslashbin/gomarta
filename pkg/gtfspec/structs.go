package gtfspec

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url
// MARTA,Metropolitan Atlanta Rapid Transit Authority,https://www.itsmarta.com,America/New_York,en,404-848-5000,https://www.itsmarta.com/fare-programs.aspx
type Agency struct {
	gorm.Model
	AgencyId string `json:"agency_id" gorm:"primary_key"`
	Name     string `json:"agency_name"`
	Url      string `json:"agency_url"`
	Timezone string `json:"agency_timezone"`
	Lang     string `json:"agency_lang"`
	Phone    string `json:"agency_phone"`
	FareUrl  string `json:"agency_fare_url"`
}

// service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date
// 20,0,0,0,0,0,0,0,20220423,20220812
type Calendar struct {
	gorm.Model
	ServiceId int       `json:"service_id" gorm:"primary_key"`
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
	gorm.Model
	ServiceId     int       `json:"service_id" gorm:"primary_key"`
	Date          time.Time `json:"date" gorm:"primary_key"`
	ExceptionType int       `json:"exception_type"`
}

// route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color
// 16883,MARTA,1,Marietta Blvd/Joseph E Lowery Blvd,,3,https://itsmarta.com/1.aspx,FF00FF,000000
type Route struct {
	gorm.Model
	RouteID   int     `json:"route_id" gorm:"primary_key"`
	AgencyId  string  `json:"agency_id"`
	ShortName string  `json:"route_short_name"`
	LongName  string  `json:"route_long_name"`
	Desc      string  `json:"route_desc"`
	Type      int     `json:"route_type"`
	Url       string  `json:"route_url"`
	Color     []uint8 `json:"route_color"`
	TextColor []uint8 `json:"route_text_color"`
}

// shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
// 100095,33.818860,-84.450519,1,0.0000
type Shape struct {
	gorm.Model
	ShapeId  int     `json:"shape_id" gorm:"primary_key"`
	Lat      float64 `json:"shape_pt_lat"`
	Lon      float64 `json:"shape_pt_lon"`
	Sequence int     `json:"shape_pt_sequence" gorm:"primary_key"`
	Distance float64 `json:"shape_dist_traveled"`
}

// trip_id,arrival_time,departure_time,stop_id,stop_sequence,stop_headsign,pickup_type,drop_off_type,shape_dist_traveled,timepoint
// 7142673, 6:43:00, 6:43:00,27,1,,0,0,,1
type StopTime struct {
	gorm.Model
	TripId            int     `json:"trip_id" gorm:"primary_key"`
	ArrivalTime       string  `json:"arrival_time"`
	DepartureTime     string  `json:"departure_time"`
	StopId            int     `json:"stop_id"`
	StopSequence      int     `json:"stop_sequence" gorm:"primary_key"`
	StopHeadsign      string  `json:"stop_headsign"`
	PickupType        int     `json:"pickup_type"`
	DropOffType       int     `json:"drop_off_type"`
	ShapeDistTraveled float64 `json:"shape_dist_traveled"`
	Timepoint         int     `json:"timepoint"`
}

// stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station,stop_timezone,wheelchair_boarding
// 27,907933,HAMILTON E HOLMES STATION,70 HAMILTON E HOLMES DR NW & CSX TRANSPORTATION,33.754553,-84.469302,,,,,,1
type Stop struct {
	gorm.Model
	StopId             int     `json:"stop_id" gorm:"primary_key"`
	Code               int     `json:"stop_code"`
	Name               string  `json:"stop_name"`
	Desc               string  `json:"stop_desc"`
	Lat                float64 `json:"stop_lat"`
	Lon                float64 `json:"stop_lon"`
	ZoneId             int     `json:"zone_id"`
	Url                string  `json:"stop_url"`
	LocationType       int     `json:"location_type"`
	ParentStation      string  `json:"parent_station"`
	Timezone           string  `json:"stop_timezone"`
	WheelchairBoarding bool    `json:"wheelchair_boarding"`
}

// route_id,service_id,trip_id,trip_headsign,trip_short_name,direction_id,block_id,shape_id,wheelchair_accessible,bikes_allowed
// 17114,2,7142675,BLUE EASTBOUND TO INDIAN CREEK STATION,,0,1075016,100750,0,0
type Trip struct {
	gorm.Model
	RouteId      int    `json:"route_id"`
	ServiceId    int    `json:"service_id"`
	TripID       int    `json:"trip_id" gorm:"primary_key"`
	Headsign     string `json:"trip_headsign"`
	ShortName    string `json:"trip_short_name"`
	DirectionId  int    `json:"direction_id"`
	BlockId      int    `json:"block_id"`
	ShapeId      int    `json:"shape_id"`
	Wheelchair   bool   `json:"wheelchair_accessible"`
	BikesAllowed bool   `json:"bikes_allowed"`
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
type Position struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Bearing   float32 `json:"bearing"`
}

type VehiclePosition struct {
	gorm.Model
	ID              string `json:"id" gorm:"primary_key"`
	VehicleID       string `json:"vehicle"`
	VehicleLabel    string `json:"vehicle_label"`
	TripID          int
	Trip            Trip
	RouteID         int
	Route           Route
	TripStartDate   string   `json:"start_date"`
	Position        Position `gorm:"embedded;embeddedPrefix:pos_"`
	Timestamp       uint64   `json:"timestamp"`
	OccupancyStatus string   `json:"occupancy_status"`
}

func (a *Agency) Add(record []string) error {
	a.AgencyId = record[0]
	a.Name = record[1]
	a.Url = record[2]
	a.Timezone = record[3]
	a.Lang = record[4]
	a.Phone = record[5]
	a.FareUrl = record[6]

	return nil
}

func (c *Calendar) Add(record []string) error {
	var err error
	if c.ServiceId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("service_id: %v", err)
	}
	if c.Monday, err = strconv.Atoi(record[1]); err != nil {
		return fmt.Errorf("monday: %v", err)
	}
	if c.Tuesday, err = strconv.Atoi(record[2]); err != nil {
		return fmt.Errorf("tuesday: %v", err)
	}
	if c.Wednesday, err = strconv.Atoi(record[3]); err != nil {
		return fmt.Errorf("wednesday: %v", err)
	}
	if c.Thursday, err = strconv.Atoi(record[4]); err != nil {
		return fmt.Errorf("thursday: %v", err)
	}
	if c.Friday, err = strconv.Atoi(record[5]); err != nil {
		return fmt.Errorf("friday: %v", err)
	}
	if c.Saturday, err = strconv.Atoi(record[6]); err != nil {
		return fmt.Errorf("saturday: %v", err)
	}
	if c.Sunday, err = strconv.Atoi(record[7]); err != nil {
		return fmt.Errorf("sunday: %v", err)
	}
	if c.StartDate, err = time.Parse("20060102", record[8]); err != nil {
		return fmt.Errorf("start_date: %v", err)
	}
	if c.EndDate, err = time.Parse("20060102", record[9]); err != nil {
		return fmt.Errorf("end_date: %v", err)
	}

	return nil
}

func (c *CalendarDate) Add(record []string) error {
	var err error
	if c.ServiceId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("service_id: %v", err)
	}
	if c.Date, err = time.Parse("20060102", record[1]); err != nil {
		return fmt.Errorf("date: %v", err)
	}
	if c.ExceptionType, err = strconv.Atoi(record[2]); err != nil {
		return fmt.Errorf("exception_type: %v", err)
	}

	//pair := EncodePair(int64(c.ServiceId), int64(c.Date.Unix()))
	return nil
}

func (r *Route) Add(record []string) error {
	var err error

	r.AgencyId = record[1]
	r.ShortName = record[2]
	r.LongName = record[3]
	r.Desc = record[4]
	r.Url = record[6]

	if r.RouteID, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("route id: %v", err)
	}
	if r.Type, err = strconv.Atoi(record[5]); err != nil {
		return fmt.Errorf("route type: %v", err)
	}
	if color, err := hex.DecodeString(record[7]); err != nil {
		return fmt.Errorf("route color: %v", err)
	} else {
		r.Color = color
	}
	if textColor, err := hex.DecodeString(record[8]); err != nil {
		return fmt.Errorf("route text color: %v", err)
	} else {
		r.TextColor = textColor
	}

	return nil
}

func (s *Shape) Add(record []string) error {
	var err error

	if s.ShapeId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("shape id: %v", err)
	}
	if s.Lat, err = strconv.ParseFloat(record[1], 64); err != nil {
		return fmt.Errorf("shape lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[2], 64); err != nil {
		return fmt.Errorf("shape lon: %v", err)
	}
	if s.Sequence, err = strconv.Atoi(record[3]); err != nil {
		return fmt.Errorf("shape sequence: %v", err)
	}
	if s.Distance, err = strconv.ParseFloat(record[4], 64); err != nil {
		return fmt.Errorf("shape distance: %v", err)
	}

	return nil
}

func (s *StopTime) Add(record []string) error {
	var err error
	s.ArrivalTime = record[1]
	s.DepartureTime = record[2]
	s.StopHeadsign = record[5]

	if s.TripId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("trip id: %v", err)
	}
	if s.StopId, err = strconv.Atoi(record[3]); err != nil {
		return fmt.Errorf("stop id: %v", err)
	}
	if s.StopSequence, err = strconv.Atoi(record[4]); err != nil {
		return fmt.Errorf("stop sequence: %v", err)
	}
	if s.PickupType, err = strconv.Atoi(record[6]); err != nil {
		return fmt.Errorf("pickup type: %v", err)
	}
	if s.DropOffType, err = strconv.Atoi(record[7]); err != nil {
		return fmt.Errorf("drop off type: %v", err)
	}
	if s.ShapeDistTraveled, err = strconv.ParseFloat(strings.TrimSpace(record[8]), 64); err != nil {
		if strings.TrimSpace(record[8]) == "" {
			s.ShapeDistTraveled = 0
		} else {
			return fmt.Errorf("shape dist traveled: %v", err)
		}
	}
	if s.Timepoint, err = strconv.Atoi(record[9]); err != nil {
		return fmt.Errorf("timepoint: %v", err)
	}

	return nil
}

func (s *Stop) Add(record []string) error {
	var err error
	s.Name = record[2]
	s.Desc = record[3]
	s.Url = record[7]
	s.ParentStation = record[9]

	if s.StopId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("stop id: %v", err)
	}
	if s.Code, err = strconv.Atoi(record[1]); err != nil {
		return fmt.Errorf("stop code: %v", err)
	}
	if s.Lat, err = strconv.ParseFloat(record[4], 64); err != nil {
		return fmt.Errorf("stop lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[5], 64); err != nil {
		return fmt.Errorf("stop lon: %v", err)
	}
	if s.ZoneId, err = strconv.Atoi(record[6]); err != nil {
		if strings.TrimSpace(record[6]) == "" {
			s.ZoneId = 0
		} else {
			return fmt.Errorf("stop zone id: %v", err)
		}
	}
	if s.LocationType, err = strconv.Atoi(record[8]); err != nil {
		if strings.TrimSpace(record[8]) == "" {
			s.LocationType = 0
		} else {
			return fmt.Errorf("stop location type: %v", err)
		}
	}
	if strings.TrimSpace(record[10]) == "1" {
		s.WheelchairBoarding = true
	} else {
		s.WheelchairBoarding = false
	}
	return nil
}

func (t *Trip) Add(record []string) error {
	var err error
	t.Headsign = strings.TrimSpace(record[3])
	t.ShortName = strings.TrimSpace(record[4])

	if t.RouteId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("route id: %v", err)
	}
	if t.ServiceId, err = strconv.Atoi(record[1]); err != nil {
		return fmt.Errorf("service id: %v", err)
	}
	if t.BlockId, err = strconv.Atoi(record[5]); err != nil {
		return fmt.Errorf("block id: %v", err)
	}
	if t.DirectionId, err = strconv.Atoi(record[6]); err != nil {
		return fmt.Errorf("direction id: %v", err)
	}
	if t.ShapeId, err = strconv.Atoi(record[7]); err != nil {
		return fmt.Errorf("shape id: %v", err)
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
	if t.TripID, err = strconv.Atoi(record[2]); err != nil {
		return fmt.Errorf("trip id: %v", err)
	}
	return nil
}
