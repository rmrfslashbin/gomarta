/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateSpecsCmd represents the validateSpecs command
var validateSpecsCmd = &cobra.Command{
	Use:   "validateSpecs",
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
		if err := validateSpecs(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateSpecsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateSpecsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateSpecsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validateSpecsCmd.Flags().String("datadir", "", "The directory containing the GTFS files")
	viper.BindPFlag("datadir", updateSpecsCmd.Flags().Lookup("data"))
}

func validateSpecs() error {
	if viper.GetString("datadir") == "" {
		return errors.New("datadir directory not set")
	}

	dataDir, err := filepath.Abs(viper.GetString("datadir"))
	if err != nil {
		return err
	}

	dataFile := path.Join(dataDir, "data.gob.gz")

	data := gtfspec.Data{}
	if err := data.Read(dataFile); err != nil {
		return err
	}

	fmt.Printf("agency: %d\n", len(data.Agencies))
	fmt.Printf("calendar: %d\n", len(data.Calendars))
	fmt.Printf("calendar_dates: %d\n", len(data.CalendarDates))
	fmt.Printf("routes: %d\n", len(data.Routes))
	fmt.Printf("shapes: %d\n", len(data.Shapes))
	fmt.Printf("stop_times: %d\n", len(data.StopTimes))
	fmt.Printf("stops: %d\n", len(data.Stops))
	fmt.Printf("trips: %d\n", len(data.Trips))

	return nil
}
