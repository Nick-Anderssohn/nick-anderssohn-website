package api

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"full-share/server/db"
	"full-share/server/file"
	"full-share/server/rest"
	"log"
	"net/http"
	"strconv"
	"io/ioutil"
)

const (
	FullShareBase = "/fullShare"
	downloadFile  = "/download"
	domain = "nickanderssohn.com"
)

type SaveFileReq struct {
	Data            []byte `json:"Data"`
	FileSize        int    `json:"FileSize"`
	FileName string `json:"FileName"`
}

type SaveFileResp struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	DownloadLink    string `json:"DownloadLink"`
}

func recoverAndLog() {
	if e := recover(); e != nil {
		log.Println(e)
	}
}

func SaveFile(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("In save file!")
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
	fmt.Println("file name:", saveFileReq.FileName)

	if saveFileReq.Data, err = ioutil.ReadAll(req.Body); err != nil {
		resp.Message = "Could not read request body."
		return
	}

	var success, alreadyExists bool
	var code string

	for !success {
		for !success {
			success = true
			code = getUUID()
			if db.FileInfoExists(code) {
				success = false
			}
		}
		if alreadyExists, _, err = file.SaveFile(code, saveFileReq.FileName, saveFileReq.Data); err != nil || alreadyExists {
			if alreadyExists {
				log.Println("Code already in use: ", code)
				success = false
			} else {
				resp.Message = "Could not save file."
				return
			}
		}
	}

	db.InsertNewFileInfo(code,  saveFileReq.FileName, fileSize)

	resp.Success = true
	resp.DownloadLink = "https://" + domain + downloadFile + "/" + code + "/" + saveFileReq.FileName
}

func sendJson(writer http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	writer.Write(b)
}

func getUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func GetEndpoints() []*rest.Endpoint {
	return []*rest.Endpoint{
		{
			Path:       "/",
			HandleFunc: SaveFile,
		},
	}
}
