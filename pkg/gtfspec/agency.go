package gtfspec

import (
	"fmt"

	"gorm.io/gorm"
)

// agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url
// MARTA,Metropolitan Atlanta Rapid Transit Authority,https://www.itsmarta.com,America/New_York,en,404-848-5000,https://www.itsmarta.com/fare-programs.aspx
type Agency struct {
	gorm.Model
	AgencyId string `json:"agency_id" gorm:"primaryKey"`
	Name     string `json:"agency_name"`
	Url      string `json:"agency_url"`
	Timezone string `json:"agency_timezone"`
	Lang     string `json:"agency_lang"`
	Phone    string `json:"agency_phone"`
	FareUrl  string `json:"agency_fare_url"`
}

func (a *Agency) Add(record []string) error {
	if len(record) != 7 {
		return fmt.Errorf("invalid agency record length: %d", len(record))
	}
	a.AgencyId = record[0]
	a.Name = record[1]
	a.Url = record[2]
	a.Timezone = record[3]
	a.Lang = record[4]
	a.Phone = record[5]
	a.FareUrl = record[6]

	return nil
}
