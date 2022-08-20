/*
Copyright © 2022 Robert Sigler <sigler@improvisedscience.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Direction int

const (
	Unknown Direction = iota
	North
	South
	East
	West
)

type TrainArrival struct {
	Destination string    `json:"DESTINATION"`
	Direction   Direction `json:"DIRECTION"`
	EventTime   time.Time `json:"EVENT_TIME"`
	Line        string    `json:"LINE"`
	NextArrival time.Time `json:"NEXT_ARR"`
	Station     string    `json:"STATION"`
	TrainID     string    `json:"TRAIN_ID"`
}

var loc *time.Location

// trainsCmd represents the trip command
var trainsCmd = &cobra.Command{
	Use:   "train",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Catch errors
		var err error
		defer func() {
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("main crashed")
			}
		}()
		if err := getTrains(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(trainsCmd)
	// Atlanta is located in the Eastern time zone.
	loc, _ = time.LoadLocation("America/New_York")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trainsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trainsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTrains() error {
	url := viper.GetString("rail.rest")
	if url == "" {
		return errors.New("rail.rest URL is empty- check the config file")
	}
	log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("getting trip data")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var trains []TrainArrival
	if err := json.Unmarshal(body, &trains); err != nil {
		return err
	}
	for _, train := range trains {
		fmt.Printf("Destination: %s\n", train.Destination)
		fmt.Printf("Direction: %s\n", train.GetDirection())
		fmt.Printf("EventTime: %s\n", train.EventTime.In(loc))
		fmt.Printf("Line: %s\n", train.Line)
		fmt.Printf("NextArrival: %s\n", train.NextArrival.In(loc))
		fmt.Printf("Station: %s\n", train.Station)
		fmt.Printf("TrainID: %s\n", train.TrainID)
		fmt.Printf("Wait time: %s\n", train.NextArrival.Sub(train.EventTime))
		fmt.Println("")
	}

	return nil
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
