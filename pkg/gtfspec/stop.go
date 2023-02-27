package gtfspec

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station,stop_timezone,wheelchair_boarding
// 27,907933,HAMILTON E HOLMES STATION,70 HAMILTON E HOLMES DR NW & CSX TRANSPORTATION,33.754553,-84.469302,,,,,,1
type Stop struct {
	gorm.Model
	StopId             int     `json:"stop_id" gorm:"primaryKey"`
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

func (s *Stop) Add(record []string) error {
	if len(record) != 12 {
		return fmt.Errorf("invalid stop record length: %d", len(record))
	}

	var err error
	s.Name = record[2]
	s.Desc = record[3]
	s.Url = record[7]
	s.ParentStation = record[9]

	s.StopId, err = strconv.Atoi(record[0])
	if err != nil {
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
