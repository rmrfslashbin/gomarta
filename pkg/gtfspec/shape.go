package gtfspec

import (
	"fmt"
	"strconv"
)

// shape_id,shape_pt_lat,shape_pt_lon,shape_pt_sequence,shape_dist_traveled
// 100095,33.818860,-84.450519,1,0.0000
type Shape struct {
	//ShapeId  int     `json:"shape_id"`
	Lat float64 `json:"shape_pt_lat"`
	Lon float64 `json:"shape_pt_lon"`
	//Sequence int     `json:"shape_pt_sequence"`
	Distance float64 `json:"shape_dist_traveled"`
}

type ShapeData struct {
	ShapeId  int
	Sequence int
}

func (s *Shape) Add(record []string) (*ShapeData, error) {
	if len(record) != 5 {
		return nil, fmt.Errorf("invalid shape record length: %d", len(record))
	}

	var err error

	shapeId, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("shape id: %v", err)
	}

	sequence, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, fmt.Errorf("shape sequence: %v", err)
	}

	if s.Lat, err = strconv.ParseFloat(record[1], 64); err != nil {
		return nil, fmt.Errorf("shape lat: %v", err)
	}
	if s.Lon, err = strconv.ParseFloat(record[2], 64); err != nil {
		return nil, fmt.Errorf("shape lon: %v", err)
	}

	if s.Distance, err = strconv.ParseFloat(record[4], 64); err != nil {
		return nil, fmt.Errorf("shape distance: %v", err)
	}

	return &ShapeData{
		ShapeId:  shapeId,
		Sequence: sequence,
	}, nil
}
