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

		// process only unique ids, so we don't send same id twice or thrice to server, for example if would receive 1,000,000 ids and 1/3 of them would be duplicates, then we would send 333,333 useless requests and waste time
		uniqueIDs := make(map[int32]struct{})
		ids := make([]int32, 0, len(body.ObjectIDs))
		for _, id := range body.ObjectIDs {
			_, ok := uniqueIDs[id]
			if !ok {
				uniqueIDs[id] = struct{}{}
				ids = append(ids, id)
			}
		}

		// timeout can be increased if server responds longer than 4 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		// send request to tester service and get online statuses for objects
		objStatuses := srv.cli.Do(ctx, ids)

		onlineIDs := make([]int32, 0, len(ids))
		offlineIDs := make([]int32, 0, len(ids))
		for _, obj := range objStatuses {
			if obj.Online {
				onlineIDs = append(onlineIDs, obj.ID)
			} else {
				offlineIDs = append(offlineIDs, obj.ID)
			}
		}

		// insert/update and delete objects
		insertedIDs, updatedIDs, err := srv.database.InsertObjectsOrUpdate(ctx, onlineIDs, offlineIDs)
		if err != nil {
			log.Logger.Err(err).Msg("failed to process objects in database")
			return
		}

		log.Logger.Info().Ints32("inserted_ids", insertedIDs).Ints32("updated_ids", updatedIDs).Msg("succeeded to process objects in database")
	}()

	// notify tester_service that we received objects successfully
	rw.WriteHeader(http.StatusOK)
}
