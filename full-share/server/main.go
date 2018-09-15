package main

import (
	"fmt"
	"nick-anderssohn-website/full-share/server/config"
	"nick-anderssohn-website/full-share/server/fullshare"
	"nick-anderssohn-website/full-share/server/reaper"
)

func main() {
	server := fullshare.NewServer("0.0.0.0", "8080")
	fmt.Println("Starting reaper...")
	if config.ApplicationConfig.RunReaper {
		go reaper.ReapEvery48Hours()
	}
	fmt.Println("Starting server...")
	if err := server.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
