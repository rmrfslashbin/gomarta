package gtfspec

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// trip_id,arrival_time,departure_time,stop_id,stop_sequence,stop_headsign,pickup_type,drop_off_type,shape_dist_traveled,timepoint
// 7142673, 6:43:00, 6:43:00,27,1,,0,0,,1
type StopTime struct {
	gorm.Model
	TripId            int     `json:"trip_id" gorm:"primaryKey"`
	ArrivalTime       string  `json:"arrival_time"`
	DepartureTime     string  `json:"departure_time"`
	StopId            int     `json:"stop_id" gorm:"primaryKey"`
	StopSequence      int     `json:"stop_sequence" `
	StopHeadsign      string  `json:"stop_headsign"`
	PickupType        int     `json:"pickup_type"`
	DropOffType       int     `json:"drop_off_type"`
	ShapeDistTraveled float64 `json:"shape_dist_traveled"`
	Timepoint         int     `json:"timepoint"`
}

func (s *StopTime) Add(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("invalid stop time record length: %d", len(record))
	}

	var err error
	s.ArrivalTime = record[1]
	s.DepartureTime = record[2]
	s.StopHeadsign = record[5]

	s.TripId, err = strconv.Atoi(record[0])
	if err != nil {
		return fmt.Errorf("trip id: %v", err)
	}
	s.StopId, err = strconv.Atoi(record[3])
	if err != nil {
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
