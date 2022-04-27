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

	sqliteDBFile := filepath.Join(dataDir, "gtfs.db")
	os.Remove(sqliteDBFile)

	db, err = gorm.Open(sqlite.Open(sqliteDBFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	db.AutoMigrate(
		&gtfspec.Agency{},
		&gtfspec.Calendar{},
		&gtfspec.CalendarDate{},
		&gtfspec.Route{},
		&gtfspec.Shape{},
		&gtfspec.StopTime{},
		&gtfspec.Stop{},
		&gtfspec.Trip{},
	)

	if err := filepath.Walk(dataDir, processFile); err != nil {
		return err
	}

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
			if err := agency.Add(record); err != nil {
				return err
			} else {
				db.Create(agency)
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
			if err := calendar.Add(record); err != nil {
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
		for {
			calendarDate := &gtfspec.CalendarDate{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := calendarDate.Add(record); err != nil {
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
		for {
			route := &gtfspec.Route{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := route.Add(record); err != nil {
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
		startTime := time.Now()
		s := NewShapes()
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
			if err := shape.Add(record); err != nil {
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
		for {
			stopTime := &gtfspec.StopTime{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := stopTime.Add(record); err != nil {
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
		for {
			stop := &gtfspec.Stop{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := stop.Add(record); err != nil {
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
		for {
			trip := &gtfspec.Trip{}
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if err := trip.Add(record); err != nil {
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
	return nil
}
