package api

import (
	"encoding/json"
	"fmt"
	"full-share/server/codegen"
	"full-share/server/file"
	"full-share/server/rest"
	"log"
	"net/http"
)

const (
	FullShareBase = "/fullShare"
	saveFile      = "/saveFile"
)

var generator *codegen.Generator

func init() {
	generator = codegen.NewDefaultGenerator()
	go generator.RunAndWaitForCodeRequests()
}

type SaveFileReq struct {
	Data []byte `json:"data"`
}

type SaveFileResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func SaveFile(writer http.ResponseWriter, req *http.Request) {
	var saveFileReq SaveFileReq
	var resp SaveFileResp

	fmt.Println("In save file handler")

	defer sendJson(writer, &resp)

	if err := json.NewDecoder(req.Body).Decode(&saveFileReq); err != nil {
		log.Println(err.Error())
		resp.Message = "Bad request format."
		return
	}

	var success bool
	var code string
	for !success {
		success = true
		code = atomicGetCode()
		if alreadyExists, err := file.SaveFile(code, saveFileReq.Data); err != nil || alreadyExists {
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

	resp.Success = true
	resp.Code = code
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
	}
}
