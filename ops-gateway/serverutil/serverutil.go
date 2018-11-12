package serverutil

import (
	"context"
	"net/http"
	"time"
)

// Endpoint wraps a path and a HandlerFunc
type Endpoint struct {
	Path       string
	HandleFunc http.HandlerFunc
}

// Server is a basic rest server
type Server struct {
	server    *http.Server
	Endpoints []*Endpoint
	Context   context.Context
	Address   string
	Port      string
}

// NewServer returns a new server with the endpoints setup to be handled
func NewServer(endpoints []*Endpoint, address, port string) *Server {
	server := &Server{
		server: &http.Server{
			Addr: address + ":" + port,
		},
		Endpoints: endpoints,
		Address:   address,
		Port:      port,
	}

	for _, endpoint := range endpoints {
		http.HandleFunc(endpoint.Path, endpoint.HandleFunc)
	}

	return server
}

// Run runs the server on address:port
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

// ShutDown attempts to gracefully shutdown within 5 seconds
func (s *Server) ShutDown() error {
	ctx, cancel := context.WithTimeout(s.Context, time.Duration(5000))
	defer cancel()
	return s.server.Shutdown(ctx)
}