package fullshare

import (
	"nick-anderssohn-website/full-share/server/api"
	"nick-anderssohn-website/full-share/server/serverutil"
)

// Server is a simple http server with all the endpoints for nickanderssohn.com
type Server struct {
	serverutil.Server
}

// NewServer returns a new server setup with full share's endpoints.
func NewServer(address, port string) *Server {
	return &Server{
		Server: *serverutil.NewServer(api.GetEndpoints(), address, port),
	}
}
