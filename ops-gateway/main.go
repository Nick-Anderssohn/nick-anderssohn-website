package main

import (
	"log"
	"nick-anderssohn-website/ops-gateway/auth"
	"nick-anderssohn-website/ops-gateway/serverutil"
)

func main() {
	log.Println("Enter main")
	endpoints := []*serverutil.Endpoint{
		auth.LoginPageEndpoint,
		auth.LoginEndpoint,
		auth.KibanaEndpoint,
	}
	log.Println("New server")
	server := serverutil.NewServer(endpoints, "0.0.0.0", "8080")

	log.Println("starting server")
	server.Run()
}