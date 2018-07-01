package main

import (
	"full-share/server/fullshare"
	"fmt"
)

func main() {
	server := fullshare.NewServer("0.0.0.0", "8080")
	fmt.Println("Starting server...")
	if err := server.Run(); err != nil {
		fmt.Println(err.Error())
	}
}