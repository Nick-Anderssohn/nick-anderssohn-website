package reaper

import (
	"log"
	"nick-anderssohn-website/full-share/server/db"
	"nick-anderssohn-website/full-share/server/file"
	"os"
	"time"
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
	return
}

func reap() {
	codesToDelete, err := db.DeleteFilesOlderThan(2)
	if err != nil {
		log.Println("could not delete file entries from db ", err)
	}

	for _, code := range codesToDelete {
		err = deleteFile(code)
		if err != nil {
			log.Println("could not delete file ", err)
		}
	}
}
