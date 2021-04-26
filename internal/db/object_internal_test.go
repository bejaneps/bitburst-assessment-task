package db

import (
	"context"
	"database/sql"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// startDatabase starts postgres inside container and returns it's connection url,
// for use in tests
func startDatabase(tb testing.TB, zlog *zerolog.Logger) *Config {
	tb.Helper()

	connURL := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword("postgres", "postgres"),
		Path:   "postgres",
	}
	connURL.Query().Add("sslmode", "disable")

	pool, err := dockertest.NewPool("")
	require.Nil(tb, err, "failed to connect to docker")

	psw, _ := connURL.User.Password()
	envs := []string{
		"POSTGRES_USER=" + connURL.User.Username(),
		"POSTGRES_PASSWORD=" + psw,
		"POSTGRES_DB=" + connURL.Path,
	}

	resource, err := pool.Run("postgres", "13-alpine", envs)
	require.Nil(tb, err, "failed to start postgres container")
	// release container resource when job is done on behalf of it
	tb.Cleanup(func() {
		assert.Nil(tb, pool.Purge(resource), "failed to purge postgres container")
	})

	// set container host to connection url host
	connURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		connURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	// attach to container zerolog logger
	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: zlog,
		ErrorStream:  zlog,
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	require.Nil(tb, err, "failed to connect to postgres container zerolog logger")
	// close logger when job is done on behalf of container
	tb.Cleanup(func() {
		if assert.Nil(tb, logWaiter.Close(), "failed to close container zerolog logger") {
			assert.Nil(tb, logWaiter.Wait(), "failed to wait for posgres container logger to close")
		}
	})

	// set a retry function
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := sql.Open("pgx", connURL.String())
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()

		return db.Ping()
	})
	require.Nil(tb, err, "failed to connect to postgres container")

	dbConf := &Config{
		Host:             strings.Split(connURL.Host, ":")[0],
		Port:             strings.Split(connURL.Host, ":")[1],
		Username:         connURL.User.Username(),
		Password:         psw,
		Name:             connURL.Path,
		MigrationVersion: 3,
		SSLmode:          "disable",
	}

	return dbConf
}

func TestInsertObjectsOrUpdate(t *testing.T) {
	t.Parallel()

	zlog := log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}).With().Logger().Level(zerolog.DebugLevel)

	database, err := New(startDatabase(t, &zlog), &zlog)
	require.Nil(t, err, "failed to establish connection")
	t.Cleanup(func() {
		assert.Nil(t, database.Close(), "failed to close connection")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("default", func(t *testing.T) {
		onlineIDs := make([]int32, 0, 100)
		offlineIDs := make([]int32, 0, 100)
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				onlineIDs = append(onlineIDs, int32(i))
			} else {
				offlineIDs = append(offlineIDs, int32(i))
			}
		}

		insertedIDs, updatedIDs, err := database.InsertObjectsOrUpdate(ctx, onlineIDs, offlineIDs)
		require.Nil(t, err, "failed to process objects")

		assert.True(t, len(insertedIDs) == len(onlineIDs), "length of inserted ids isn't equal to len of online ids")
		assert.Equal(t, len(updatedIDs), 0, "length of updated ids isn't equal to 0")

		// compare if two initial ids are equal to modified ids
		onlineIDsMap := make(map[int32]struct{}, 100)
		for _, id := range onlineIDs {
			onlineIDsMap[id] = struct{}{}
		}
		for _, id := range insertedIDs {
			_, ok := onlineIDsMap[id]
			assert.Truef(t, ok, "%d id exists in inserted id slice and not in online id slice", id)
		}

		offlineIDsMap := make(map[int32]struct{}, 100)
		for _, id := range offlineIDs {
			offlineIDsMap[id] = struct{}{}
		}
		for _, id := range updatedIDs {
			_, ok := offlineIDsMap[id]
			assert.Truef(t, ok, "%d id exists in updated id slice and not in offline id slice", id)
		}

		// query database and get all objects
		rows, err := database.sqlDB.QueryContext(ctx, `SELECT o_id, online FROM bitburst."objects";`)
		require.Nil(t, err, "failed to select o_id, online")
		defer func() {
			assert.Nil(t, rows.Close(), "failed to close result rows")
		}()

		type oidonline struct {
			OID    int32 `db:"o_id"`
			Online bool  `db:"online"`
		}
		rs := make([]*oidonline, 0, 100)
		for rows.Next() {
			r := &oidonline{}
			err := rows.Scan(&r.OID, &r.Online)
			require.Nil(t, err, "failed to scan result rows")
			rs = append(rs, r)
		}
		require.Nil(t, rows.Err(), "failed to read result rows")

		// check if all online ids were saved correct in database
		for _, r := range rs {
			assert.True(t, r.Online == true, "%d id has online false, when it should be true")
		}
	})
}

