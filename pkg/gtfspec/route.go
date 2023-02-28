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

func (r *Route) Add(headers map[string]int, record []string) error {
	if len(record) != 9 {
		return fmt.Errorf("invalid route record length: %d", len(record))
	}

	var err error

	r.AgencyId = record[headers["agency_id"]]
	r.ShortName = record[headers["route_short_name"]]
	r.LongName = record[headers["route_long_name"]]
	r.Desc = record[headers["route_desc"]]
	r.Url = record[headers["route_url"]]

	r.RouteId, err = strconv.Atoi(record[headers["route_id"]])
	if err != nil {
		return fmt.Errorf("route id: %v", err)
	}
	if r.RouteType, err = strconv.Atoi(record[headers["route_type"]]); err != nil {
		return fmt.Errorf("route type: %v", err)
	}
	if color, err := hex.DecodeString(record[headers["route_color"]]); err != nil {
		return fmt.Errorf("route color: %v", err)
	} else {
		r.Color = color
	}
	if textColor, err := hex.DecodeString(record[headers["route_text_color"]]); err != nil {
		return fmt.Errorf("route text color: %v", err)
	} else {
		r.TextColor = textColor
	}

	return nil
}
