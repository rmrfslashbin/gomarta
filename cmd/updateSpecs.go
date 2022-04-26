/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

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

func init() {
	rootCmd.AddCommand(updateSpecsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateSpecsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateSpecsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	updateSpecsCmd.Flags().String("data", "", "The directory containing the GTFS files")
	viper.BindPFlag("data", updateSpecsCmd.Flags().Lookup("data"))
}

func updateSpecs() error {
	dataDir, err := filepath.Abs(viper.GetString("data"))
	if err != nil {
		return err
	}

	data := &gtfspec.Data{}
	data.Agencies = make([]gtfspec.Agency, 0)
	data.Calendars = make([]gtfspec.Calendar, 0)
	data.CalendarDates = make([]gtfspec.CalendarDate, 0)
	data.Routes = make([]gtfspec.Route, 0)
	data.Shapes = make([]gtfspec.Shape, 0)
	data.StopTimes = make([]gtfspec.StopTime, 0)
	data.Stops = make([]gtfspec.Stop, 0)
	data.Trips = make([]gtfspec.Trip, 0)

	data.Trips = make([]gtfspec.Trip, 0)

	if err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
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
			log.Debug("Parsing agency.txt")
			agency := &gtfspec.Agency{}
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
			}

		case "calendar.txt":
			log.Debug("Parsing calendar.txt")
			calendar := &gtfspec.Calendar{}
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
			}

		case "calendar_dates.txt":
			log.Debug("Parsing calendar_dates.txt")
			calendarDate := &gtfspec.CalendarDate{}
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
			}
		case "routes.txt":
			log.Debug("Parsing routes.txt")
			route := &gtfspec.Route{}
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
			}

		case "shapes.txt":
			log.Debug("Parsing shapes.txt")
		case "stop_times.txt":
			log.Debug("Parsing stop_times.txt")
		case "stops.txt":
			log.Debug("Parsing stops.txt")
		case "trips.txt":
			log.Debug("Parsing trips.txt")
		}
		return nil
	}); err != nil {
		return err
	}
	spew.Dump(data)

	return nil
}
