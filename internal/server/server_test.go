package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input *Config
		want  *Server
	}{
		"default": {
			input: &Config{
				ListenAddress: "0.0.0.0:9090",
				ReadTimeout:   time.Second * 2,
				WriteTimeout:  time.Second * 4,
			},
			want: &Server{httpServer: &http.Server{
				Addr:         "0.0.0.0:9090",
				ReadTimeout:  time.Second * 2,
				WriteTimeout: time.Second * 4,
			}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := New(tc.input, nil, nil)

			diff := ""

			diff = cmp.Diff(tc.want.httpServer.Addr, got.httpServer.Addr)
			require.Empty(t, diff, "failed to get same server listen address")

			diff = cmp.Diff(tc.want.httpServer.ReadTimeout, got.httpServer.ReadTimeout)
			require.Empty(t, diff, "failed to get same server read timeout")

			diff = cmp.Diff(tc.want.httpServer.WriteTimeout, got.httpServer.WriteTimeout)
			require.Empty(t, diff, "failed to get same server write timeout")
		})
	}
}

func TestServerStartAndClose(t *testing.T) {
	tests := map[string]struct {
		input struct {
			srv     *Server
			errChan chan error
		}
	}{
		"default": {
			input: struct {
				srv     *Server
				errChan chan error
			}{
				srv: &Server{
					conf: &Config{
						ShutdownTimeout: 2 * time.Second,
					},
					httpServer: &http.Server{
						Addr: "0.0.0.0:9090",
					},
				},
				errChan: make(chan error),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// if server.Start function doesn't return error in 2 seconds,
			// then test is counted as passed.
			// If server.Start function returns error in 2 seconds,
			// then test is counted as failed
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			go tc.input.srv.Start(tc.input.errChan)

			select {
			case <-ctx.Done():
				require.Nil(t, tc.input.srv.Close())
			case err := <-tc.input.errChan:
				require.Fail(t, err.Error(), "failed to start server")
			}
		})
	}
}
