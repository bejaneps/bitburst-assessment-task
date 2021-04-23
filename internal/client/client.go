package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	TesterServiceAddress string `mapstructure:"tester_service_address"`
}

type Client struct {
	c *http.Client

	conf *Config
}

func New(conf *Config) *Client {
	cli := &Client{}

	cli.c = &http.Client{
		// max possible response time of tester_service is 4s, but I decided to give 1 more second,
		// because request round-trip also adds time time for request
		Timeout: 5 * time.Second,
	}

	cli.conf = conf

	return cli
}

// ObjectsRespBody is a response from tester_service /objects/:id route
type ObjectsRespBody struct {
	ID     int
	Online bool
}

// Do sends a list of object ids to tester service concurrently,
// and gets their online statuses.
func (cli *Client) Do(objectIDs []int) []*ObjectsRespBody {
	objStatuses := make([]*ObjectsRespBody, 0, len(objectIDs))
	objStatusesChan := make(chan *ObjectsRespBody, len(objectIDs))
	go func() { // receive object ids from goroutines
		for obj := range objStatusesChan {
			objStatuses = append(objStatuses, obj)
		}
	}()

	// send requests to get object statuses concurrently,
	// because /objects/ route has unpredictable response time,
	// and if we do it in 1 loop, then it will be unacceptably slow
	wg := &sync.WaitGroup{}
	for _, v := range objectIDs {
		wg.Add(1)
		go func(id int, objStatusesChan chan *ObjectsRespBody, wg *sync.WaitGroup) {
			defer wg.Done()

			// get the object status
			resp, err := cli.c.Get(fmt.Sprintf("http://%s/objects/%d", cli.conf.TesterServiceAddress, id))
			if err != nil {
				// usually client requests default to timeout errors,
				// if it's not the case then report the error
				urlErr := err.(*url.Error)
				if urlErr.Timeout() {
					log.Logger.Warn().Err(err).Int("id", id).Msg("failed to get object status due to timeout")
				} else {
					log.Logger.Warn().Err(err).Int("id", id).Msg("failed to get object status due to unknown reason")
				}
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Logger.Warn().Err(err).Msg("failed to close response body")
				}
			}()

			// decode object status
			objStatus := ObjectsRespBody{}
			if err := json.NewDecoder(resp.Body).Decode(&objStatus); err != nil {
				log.Logger.Warn().Err(err).Msg("failed to decode response body")
			}

			objStatusesChan <- &objStatus
		}(v, objStatusesChan, wg)
	}

	// wait for all requests to be processesed
	wg.Wait()
	close(objStatusesChan)

	return objStatuses
}