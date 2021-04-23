package server

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Config holds configuration for server that is needed
// to set it up and run.
type Config struct {
	ListenAddress   string        `mapstructure:"listen_address"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// Server is a struct that holds http.Server and other dependencies of the app.
type Server struct {
	httpServer *http.Server

	conf *Config
}

// New constructs new server instance.
func New(conf *Config) *Server {
	srv := &Server{}

	srv.httpServer = &http.Server{
		Addr:         conf.ListenAddress,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Handler:      srv.newMux(),
	}

	srv.conf = conf

	return srv
}

// Start spins up a server in a separate goroutine and serves incoming requests.
// Start is blocking function, so run it in a separate goroutine or it will block execution of your code.
func (srv *Server) Start(errChan chan<- error) {
	if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errChan <- errors.WithMessage(err, "failed to listen and serve")
	}
}

// Close closes the server.
func (srv *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.conf.ShutdownTimeout)
	defer cancel()

	// gracefully shutdown the server
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		return errors.WithMessage(err, "failed to shutdown server gracefully")
	}

	return nil
}
