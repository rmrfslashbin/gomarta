package gtfspec

import (
	"encoding/hex"
	"strconv"
	"time"
)

type Data struct {
	Agencies      []Agency       `json:"agencies"`
	Calendars     []Calendar     `json:"calendars"`
	CalendarDates []CalendarDate `json:"calendar_dates"`
	Routes        []Route        `json:"routes"`
	Shapes        []Shape        `json:"shapes"`
	StopTimes     []StopTime     `json:"stop_times"`
	Stops         []Stop         `json:"stops"`
	Trips         []Trip         `json:"trips"`
}

// agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url
// MARTA,Metropolitan Atlanta Rapid Transit Authority,https://www.itsmarta.com,America/New_York,en,404-848-5000,https://www.itsmarta.com/fare-programs.aspx
type Agency struct {
	Id       string `json:"agency_id"`
	Name     string `json:"agency_name"`
	Url      string `json:"agency_url"`
	Timezone string `json:"agency_timezone"`
	Lang     string `json:"agency_lang"`
	Phone    string `json:"agency_phone"`
	FareUrl  string `json:"agency_fare_url"`
}

// service_id,date,exception_type
// 34,20220530,1
type CalendarDate struct {
	ServiceId     int       `json:"service_id"`
	Date          time.Time `json:"date"`
	ExceptionType int       `json:"exception_type"`
}

// service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date
// 20,0,0,0,0,0,0,0,20220423,20220812
type Calendar struct {
	ServiceId int       `json:"service_id"`
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

// route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color
// 16883,MARTA,1,Marietta Blvd/Joseph E Lowery Blvd,,3,https://itsmarta.com/1.aspx,FF00FF,000000
type Route struct {
	Id        int     `json:"route_id"`
	AgencyId  string  `json:"agency_id"`
	ShortName int     `json:"route_short_name"`
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
	Id       string  `json:"shape_id"`
	Lat      float64 `json:"shape_pt_lat"`
	Lon      float64 `json:"shape_pt_lon"`
	Sequence int     `json:"shape_pt_sequence"`
	Distance float64 `json:"shape_dist_traveled"`
}

// trip_id,arrival_time,departure_time,stop_id,stop_sequence,stop_headsign,pickup_type,drop_off_type,shape_dist_traveled,timepoint
// 7142673, 6:43:00, 6:43:00,27,1,,0,0,,1
type StopTime struct {
	TripId        string  `json:"trip_id"`
	ArrivalTime   string  `json:"arrival_time"`
	DepartureTime string  `json:"departure_time"`
	StopId        string  `json:"stop_id"`
	StopSequence  int     `json:"stop_sequence"`
	StopHeadsign  string  `json:"stop_headsign"`
	PickupType    int     `json:"pickup_type"`
	DropOffType   int     `json:"drop_off_type"`
	ShapeDist     float64 `json:"shape_dist_traveled"`
	Timepoint     int     `json:"timepoint"`
}

// stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station,stop_timezone,wheelchair_boarding
// 27,907933,HAMILTON E HOLMES STATION,70 HAMILTON E HOLMES DR NW & CSX TRANSPORTATION,33.754553,-84.469302,,,,,,1
type Stop struct {
	Id                 string  `json:"stop_id"`
	Code               string  `json:"stop_code"`
	Name               string  `json:"stop_name"`
	Desc               string  `json:"stop_desc"`
	Lat                float64 `json:"stop_lat"`
	Lon                float64 `json:"stop_lon"`
	ZoneId             string  `json:"zone_id"`
	Url                string  `json:"stop_url"`
	LocationType       int     `json:"location_type"`
	ParentStation      string  `json:"parent_station"`
	Timezone           string  `json:"stop_timezone"`
	WheelchairBoarding int     `json:"wheelchair_boarding"`
}

// route_id,service_id,trip_id,trip_headsign,trip_short_name,direction_id,block_id,shape_id,wheelchair_accessible,bikes_allowed
// 17114,2,7142675,BLUE EASTBOUND TO INDIAN CREEK STATION,,0,1075016,100750,0,0
type Trip struct {
	Id           string `json:"trip_id"`
	ServiceId    string `json:"service_id"`
	Headsign     string `json:"trip_headsign"`
	ShortName    string `json:"trip_short_name"`
	DirectionId  int    `json:"direction_id"`
	BlockId      string `json:"block_id"`
	ShapeId      string `json:"shape_id"`
	Wheelchair   int    `json:"wheelchair_accessible"`
	BikesAllowed int    `json:"bikes_allowed"`
}

func (a *Agency) Add(record []string) error {
	a.Id = record[0]
	a.Name = record[1]
	a.Url = record[2]
	a.Timezone = record[3]
	a.Lang = record[4]
	a.Phone = record[5]
	a.FareUrl = record[6]
	return nil
}

func (c *Calendar) Add(record []string) error {
	c.ServiceId, _ = strconv.Atoi(record[0])
	c.Monday, _ = strconv.Atoi(record[1])
	c.Tuesday, _ = strconv.Atoi(record[2])
	c.Wednesday, _ = strconv.Atoi(record[3])
	c.Thursday, _ = strconv.Atoi(record[4])
	c.Friday, _ = strconv.Atoi(record[5])
	c.Saturday, _ = strconv.Atoi(record[6])
	c.Sunday, _ = strconv.Atoi(record[7])
	c.StartDate, _ = time.Parse("20060102", record[8])
	c.EndDate, _ = time.Parse("20060102", record[9])

	return nil
}

func (c *CalendarDate) Add(record []string) error {
	c.ServiceId, _ = strconv.Atoi(record[0])
	c.Date, _ = time.Parse("20060102", record[1])
	c.ExceptionType, _ = strconv.Atoi(record[2])
	return nil
}

func (r *Route) Add(record []string) error {
	r.Id, _ = strconv.Atoi(record[0])
	r.AgencyId = record[1]
	r.ShortName, _ = strconv.Atoi(record[2])
	r.LongName = record[3]
	r.Desc = record[4]
	r.Type, _ = strconv.Atoi(record[5])
	r.Url = record[6]

	if color, err := hex.DecodeString(record[7]); err == nil {
		r.Color = color
	} else {
		return err
	}

	if textColor, err := hex.DecodeString(record[8]); err == nil {
		r.TextColor = textColor
	} else {
		return err
	}

	return nil
}
