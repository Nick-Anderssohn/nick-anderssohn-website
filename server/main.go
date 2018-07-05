package main

import (
	"full-share/server/fullshare"
	"fmt"
	"full-share/server/reaper"
)

func main() {
	server := fullshare.NewServer("0.0.0.0", "8080")
	fmt.Println("Starting reaper...")
	go reaper.ReapEvery48Hours()
	fmt.Println("Starting server...")
	if err := server.Run(); err != nil {
		fmt.Println(err.Error())
	}
}