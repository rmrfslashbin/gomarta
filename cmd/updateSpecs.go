/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

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
	data.Agencies = make([]gtfspec.Agency, 0)
	data.Calendars = make([]gtfspec.Calendar, 0)
	data.CalendarDates = make([]gtfspec.CalendarDate, 0)
	data.Routes = make([]gtfspec.Route, 0)
	data.Shapes = make([]gtfspec.Shape, 0)
	data.StopTimes = make([]gtfspec.StopTime, 0)
	data.Stops = make([]gtfspec.Stop, 0)
	data.Trips = make([]gtfspec.Trip, 0)

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
		agency := &gtfspec.Agency{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := agency.Add(record); err != nil {
				return err
			}
			data.Agencies = append(data.Agencies, *agency)
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed agency.txt")

	case "calendar.txt":
		log.Info("parsing calendar.txt")
		calendar := &gtfspec.Calendar{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := calendar.Add(record); err != nil {
				return err
			}
			data.Calendars = append(data.Calendars, *calendar)
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed calendar.txt")

	case "calendar_dates.txt":
		log.Info("parsing calendar_dates.txt")
		calendarDate := &gtfspec.CalendarDate{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := calendarDate.Add(record); err != nil {
				return err
			}
			data.CalendarDates = append(data.CalendarDates, *calendarDate)
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed calendar_dates.txt")
	case "routes.txt":
		log.Info("parsing routes.txt")
		route := &gtfspec.Route{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := route.Add(record); err != nil {
				return err
			}
			data.Routes = append(data.Routes, *route)
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed routes.txt")

	case "shapes.txt":
		log.Info("parsing shapes.txt")
		shape := &gtfspec.Shape{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := shape.Add(record); err != nil {
				return err
			}
			data.Shapes = append(data.Shapes, *shape)
			count++
		}
		log.WithFields(logrus.Fields{
			"records": count,
			"elapsed": time.Since(startTime),
		}).Debug("parsed shapes.txt")

	case "stop_times.txt":
		log.Info("parsing stop_times.txt")
		stopTime := &gtfspec.StopTime{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := stopTime.Add(record); err != nil {
				return err
			}
			data.StopTimes = append(data.StopTimes, *stopTime)
			count++
		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed stop_times.txt")

	case "stops.txt":
		log.Info("parsing stops.txt")
		stop := &gtfspec.Stop{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := stop.Add(record); err != nil {
				return err
			}
			data.Stops = append(data.Stops, *stop)
			count++
		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed stops.txt")
	case "trips.txt":
		log.Info("parsing trips.txt")
		trip := &gtfspec.Trip{}
		startTime := time.Now()
		count := 0
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := trip.Add(record); err != nil {
				return err
			}
			data.Trips = append(data.Trips, *trip)
			count++
		}
		log.WithFields(logrus.Fields{
			"elapsed": time.Since(startTime),
			"records": count,
		}).Debug("parsed trips.txt")
	}
	return nil
}
