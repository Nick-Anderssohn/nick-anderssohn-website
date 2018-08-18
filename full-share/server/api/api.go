package api

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"nick-anderssohn-website/full-share/server/db"
	"nick-anderssohn-website/full-share/server/file"
	"nick-anderssohn-website/full-share/server/serverutil"
	"strconv"
)

const (
	FullShareBase = "/fullShare"
	downloadFile  = "/download"
	domain        = "nickanderssohn.com"
)

type SaveFileReq struct {
	Data     []byte `json:"Data"`
	FileSize int    `json:"FileSize"`
	FileName string `json:"FileName"`
}

type SaveFileResp struct {
	Success      bool   `json:"Success"`
	Message      string `json:"Message"`
	DownloadLink string `json:"DownloadLink"`
}

func recoverAndLog() {
	if e := recover(); e != nil {
		log.Println(e)
	}
}

func SaveFile(writer http.ResponseWriter, req *http.Request) {
	var saveFileReq SaveFileReq
	var resp SaveFileResp
	defer recoverAndLog()
	defer sendJson(writer, &resp)

	fileSize, err := strconv.Atoi(req.Header.Get("FileSize"))
	if err != nil {
		resp.Message = "Could not parse file size"
		return
	}

	saveFileReq.FileName = req.Header.Get("FileName")

	if saveFileReq.Data, err = ioutil.ReadAll(req.Body); err != nil {
		resp.Message = "Could not read request body."
		return
	}

	var success, alreadyExists bool
	var code string

	for !success {
		// Grab a code until an unused one is found
		// TODO: remove potential infinite loop
		for !success {
			success = true
			code = getUUID()
			if db.FileInfoExists(code) {
				success = false
			}
		}

		// Save the file. Double check if there is a conflict.
		alreadyExists, _, err = file.SaveFile(code, saveFileReq.FileName, saveFileReq.Data)
		if alreadyExists {
			log.Println("Code already in use: ", code)
			success = false
		} else if err != nil {
			resp.Message = "Could not save file."
			return
		}
	}

	db.InsertNewFileInfo(code, saveFileReq.FileName, fileSize)

	resp.Success = true
	resp.DownloadLink = "https://" + domain + downloadFile + "/" + code + "/" + saveFileReq.FileName
}

// Does not return an error. Instead simply logs.
func sendJson(writer http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
	writer.Write(b)
}

func getUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func GetEndpoints() []*serverutil.Endpoint {
	return []*serverutil.Endpoint{
		{
			Path:       "/",
			HandleFunc: SaveFile,
		},
	}
}
