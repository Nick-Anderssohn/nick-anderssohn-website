package fullshare

import (
	"nick-anderssohn-website/full-share/server/api"
	"nick-anderssohn-website/full-share/server/serverutil"
)

type Server struct {
	serverutil.Server
}

func NewServer(address, port string) *Server {
	return &Server{
		Server: *serverutil.NewServer(api.GetEndpoints(), address, port),
	}
}
