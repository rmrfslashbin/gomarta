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
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/gomarta/pkg/buses"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// busCmd represents the bus command
var busCmd = &cobra.Command{
	Use:   "bus",
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
		if err := getBuses(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(busCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// busCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// busCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getBuses() error {
	userDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dataDir := path.Join(userDir, "gomarta")
	sqliteDBFile := filepath.Join(dataDir, "gtfs.db")
	if _, err := os.Stat(sqliteDBFile); os.IsNotExist(err) {
		return errors.New("sqlite database does not exist")
	}
	log.WithFields(logrus.Fields{
		"sqliteDBFile": sqliteDBFile,
	}).Debug("GTFS Spec sqliteDBFile")
	db, err := gorm.Open(sqlite.Open(sqliteDBFile), &gorm.Config{})
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

	url := viper.GetString("gtfs.bus.vehicles")
	if url == "" {
		return errors.New("vehicles URL is empty- check the config file")
	}

	vehicles, err := buses.GetData(url, log)
	if err != nil {
		return err
	}

	for _, vehicle := range vehicles {
		vp := gtfspec.VehiclePosition{}

		spew.Dump(vehicle)
		routeId, err := strconv.Atoi(*vehicle.Vehicle.Trip.RouteId)
		if err != nil {
			return err
		}
		//route := gtfspec.Route{RouteId: routeId}

		tripId, err := strconv.Atoi(*vehicle.Vehicle.Trip.TripId)
		if err != nil {
			return err
		}

		vp.ID = vehicle.GetId()
		vp.VehicleID = vehicle.GetVehicle().Vehicle.GetId()
		vp.VehicleLabel = vehicle.GetVehicle().Vehicle.GetLabel()
		vp.TripID = tripId
		vp.RouteID = routeId
		vp.TripStartDate = vehicle.GetVehicle().Trip.GetStartDate()
		vp.Position.Bearing = vehicle.GetVehicle().Position.GetBearing()
		vp.Position.Latitude = vehicle.GetVehicle().Position.GetLatitude()
		vp.Position.Longitude = vehicle.GetVehicle().Position.GetLongitude()
		vp.Timestamp = vehicle.GetVehicle().GetTimestamp()
		vp.OccupancyStatus = vehicle.GetVehicle().GetOccupancyStatus().String()
		db.Create(&vp)

		/*
			trip := gtfspec.Trip{TripId: tripId}
			db.First(&route)
			db.First(&trip)
			fmt.Println()
			spew.Dump(route)
			spew.Dump(trip)
		*/
		spew.Dump(vp)
		break

	}

	return nil
}
