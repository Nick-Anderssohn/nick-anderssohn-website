package main

import (
	"nick-anderssohn-website/full-share/server/config"
	"nick-anderssohn-website/full-share/server/fullshare"
	"nick-anderssohn-website/full-share/server/reaper"
	"nick-anderssohn-website/full-share/server/slog"
)

func main() {
	if config.ApplicationConfig.RunReaper {
		slog.Logger.Info("Starting reaper")
		go reaper.ReapEvery24Hours()
	}

	server := fullshare.NewServer("0.0.0.0", "8080")
	slog.Logger.Info("Starting server")
	err := server.Run()
	if err != nil {
		slog.Logger.Critical(err)
	}
}