func TestDeleteNotSeenObjects(t *testing.T) {
	t.Parallel()

	zlog := log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}).With().Logger().Level(zerolog.DebugLevel)

	database, err := New(startDatabase(t, &zlog), &zlog)
	require.Nil(t, err, "failed to establish connection")
	t.Cleanup(func() {
		assert.Nil(t, database.Close(), "failed to close connection")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("default", func(t *testing.T) {
		onlineIDs := make([]int32, 0, 100)
		offlineIDs := make([]int32, 0, 100)
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				onlineIDs = append(onlineIDs, int32(i))
			} else {
				offlineIDs = append(offlineIDs, int32(i))
			}
		}

		insertedIDs, updatedIDs, err := database.InsertObjectsOrUpdate(ctx, onlineIDs, offlineIDs)
		require.Nil(t, err, "failed to process objects")

		assert.True(t, len(insertedIDs) == len(onlineIDs), "length of inserted ids isn't equal to len of online ids")
		assert.Equal(t, len(updatedIDs), 0, "length of updated ids isn't equal to 0")

		// compare if two initial ids are equal to modified ids
		onlineIDsMap := make(map[int32]struct{}, 100)
		for _, id := range onlineIDs {
			onlineIDsMap[id] = struct{}{}
		}
		for _, id := range insertedIDs {
			_, ok := onlineIDsMap[id]
			assert.Truef(t, ok, "%d id exists in inserted id slice and not in online id slice", id)
		}

		offlineIDsMap := make(map[int32]struct{}, 100)
		for _, id := range offlineIDs {
			offlineIDsMap[id] = struct{}{}
		}
		for _, id := range updatedIDs {
			_, ok := offlineIDsMap[id]
			assert.Truef(t, ok, "%d id exists in updated id slice and not in offline id slice", id)
		}

		// start object deleter and wait 31 seconds to check if objects got deleted or no
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		go database.DeleteNotSeenObjects(newCtx)

		time.Sleep(31 * time.Second)

		rows, err := database.sqlDB.QueryContext(ctx, `SELECT o_id FROM bitburst."objects"`)
		require.Nil(t, err, "failed to select o_id")
		defer func() {
			assert.Nil(t, rows.Close(), "failed to close result rows")
		}()

		for rows.Next() {
			var rs interface{}
			require.Nil(t, rows.Scan(&rs), "failed to scan result rows")
			require.Nil(t, rs, "objects didn't get deleted from database after 30 seconds")
		}
	})
}

func BenchmarkInsertObjectsOrUpdate(b *testing.B) {
	zlog := log.Output(zerolog.ConsoleWriter{
		Out:        io.Discard,
		TimeFormat: time.Stamp,
	}).With().Logger().Level(zerolog.DebugLevel)

	database, err := New(startDatabase(b, &zlog), &zlog)
	require.Nil(b, err, "failed to establish connection")
	b.Cleanup(func() {
		assert.Nil(b, database.Close(), "failed to close connection")
	})

	ctx, cancel := context.WithCancel(context.Background())
	b.Cleanup(cancel)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// used same code from tester_service
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		// using map to remove duplicates
		var (
			onlineIDs  []int32
			offlineIDs []int32
			idsMap     map[int32]struct{}
		)
		for pb.Next() {
			idsLen := rng.Int31n(200)
			onlineIDs = make([]int32, 0, idsLen)
			offlineIDs = make([]int32, 0, idsLen)
			idsMap = make(map[int32]struct{}, idsLen)
			for i := 0; i < len(idsMap); i++ {
				num := rng.Int31() % 100
				_, ok := idsMap[num]
				if !ok {
					idsMap[num] = struct{}{}
					if num%2 == 0 {
						onlineIDs = append(onlineIDs, num)
					} else {
						offlineIDs = append(offlineIDs, num)
					}
				}
			}
		}

		_, _, err := database.InsertObjectsOrUpdate(ctx, onlineIDs, offlineIDs)
		require.Nil(b, err, "failed to insert or update objects")
	})
}
