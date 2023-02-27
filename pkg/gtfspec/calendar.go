package gtfspec

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date
// 20,0,0,0,0,0,0,0,20220423,20220812
type Calendar struct {
	gorm.Model
	ServiceId int       `json:"service_id" gorm:"primaryKey"`
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

func (c *Calendar) Add(record []string) error {
	if len(record) != 10 {
		return fmt.Errorf("invalid calendar record length: %d", len(record))
	}

	var err error

	c.ServiceId, err = strconv.Atoi(record[0])
	if err != nil {
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
