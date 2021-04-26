package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// startTx starts postgres transaction, don't forget to close connection and transaction when the work is done
func (db *DB) startTx(ctx context.Context) (*sql.Conn, *sql.Tx, error) {
	// check if database is alive before starting transaction
	if err := db.sqlDB.PingContext(ctx); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to ping database")
	}

	// acquire db connection from pool
	conn, err := db.sqlDB.Conn(ctx)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to acquire connection from pool")
	}

	// start tx
	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		if err := conn.Close(); err != nil {
			db.logger.Warn().Err(err).Msg("failed to release connection to pool")
		}
		return nil, nil, errors.WithMessage(err, "failed to begin transaction")
	}

	return conn, tx, nil
}

// DeleteNotSeenObjects deletes objects that have last_seen < 30 seconds,
// it's to run in background. When a job is done on behalf of this function, context should be canceled.
func (db *DB) DeleteNotSeenObjects(ctx context.Context) {
	tick := time.NewTicker(30 * time.Second)

	subLogger := db.logger.With().Str("func", "DeleteNotSeenObjects").Logger()
	for {
		select {
		case <-ctx.Done():
			tick.Stop()
			return
		case <-tick.C:
			newCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
			defer cancel()

			conn, tx, err := db.startTx(newCtx)
			if err != nil {
				subLogger.Warn().Msgf("%v", err)
				continue
			}
			defer func() { // release connection to pool
				if err := conn.Close(); err != nil {
					subLogger.Warn().Err(err).Msg("failed to release connection to pool")
				}
			}()
			defer func() { // rollback tx on error
				if err != nil {
					if terr := tx.Rollback(); terr != nil {
						subLogger.Warn().Err(terr).Msg("failed to rollback transaction")
					}
				}
			}()
			txQ := db.q.WithTx(tx) // attach queries in tx

			deletedIDs, err := txQ.DeleteNotSeenObjects(newCtx)
			if err != nil {
				subLogger.Warn().Err(err).Msg("failed to delete not seen objects")
				continue
			}

			// commit transaction
			if err := tx.Commit(); err != nil {
				subLogger.Warn().Err(err).Msg("failed to commit transaction")
				continue
			}

			db.logger.Info().Ints32("ids", deletedIDs).Msg("successfully deleted objects from database")
		}
	}
}

// InsertObjectsOrUpdate inserts objects in database if they don't exist,
// else it updates it's online status and last_seen date
func (db *DB) InsertObjectsOrUpdate(ctx context.Context, onlineIDs []int32, offlineIDs []int32) (insertedIDs []int32, updatedIDs []int32, err error) {
	subLogger := db.logger.With().Str("func", "InsertObjectsOrUpdate").Logger()

	conn, tx, err := db.startTx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() { // release connection to pool
		if err := conn.Close(); err != nil {
			subLogger.Warn().Err(err).Msg("failed to release connection to pool")
		}
	}()
	defer func() { // rollback tx on error
		if err != nil {
			if terr := tx.Rollback(); terr != nil {
				subLogger.Warn().Err(terr).Msg("failed to rollback transaction")
			}
		}
	}()
	txQ := db.q.WithTx(tx) // attach queries in tx

	insertedIDs, err = txQ.InsertObjectsOrUpdate(ctx, onlineIDs)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to insert/update objects")
	}

	updatedIDs, err = txQ.UpdateObjects(ctx, offlineIDs)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to update objects")
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, errors.WithMessage(err, "failed to commit transaction")
	}

	return insertedIDs, updatedIDs, nil
}