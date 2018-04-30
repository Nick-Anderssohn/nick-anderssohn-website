package fullshare

import (
	"full-share/server/api"
	"full-share/server/rest"
)

type Server struct {
	rest.Server
}

func NewServer(address, port string) *Server {
	return &Server{
		Server: *rest.NewServer(api.GetEndpoints(), address, port),
	}
}
