package train

/*
import (
	"errors"
	"fmt"

	"github.com/rmrfslashbin/gomarta/pkg/trains"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// trainsCmd represents the train command
var trainsCmd = &cobra.Command{
	Use:   "train",
	Short: "Fetch current train data",
	Long:  "Fetch current train arrival data from the Marta API",
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

// init initializes the command.
func init() {
	// Add subcommand to the root command.
	rootCmd.AddCommand(trainsCmd)
}

// getTrains fetches train data from the Marta API.
func getTrains() error {
	// Get URL or bail out.
	url := viper.GetString("rail.rest")
	if url == "" {
		return errors.New("rail.rest URL is empty- check the config file")
	}

	// Get the train data.
	trains, err := trains.GetTrains(url, log)
	if err != nil {
		return err
	}

	// Print the train data.
	for _, train := range *trains {
		fmt.Printf("Destination: %s\n", train.Destination)
		fmt.Printf("Direction: %s\n", train.GetDirection())
		fmt.Printf("EventTime: %s\n", train.EventTime)
		fmt.Printf("Line: %s\n", train.Line)
		fmt.Printf("NextArrival: %s\n", train.NextArrival)
		fmt.Printf("Station: %s\n", train.Station)
		fmt.Printf("TrainID: %s\n", train.TrainID)
		fmt.Printf("Wait time: %s\n", train.NextArrival.Sub(train.EventTime))
		fmt.Println("")
	}

	return nil
}
*/
