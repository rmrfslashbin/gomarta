/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const MAX_RECORDS = 2000

type StopTimes struct {
	stopTimes []*gtfspec.StopTime
}

type Shapes struct {
	shapes []*gtfspec.Shape
}

func NewShapes() *Shapes {
	return &Shapes{
		shapes: make([]*gtfspec.Shape, 0, MAX_RECORDS),
	}
}

func (s *Shapes) Add(shape *gtfspec.Shape) {
	s.shapes = append(s.shapes, shape)
	if len(s.shapes) > MAX_RECORDS {
		db.Create(s.shapes)
		s.shapes = make([]*gtfspec.Shape, 0, MAX_RECORDS)
	}
}

func (s *Shapes) Flush() {
	db.Create(s.shapes)
}

func NewStopTimes() *StopTimes {
	return &StopTimes{
		stopTimes: make([]*gtfspec.StopTime, 0, MAX_RECORDS),
	}
}

func (s *StopTimes) Add(stopTime *gtfspec.StopTime) {
	s.stopTimes = append(s.stopTimes, stopTime)
	if len(s.stopTimes) > MAX_RECORDS {
		db.Create(s.stopTimes)
		s.stopTimes = make([]*gtfspec.StopTime, 0, MAX_RECORDS)
	}
}

func (s *StopTimes) Flush() {
	db.Create(s.stopTimes)
}

// updateSpecsCmd represents the updateSpecs command
var updateSpecsCmd = &cobra.Command{
	Use:   "updateSpecs",
	Short: "Update the GTFS specs",
	Long:  "Update the GTFS specs to include changes made to the GTFS spec such as routes, time tables, routes, etc",
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

// var data gtfspec.Data
var db *gorm.DB

func init() {
	rootCmd.AddCommand(updateSpecsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateSpecsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateSpecsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	updateSpecsCmd.Flags().String("datadir", "", "The directory containing the GTFS database")
	viper.BindPFlag("datadir", updateSpecsCmd.Flags().Lookup("data"))
}

func updateSpecs() error {
	var dataDir string
	var err error

	if viper.GetString("datadir") != "" {
		dataDir, err = filepath.Abs(viper.GetString("datadir"))
		if err != nil {
			return err
		}
	} else {
		userDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		dataDir = path.Join(userDir, "gomarta")
	}

	sqliteDBFile := filepath.Join(dataDir, "gtfs.db")
	log.WithFields(logrus.Fields{
		"sqliteDBFile": sqliteDBFile,
	}).Debug("GTFS Spec sqliteDBFile")
	os.Remove(sqliteDBFile)

	db, err = gorm.Open(sqlite.Open(sqliteDBFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(
		&gtfspec.Agency{},
		&gtfspec.Calendar{},
		&gtfspec.CalendarDate{},
		&gtfspec.Route{},
		&gtfspec.Shape{},
		&gtfspec.StopTime{},
		&gtfspec.Stop{},
		&gtfspec.Trip{},
		&gtfspec.VehiclePosition{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := processZipData(); err != nil {
		return err
	}

	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func processZipData() error {
	specURL := viper.GetString("gtfs.specs")
	if specURL == "" {
		return fmt.Errorf("no spec URL provided")
	}
	resp, err := http.Get(specURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range zipReader.File {
		zipData, err := readZipFile(file)
		if err != nil {
			return err
		}
		csvData, err := csv.NewReader(bytes.NewReader(zipData)).ReadAll()
		if err != nil {
			return err
		}

		switch file.Name {
		case "agency.txt":
			log.Info("parsing agency.txt")
			count := 0
			startTime := time.Now()

			for _, row := range csvData[1:] {
				agency := &gtfspec.Agency{}
				if err := agency.Add(row); err != nil {
					return err
				} else {
					db.Create(agency)
					count++
				}
			}
			log.WithFields(logrus.Fields{
				"records": count,
				"elapsed": time.Since(startTime),
			}).Debug("parsed agency.txt")

		case "calendar.txt":
			log.Info("parsing calendar.txt")
			count := 0
			startTime := time.Now()

			for _, row := range csvData[1:] {
				calendar := &gtfspec.Calendar{}

				if err := calendar.Add(row); err != nil {
					return err
				} else {
					db.Create(calendar)
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

			for _, row := range csvData[1:] {
				calendarDate := &gtfspec.CalendarDate{}

				if err := calendarDate.Add(row); err != nil {
					return err
				} else {
					db.Create(calendarDate)
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

			for _, row := range csvData[1:] {
				route := &gtfspec.Route{}
				if err := route.Add(row); err != nil {
					return err
				} else {
					db.Create(route)
					count++
				}
			}
			log.WithFields(logrus.Fields{
				"records": count,
				"elapsed": time.Since(startTime),
			}).Debug("parsed routes.txt")

		case "shapes.txt":
			log.Info("parsing shapes.txt")
			s := NewShapes()
			startTime := time.Now()
			count := 0

			for _, row := range csvData[1:] {
				shape := &gtfspec.Shape{}
				if err := shape.Add(row); err != nil {
					return err
				} else {
					s.Add(shape)
					count++
				}
			}
			s.Flush()
			log.WithFields(logrus.Fields{
				"records": count,
				"elapsed": time.Since(startTime),
			}).Debug("parsed shapes.txt")

		case "stop_times.txt":
			log.Info("parsing stop_times.txt")
			startTime := time.Now()
			count := 0
			st := NewStopTimes()
			for _, row := range csvData[1:] {
				stopTime := &gtfspec.StopTime{}
				if err := stopTime.Add(row); err != nil {
					return err
				} else {
					st.Add(stopTime)
					count++
				}
			}
			st.Flush()
			log.WithFields(logrus.Fields{
				"elapsed": time.Since(startTime),
				"records": count,
			}).Debug("parsed stop_times.txt")

		case "stops.txt":
			log.Info("parsing stops.txt")
			startTime := time.Now()
			count := 0
			for _, row := range csvData[1:] {
				stop := &gtfspec.Stop{}
				if err := stop.Add(row); err != nil {
					return err
				} else {
					db.Create(stop)
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
			for _, row := range csvData[1:] {
				trip := &gtfspec.Trip{}
				if err := trip.Add(row); err != nil {
					return err
				} else {
					db.Create(trip)
					count++
				}
			}
			log.WithFields(logrus.Fields{
				"elapsed": time.Since(startTime),
				"records": count,
			}).Debug("parsed trips.txt")
		}
	}
	return nil
}
