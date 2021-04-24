package db

import (
	"bitburst-assessment-task/internal/db/objects"
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ProcessObjects inserts objects in database if they don't exist,
// else it updates it's online status and last_seen date. Finally, it deletes objects that have last_seen < 30 seconds.
func (db *DB) ProcessObjects(ctx context.Context, ids []int32, onlines []bool) (modifiedIDs []int32, deletedIDs []int32, err error) {
	// check if database is alive before starting transaction
	if err := db.sqlDB.PingContext(ctx); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to ping database")
	}
	
	// acquire db connection from pool
	conn, err := db.sqlDB.Conn(ctx)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to acquire connection from pool")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Logger.Warn().Msg("failed to release connection to pool")
		}
	}()

	// start tx
	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to begin transaction")
	}
	defer func() { // rollback tx on error
		if err != nil {
			if terr := tx.Rollback(); terr != nil {
				log.Logger.Warn().Ints32("ids", ids).Msg("failed to rollback transaction")
			}
		}
	}()
	txQ := db.q.WithTx(tx) // attach queries in tx

	modifiedIDs, err = txQ.InsertObjectsOrUpdate(ctx, objects.InsertObjectsOrUpdateParams{
		Column1: ids,
		Column2: onlines,
	})
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to insert/update objects")
	}

	deletedIDs, err = txQ.DeleteNotSeenObjects(ctx)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to delete not seen objects")
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to commit transaction")
	}
	
	return modifiedIDs, deletedIDs, nil
}