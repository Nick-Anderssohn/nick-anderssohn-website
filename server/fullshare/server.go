package fullshare

import (
	"full-share/server/api"
	"full-share/server/rest"
	"net/http"
)

var imageFileServer http.Handler

func init() {
	imageFileServer = http.StripPrefix(api.FullShareBase+"/data/", http.FileServer(http.Dir("data")))
}

type Server struct {
	rest.Server
}

func NewServer(address, port string) *Server {
	return &Server{
		Server: *rest.NewServer(getEndpoints(), address, port),
	}
}

func getEndpoints() []*rest.Endpoint {
	endpoints := []*rest.Endpoint{
		{
			Path:       api.FullShareBase + "/data/",
			HandleFunc: imageFileServer.ServeHTTP,
		},
	}

	return append(endpoints, api.GetEndpoints()...)
}
