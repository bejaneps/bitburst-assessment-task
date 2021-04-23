package server

import "net/http"

// newMux constructs new server multiplexer and registers all routes to it
func (srv *Server) newMux() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
