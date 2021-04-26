package client

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	cli := New(&Config{})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	mux := http.NewServeMux()
	// used same code from tester_service
	mux.HandleFunc("/objects/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rng.Int63n(4000)+300) * time.Millisecond)

		idRaw := strings.TrimPrefix(r.URL.Path, "/objects/")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		w.Write([]byte(fmt.Sprintf(`{"id":%d,"online":%v}`, id, id%2 == 0)))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	cli.conf.TesterServiceAddress = srv.URL

	ids := make([]int32, 0, 100)
	for i := 0; i < cap(ids); i++ {
		ids = append(ids, int32(i))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	objs := cli.Do(ctx, ids)

	require.Equal(t, len(ids), len(objs), "length of received objects isn't equal to actual ids")

	// check if client received ids and their online are correct
	for _, obj := range objs {
		if obj.ID % 2 == 0 {
			require.Equalf(t, true, obj.Online, "% id has online false, when it should be true", obj.ID)
		} else {
			require.Equal(t, false, obj.Online, "% id has online true, when it should be false")
		}
	}
}