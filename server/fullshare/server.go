package fullshare

import (
	"full-share/server/rest"
	"full-share/server/api"
)

type Server struct {
	rest.Server
}

func NewServer(address, port string) *Server {
	return &Server{
		Server: *rest.NewServer(api.GetEndpoints(), address, port),
	}
}