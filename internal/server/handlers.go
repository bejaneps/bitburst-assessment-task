package server

import (
	"fmt"
	"net/http"

	json "github.com/json-iterator/go"

	"github.com/rs/zerolog/log"
)

type callbackReqBody struct {
	ObjectIDs []int `json:"object_ids"`
}

// handleCallback handles all requests coming on /callback route
func (srv *Server) handleCallback(rw http.ResponseWriter, r *http.Request) {
	log.Logger.Info().Msg("received request")
	defer log.Logger.Info().Msg("finished request")

	// unmarshal request body
	var body callbackReqBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Logger.Err(err).Msg("failed to decode request body")
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// do the job in background, so we won't keep busy the client
	go func() {
		log.Logger.Debug().Ints("object_ids", body.ObjectIDs).Msg("request body")

		// send request to tester service and get online statuses for objects
		objStatuses := srv.cli.Do(body.ObjectIDs)
		statuses := ""
		for i, obj := range objStatuses { // output object ids and their statuses for debugging purposes
			if i == len(objStatuses) -1 {
				statuses = statuses + fmt.Sprintf("%d-%t", obj.ID, obj.Online)
				break
			}
			statuses = statuses + fmt.Sprintf("%d-%t, ", obj.ID, obj.Online)
		}

		log.Logger.Debug().Str("statuses", statuses).Msg("object statuses")
	}()

	// notify tester_service that we received objects successfully
	rw.WriteHeader(http.StatusOK)
}
