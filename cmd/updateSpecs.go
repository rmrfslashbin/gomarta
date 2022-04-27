/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateSpecsCmd represents the updateSpecs command
var updateSpecsCmd = &cobra.Command{
	Use:   "updateSpecs",
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
		if err := updateSpecs(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

var data gtfspec.Data

func init() {
	rootCmd.AddCommand(updateSpecsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateSpecsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateSpecsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	updateSpecsCmd.Flags().String("datadir", "", "The directory containing the GTFS files")
	viper.BindPFlag("datadir", updateSpecsCmd.Flags().Lookup("data"))
}

func updateSpecs() error {
	if viper.GetString("datadir") == "" {
		return errors.New("datadir directory not set")
	}

	dataDir, err := filepath.Abs(viper.GetString("datadir"))
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"dataDir": dataDir,
	}).Debug("dataDir")

	data = gtfspec.Data{}
	data.Agencies = make(map[string]*gtfspec.Agency, 0)
	data.Calendars = make(map[int]*gtfspec.Calendar, 0)
	data.CalendarDates = make(map[int64]*gtfspec.CalendarDate, 0)
	data.Routes = make(map[int]*gtfspec.Route, 0)
	data.Shapes = make(map[int64]*gtfspec.Shape, 0)
	data.StopTimes = make(map[int64]*gtfspec.StopTime, 0)
	data.Stops = make(map[int]*gtfspec.Stop, 0)
	data.Trips = make(map[int]*gtfspec.Trip, 0)

	if err := filepath.Walk(dataDir, processFile); err != nil {
		return err
	}

	dataFile := path.Join(dataDir, "data.gob.gz")
	if err := data.Write(dataFile); err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"dataFile": dataFile,
	}).Info("wrote data file")

	return nil
}

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	r := csv.NewReader(fh)
	_, _ = r.Read() // Skip header

	switch info.Name() {
	case "agency.txt":
		log.Info("parsing agency.txt")
		startTime := time.Now()
		count := 0
		for {
			agency := &gtfspec.Agency{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if id, err := agency.Add(record); err != nil {
				return err
			} else {
				data.Agencies[*id] = agency
			}
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed agency.txt")

	case "calendar.txt":
		log.Info("parsing calendar.txt")
		startTime := time.Now()
		count := 0
		for {
			calendar := &gtfspec.Calendar{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if serviceId, err := calendar.Add(record); err != nil {
				return err
			} else {
				data.Calendars[*serviceId] = calendar
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed calendar.txt")

	case "calendar_dates.txt":
		log.Info("parsing calendar_dates.txt")
		startTime := time.Now()
		count := 0
		for {
			calendarDate := &gtfspec.CalendarDate{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if hash, err := calendarDate.Add(record); err != nil {
				return err
			} else {
				fmt.Println(hash)
				spew.Dump(calendarDate)
				data.CalendarDates[*hash] = calendarDate
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed calendar_dates.txt")

	case "routes.txt":
		log.Info("parsing routes.txt")
		startTime := time.Now()
		count := 0
		for {
			route := &gtfspec.Route{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if routeId, err := route.Add(record); err != nil {
				return err
			} else {
				data.Routes[*routeId] = route
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed routes.txt")

	case "shapes.txt":
		log.Info("parsing shapes.txt")
		startTime := time.Now()
		count := 0
		for {
			shape := &gtfspec.Shape{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if shapeId, err := shape.Add(record); err != nil {
				return err
			} else {
				data.Shapes[*shapeId] = shape
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed shapes.txt")

	case "stop_times.txt":
		log.Info("parsing stop_times.txt")
		startTime := time.Now()
		count := 0
		for {
			stopTime := &gtfspec.StopTime{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if hash, err := stopTime.Add(record); err != nil {
				return err
			} else {
				data.StopTimes[*hash] = stopTime
				count++
			}

		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed stop_times.txt")

	case "stops.txt":
		log.Info("parsing stops.txt")
		startTime := time.Now()
		count := 0
		for {
			stop := &gtfspec.Stop{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if stopId, err := stop.Add(record); err != nil {
				return err
			} else {
				data.Stops[*stopId] = stop
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed stops.txt")

	case "trips.txt":
		log.Info("parsing trips.txt")
		startTime := time.Now()
		count := 0
		for {
			trip := &gtfspec.Trip{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if tripId, err := trip.Add(record); err != nil {
				return err
			} else {
				data.Trips[*tripId] = trip
				count++
			}
		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed trips.txt")
	}
	return nil
}
