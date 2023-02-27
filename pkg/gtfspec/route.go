package gtfspec

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color
// 16883,MARTA,1,Marietta Blvd/Joseph E Lowery Blvd,,3,https://itsmarta.com/1.aspx,FF00FF,000000
type Route struct {
	gorm.Model
	RouteId   int     `json:"route_id" gorm:"primaryKey"`
	AgencyId  string  `json:"agency_id"`
	ShortName string  `json:"route_short_name"`
	LongName  string  `json:"route_long_name"`
	Desc      string  `json:"route_desc"`
	RouteType int     `json:"route_type"`
	Url       string  `json:"route_url"`
	Color     []uint8 `json:"route_color"`
	TextColor []uint8 `json:"route_text_color"`
}

func (r *Route) Add(record []string) error {
	if len(record) != 9 {
		return fmt.Errorf("invalid route record length: %d", len(record))
	}

	var err error

	r.AgencyId = record[1]
	r.ShortName = record[2]
	r.LongName = record[3]
	r.Desc = record[4]
	r.Url = record[6]

	r.RouteId, err = strconv.Atoi(record[0])
	if err != nil {
		return fmt.Errorf("route id: %v", err)
	}
	if r.RouteType, err = strconv.Atoi(record[5]); err != nil {
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
