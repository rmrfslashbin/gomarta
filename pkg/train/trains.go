package train

import (
	"encoding/json"
	"strings"
	"time"
)

// Train direction.
type Direction int

// Train directions
const (
	Unknown Direction = iota
	North
	South
	East
	West
)

// Var to store the location of the time zone.
var loc *time.Location

// TrainArrival is the data structure for a train arrival.
type TrainArrival struct {
	Destination string    `json:"DESTINATION"`
	Direction   Direction `json:"DIRECTION"`
	EventTime   time.Time `json:"EVENT_TIME"`
	Line        string    `json:"LINE"`
	NextArrival time.Time `json:"NEXT_ARR"`
	Station     string    `json:"STATION"`
	TrainID     string    `json:"TRAIN_ID"`
}

// GetDirection returns the direction of the train.
func (t *TrainArrival) GetDirection() string {
	switch t.Direction {
	case North:
		return "north"
	case South:
		return "south"
	case East:
		return "east"
	case West:
		return "west"
	default:
		return "unknown"
	}
}

// UnmarshalJSON unmarshals the JSON data into the TrainArrival struct.
func (t *TrainArrival) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case map[string]interface{}:
		var err error
		m := v.(map[string]interface{})
		t.Destination = strings.ToLower(m["DESTINATION"].(string))
		t.Line = strings.ToLower(m["LINE"].(string))
		t.Station = strings.ToLower(m["STATION"].(string))
		t.TrainID = strings.ToLower(m["TRAIN_ID"].(string))

		// Parse the event time data; assume Atlanta is in the Eastern time zone.
		t.EventTime, err = time.ParseInLocation("1/2/2006 15:04:05 PM", m["EVENT_TIME"].(string), loc)
		if err != nil {
			return err
		}

		// Marta retuns some very useless data realted to train wait times.
		// Let's use the "WAITING_SECONDS" value to calculate the arrival time of the train.
		waitingSeconds, err := time.ParseDuration(m["WAITING_SECONDS"].(string) + "s")
		if err != nil {
			return err
		}

		// Add the waiting time to the event time to get the arrival time.
		t.NextArrival = t.EventTime.Add(waitingSeconds)

		// Parse the direction.
		switch strings.ToLower(m["DIRECTION"].(string)) {
		default:
			t.Direction = Unknown
		case "n":
			t.Direction = North
		case "s":
			t.Direction = South
		case "e":
			t.Direction = East
		case "w":
			t.Direction = West
		}
	}
	return nil
}

func init() {
	var err error
	// Atlanta is located in the Eastern time zone.
	loc, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
}
