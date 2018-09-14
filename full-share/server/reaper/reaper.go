package reaper

import (
	"nick-anderssohn-website/full-share/server/file"
	"os"
	"nick-anderssohn-website/full-share/server/db"
	"time"
	"log"
)

func deleteFile(code string) {
	if _, err := os.Stat(file.GetPath(code)); err == nil {
		if err = file.DeleteFile(code); err != nil {
			log.Println(err.Error())
		}
	}
}

func reap() {
	codesToDelete := db.DeleteFilesOlderThan(2)
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