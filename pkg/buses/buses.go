package buses

import (
	"io"
	"net/http"

	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// GetData gets the current bus data from the API.
func GetData(url string, log *logrus.Logger) ([]*gtfsrt.FeedEntity, error) {
	log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("getting requested data")

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
	}).Debug("got bus data from endpoint")

	feed := &gtfsrt.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"items": len(feed.GetEntity()),
	}).Debug("parsed bus data from protobuf")

	return feed.GetEntity(), nil
}
