package buses

import (
	"io"
	"net/http"

	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// GetTrips gets the current bus trip data from the API.
func GetTrips(url string, log *logrus.Logger) ([]*gtfsrt.FeedEntity, error) {
	log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("getting trip data")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"body": len(body),
	}).Debug("got bus trip data from endpoint")

	feed := &gtfsrt.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"trips": len(feed.GetEntity()),
	}).Debug("parsed bus trips from protobuf")

	return feed.GetEntity(), nil
}
