package reaper

import (
	"nick-anderssohn-website/full-share/server/db"
	"nick-anderssohn-website/full-share/server/file"
	"nick-anderssohn-website/full-share/server/slog"
	"os"
	"time"

	"github.com/Nick-Anderssohn/sherlog"
)

// ReapEvery24Hours will delete files older than 2 days every 24 hours
func ReapEvery24Hours() {
	for {
		reap()
		time.Sleep(24 * time.Hour)
	}
}

func deleteFile(code string) (err error) {
	_, err = os.Stat(file.GetPath(code))
	if err != nil {
		return err
	}

	err = file.DeleteFile(code)
	return sherlog.AsError("could not delete file: ", err)
}

func reap() {
	codesToDelete, err := db.DeleteFilesOlderThan(2)
	if err != nil {
		slog.Logger.Log(err)
	}

	for _, code := range codesToDelete {
		err = deleteFile(code)
		if err != nil {
			slog.Logger.Log(err)
		}
	}
}
