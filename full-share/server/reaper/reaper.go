package reaper

import (
	"log"
	"nick-anderssohn-website/full-share/server/db"
	"nick-anderssohn-website/full-share/server/file"
	"os"
	"time"
)

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

func ReapEvery48Hours() {
	for {
		reap()
		time.Sleep(48 * time.Hour)
	}
}
