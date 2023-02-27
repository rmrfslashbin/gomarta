package train

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// getTrains gets the current train data from the API.
func GetTrains(url string, log *logrus.Logger) (*[]TrainArrival, error) {
	log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("getting train data")

	// Fetch the data from the API.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the data into a byte array.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"body": len(body),
	}).Debug("got train data from http")

	// Unmarshal the data into a TrainArrival struct.
	var trains []TrainArrival
	if err := json.Unmarshal(body, &trains); err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"trains": len(trains),
	}).Debug("parsed trains from json")

	return &trains, nil
}
