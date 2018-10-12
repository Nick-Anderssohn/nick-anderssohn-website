package slog

import "github.com/Nick-Anderssohn/sherlog"

var Logger sherlog.Logger

func init() {
	var err error
	Logger, err = sherlog.NewRollingFileLoggerWithSizeLimit("/logs/full-share.log", 1000)
	if err != nil {
		panic(err)
	}
}