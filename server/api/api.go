package api

import (
	"encoding/json"
	"fmt"
	"full-share/server/codegen"
	"full-share/server/db"
	"full-share/server/file"
	"full-share/server/rest"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	FullShareBase = "/fullShare"
	saveFile      = "/saveFile"
	downloadFile  = "/download"
)

var generator *codegen.Generator

func init() {
	generator = codegen.NewDefaultGenerator()
	go generator.RunAndWaitForCodeRequests()
}

type SaveFileReq struct {
	Data            []byte `json:"data"`
	IsCompressed    bool   `json:"is_compressed"`
	CompressionType string `json:"compression_type"`
	FileSize        int    `json:"file_size"`
}

type SaveFileResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func strToBool(str string) bool {
	return strings.ToLower(str) == "true"
}

func boolToStr(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func SaveFile(writer http.ResponseWriter, req *http.Request) {
	var saveFileReq SaveFileReq
	var resp SaveFileResp
	var err error

	defer sendJson(writer, &resp)

	saveFileReq.IsCompressed = strToBool(req.Header.Get("Is_compressed"))
	saveFileReq.CompressionType = req.Header.Get("Compression_type")
	fileSize, err := strconv.Atoi(req.Header.Get("File_size"))
	if err != nil {
		log.Println(err.Error())
		resp.Message = "Could not parse file size"
	}
	saveFileReq.Data, err = ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println(err.Error())
		resp.Message = "Could not read request body."
		return
	}

	var success, alreadyExists bool
	var code, path string
	for !success {
		for !success {
			success = true
			code = atomicGetCode()
			if db.FileInfoExists(code) {
				success = false
			}
		}
		if alreadyExists, path, err = file.SaveFile(code, saveFileReq.Data); err != nil || alreadyExists {
			if alreadyExists {
				log.Println("Code already in use: ", code)
				success = false
			} else {
				log.Println(err.Error())
				resp.Message = "Could not save file."
				if !alreadyExists {
					generator.RestoreCode(code)
				}
				return
			}
		}
	}

	db.InsertNewFileInfo(code, saveFileReq.CompressionType, path, saveFileReq.IsCompressed, fileSize)

	resp.Success = true
	resp.Code = code
}

func DownloadFile(writer http.ResponseWriter, req *http.Request) {
	code := req.Header.Get("Code")

	if code == "" {
		writer.Header().Set("Success", "false")
		writer.Header().Set("Err_message", "no code provided")
		writer.Write([]byte{0})
		return
	}

	// get file info
	fileInfo := db.SelectFileInfo(code)
	if fileInfo == nil {
		writer.Header().Set("Success", "false")
		writer.Header().Set("Err_message", "There was a problem retrieving the file associated with the provided code. It is possible that the file does not exist.")
		writer.Write([]byte{0})
		return
	}
	fileData, err := file.GetFile(fileInfo.Path)
	if err != nil {
		writer.Header().Set("Success", "false")
		writer.Header().Set("Err_message", "There was a problem retrieving the file associated with the provided code. It is possible that the file does not exist.")
		writer.Write([]byte{0})
		return
	}

	writer.Header().Set("Success", "true")
	writer.Header().Set("Is_compressed", boolToStr(fileInfo.IsCompressed))
	writer.Header().Set("Compression_type", fileInfo.CompressionType)
	writer.Header().Set("Code", fileInfo.Code)
	writer.Header().Set("File_size", strconv.Itoa(fileInfo.FileSize))
	writer.Write(fileData)
}

func sendJson(writer http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	writer.Write(b)
}

func atomicGetCode() string {
	generator.RequestCodeChan <- true
	return <-generator.CodeChan
}

func GetEndpoints() []*rest.Endpoint {
	return []*rest.Endpoint{
		{
			Path:       FullShareBase + saveFile,
			HandleFunc: SaveFile,
		},
		{
			Path:       FullShareBase + downloadFile,
			HandleFunc: DownloadFile,
		},
	}
}
