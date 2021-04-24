package server

import (
	"context"
	"net/http"
	"time"

	json "github.com/json-iterator/go"

	"github.com/rs/zerolog/log"
)

type callbackReqBody struct {
	ObjectIDs []int32 `json:"object_ids"`
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
	// and miss any callback
	go func() {
		log.Logger.Debug().Ints32("object_ids", body.ObjectIDs).Msg("request body")

		// send request to tester service and get online statuses for objects
		objStatuses := srv.cli.Do(body.ObjectIDs)
		
		// process only unique ids and store them in array for sql query
		uniqueIDs := make(map[int32]struct{})
		ids := make([]int32, 0, len(objStatuses))
		onlines := make([]bool, 0, len(objStatuses))
		for _, v := range objStatuses {
			_, ok := uniqueIDs[v.ID]
			if !ok {
				uniqueIDs[v.ID] = struct{}{}
				ids = append(ids, v.ID)
				onlines = append(onlines, v.Online)
			}
		}

		// insert/update and delete objects
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		modifiedIDs, deletedIDs, err := srv.database.ProcessObjects(ctx, ids, onlines)
		if err != nil {
			log.Logger.Err(err).Msg("failed to process objects in database")
			return
		}

		log.Logger.Debug().Ints32("modified_ids", modifiedIDs).Ints32("deleted_ids", deletedIDs).Msg("succeeded to process objects in database")
	}()

	// notify tester_service that we received objects successfully
	rw.WriteHeader(http.StatusOK)
}
