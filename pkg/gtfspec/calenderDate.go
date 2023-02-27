package gtfspec

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// service_id,date,exception_type
// 34,20220530,1
type CalendarDate struct {
	gorm.Model
	ServiceId     int       `json:"service_id" gorm:"primaryKey"`
	Date          time.Time `json:"date" gorm:"primaryKey"`
	ExceptionType int       `json:"exception_type"`
}

func (c *CalendarDate) Add(record []string) error {
	if len(record) != 3 {
		return fmt.Errorf("invalid calendar_date record length: %d", len(record))
	}

	var err error
	if c.ServiceId, err = strconv.Atoi(record[0]); err != nil {
		return fmt.Errorf("service_id: %v", err)
	}
	if c.Date, err = time.Parse("20060102", record[1]); err != nil {
		return fmt.Errorf("date: %v", err)
	}
	if c.ExceptionType, err = strconv.Atoi(record[2]); err != nil {
		return fmt.Errorf("exception_type: %v", err)
	}

	return nil
}
