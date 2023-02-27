package gtfspec

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
// 100095,33.818860,-84.450519,1,0.0000
type Shape struct {
	gorm.Model
	ShapeId  int     `json:"shape_id" gorm:"primaryKey"`
	Lat      float64 `json:"shape_pt_lat"`
	Lon      float64 `json:"shape_pt_lon"`
	Sequence int     `json:"shape_pt_sequence" gorm:"primaryKey"`
	Distance float64 `json:"shape_dist_traveled"`
}

func (s *Shape) Add(record []string) error {
	if len(record) != 5 {
		return fmt.Errorf("invalid shape record length: %d", len(record))
	}

	var err error

	s.ShapeId, err = strconv.Atoi(record[0])
	if err != nil {
		return fmt.Errorf("shape id: %v", err)
	}

	s.Sequence, err = strconv.Atoi(record[3])
	if err != nil {
		return fmt.Errorf("shape sequence: %v", err)
	}

	if s.Lat, err = strconv.ParseFloat(record[1], 64); err != nil {
		return fmt.Errorf("shape lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[2], 64); err != nil {
		return fmt.Errorf("shape lon: %v", err)
	}

	if s.Distance, err = strconv.ParseFloat(record[4], 64); err != nil {
		return fmt.Errorf("shape distance: %v", err)
	}

	return nil
}
