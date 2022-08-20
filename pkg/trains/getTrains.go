package trains

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Direction int

const (
	Unknown Direction = iota
	North
	South
	East
	West
)

var loc *time.Location

type TrainArrival struct {
	Destination string    `json:"DESTINATION"`
	Direction   Direction `json:"DIRECTION"`
	EventTime   time.Time `json:"EVENT_TIME"`
	Line        string    `json:"LINE"`
	NextArrival time.Time `json:"NEXT_ARR"`
	Station     string    `json:"STATION"`
	TrainID     string    `json:"TRAIN_ID"`
}

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

		t.EventTime, err = time.ParseInLocation("1/2/2006 15:04:05 PM", m["EVENT_TIME"].(string), loc)
		if err != nil {
			return err
		}

		waitingSeconds, err := time.ParseDuration(m["WAITING_SECONDS"].(string) + "s")
		if err != nil {
			return err
		}

		t.NextArrival = t.EventTime.Add(waitingSeconds)
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

func GetTrains(url string, log *logrus.Logger) (*[]TrainArrival, error) {
	log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("getting train data")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"body": len(body),
	}).Debug("got train data from http")

	var trains []TrainArrival
	if err := json.Unmarshal(body, &trains); err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"trains": len(trains),
	}).Debug("parsed trains from json")

	return &trains, nil
}
