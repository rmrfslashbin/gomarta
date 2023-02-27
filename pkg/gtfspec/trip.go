package gtfspec

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// route_id,service_id,trip_id,trip_headsign,trip_short_name,direction_id,block_id,shape_id,wheelchair_accessible,bikes_allowed
// 17114,2,7142675,BLUE EASTBOUND TO INDIAN CREEK STATION,,0,1075016,100750,0,0
type Trip struct {
	gorm.Model
	RouteId      int    `json:"route_id" gorm:"primaryKey"`
	ServiceId    int    `json:"service_id"`
	TripID       int    `json:"trip_id" gorm:"primaryKey"`
	Headsign     string `json:"trip_headsign"`
	ShortName    string `json:"trip_short_name"`
	DirectionId  int    `json:"direction_id"`
	BlockId      int    `json:"block_id"`
	ShapeId      int    `json:"shape_id"`
	Wheelchair   bool   `json:"wheelchair_accessible"`
	BikesAllowed bool   `json:"bikes_allowed"`
}

func (t *Trip) Add(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("invalid trip record length: %d", len(record))
	}

	var err error
	t.Headsign = strings.TrimSpace(record[3])
	t.ShortName = strings.TrimSpace(record[4])

	t.RouteId, err = strconv.Atoi(record[0])
	if err != nil {
		return fmt.Errorf("route id: %v", err)
	}
	t.TripID, err = strconv.Atoi(record[2])
	if err != nil {
		return fmt.Errorf("trip id: %v", err)
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

	return nil
}
