package buses

import (
	"io"
	"net/http"

	"github.com/rmrfslashbin/gomarta/pkg/gtfsrt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Input struct {
	Url string
	Log *zerolog.Logger
}

// GetData gets the current bus data from the API.
func GetData(input *Input) ([]*gtfsrt.FeedEntity, error) {
	log.Debug().
		Str("url", input.Url).
		Str("function", "pkg/buses.GetData()").
		Msg("getting requested data")

	resp, err := http.Get(input.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("url", input.Url).
		Str("function", "pkg/buses.GetData()").
		Str("body", string(body)).
		Msg("fetched data from API")

	feed := &gtfsrt.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	log.Debug().
		Str("url", input.Url).
		Str("function", "pkg/buses.GetData()").
		Int("itmes", len(feed.GetEntity())).
		Msg("unmarshalled protobuf data")

	return feed.GetEntity(), nil
}
