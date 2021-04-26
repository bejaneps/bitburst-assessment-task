package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type Config struct {
	TesterServiceAddress string `mapstructure:"tester_service_address"`
}

type Client struct {
	c *http.Client

	conf *Config

	logger *zerolog.Logger
}

func New(conf *Config, logger *zerolog.Logger) *Client {
	cli := &Client{}

	cli.c = &http.Client{
		// max possible response time of tester_service is 4s, but I decided to give 1 more second,
		// because request round-trip also adds time time for request
		Timeout: 5 * time.Second,
	}

	cli.conf = conf

	if !strings.Contains(cli.conf.TesterServiceAddress, "http://") {
		if !strings.Contains(cli.conf.TesterServiceAddress, "https://") {
			cli.conf.TesterServiceAddress = "http://" + cli.conf.TesterServiceAddress
		}
	}

	cli.logger = logger

	return cli
}

// ObjectsRespBody is a response from tester_service /objects/:id route
type ObjectsRespBody struct {
	ID     int32
	Online bool
}

// Do sends a list of object ids to tester service concurrently,
// and gets their online statuses.
func (cli *Client) Do(ctx context.Context, objectIDs []int32) []*ObjectsRespBody {
	objStatuses := make([]*ObjectsRespBody, 0, len(objectIDs))
	objStatusesChan := make(chan *ObjectsRespBody, len(objectIDs))
	// receive object ids from goroutines concurrently, so we read ids from buffer at the same time they are sent
	go func() {
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
		go func(id int32, objStatusesChan chan *ObjectsRespBody, wg *sync.WaitGroup) {
			defer wg.Done()

			// get the object status
			resp, err := cli.c.Get(fmt.Sprintf("%s/objects/%d", cli.conf.TesterServiceAddress, id))
			if err != nil {
				// usually client requests default to timeout errors,
				// if it's not the case then report the error
				urlErr := err.(*url.Error)
				if urlErr.Timeout() {
					cli.logger.Warn().Err(err).Int32("id", id).Msg("failed to get object status due to timeout")
				} else {
					cli.logger.Warn().Err(err).Int32("id", id).Msg("failed to get object status due to unknown reason")
				}
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					cli.logger.Warn().Err(err).Msg("failed to close response body")
				}
			}()

			// decode object status
			objStatus := ObjectsRespBody{}
			if err := json.NewDecoder(resp.Body).Decode(&objStatus); err != nil {
				cli.logger.Warn().Err(err).Msg("failed to decode response body")
			}

			objStatusesChan <- &objStatus
		}(v, objStatusesChan, wg)
	}

	// wait for all requests to be processesed
	wg.Wait()

	// wait until the last obj is received from channel and appended to objStatuses slice, or until context timeout is reached
LOOP:
	for {
		select {
		case <-ctx.Done():
			close(objStatusesChan)
			break LOOP
		default:
			if len(objStatuses) == len(objectIDs) {
				close(objStatusesChan)
				break LOOP
			}
		}
	}

	cli.c.CloseIdleConnections()

	return objStatuses
}
