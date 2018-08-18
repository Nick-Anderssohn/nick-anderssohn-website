package reaper

import (
	"nick-anderssohn-website/full-share/server/file"
	"os"
	"nick-anderssohn-website/full-share/server/db"
	"log"
	"time"
)

func deleteFile(code string) {
	if _, err := os.Stat(file.GetPath(code)); err == nil {
		file.DeleteFile(code)
	}
}

func reap() {
	codesToDelete := db.DeleteFilesOlderThan(2)
	log.Println("codesToDelete: ", codesToDelete)
	for _, code := range codesToDelete {
		deleteFile(code)
	}
}

func ReapEvery48Hours() {
	for {
		reap()
		time.Sleep(48 * time.Hour)
	}
}