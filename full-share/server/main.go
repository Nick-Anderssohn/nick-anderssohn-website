package main

import (
	"fmt"
	"nick-anderssohn-website/full-share/server/config"
	"nick-anderssohn-website/full-share/server/fullshare"
	"nick-anderssohn-website/full-share/server/reaper"
)

func main() {
	if config.ApplicationConfig.RunReaper {
		fmt.Println("Starting reaper...")
		go reaper.ReapEvery24Hours()
	}

	server := fullshare.NewServer("0.0.0.0", "8080")
	fmt.Println("Starting server...")
	err := server.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
