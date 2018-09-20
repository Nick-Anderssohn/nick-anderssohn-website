package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"nick-anderssohn-website/full-share/server/db"
	"nick-anderssohn-website/full-share/server/file"
	"nick-anderssohn-website/full-share/server/serverutil"
	"os"
	"strconv"
	"unicode"

	"github.com/satori/go.uuid"
)

const (
	FullShareBase = "/fullShare"
	downloadFile  = "/download"
)

var domain string
var protocol string

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

func init() {
	domain = os.Getenv("WEBSITE_DOMAIN")
	protocol = os.Getenv("WEBSITE_PROTOCOL")
}

func recoverAndLog() {
	if e := recover(); e != nil {
		log.Println(e)
	}
}

// saveFile is an http handler func that saves a file
func saveFile(writer http.ResponseWriter, req *http.Request) {
	var saveFileReq SaveFileReq
	var resp SaveFileResp
	defer recoverAndLog()
	defer sendJson(writer, &resp)

	fileSize, err := strconv.Atoi(req.Header.Get("FileSize"))
	if err != nil {
		log.Println("could not parse file size")
		resp.Message = "Could not parse file size"
		return
	}

	saveFileReq.FileName = sanitizeFileName(req.Header.Get("FileName"))

	saveFileReq.Data, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("could not read request body")
		resp.Message = "Could not read request body."
		return
	}

	var success, alreadyExists bool
	var code string

	for retryCount := 0; !success; retryCount++ {
		// Grab a code until an unused one is found
		for retryCount2 := 0; !success; retryCount2++ {
			success = true
			code = getUuid()

			exists, err := db.FileInfoExists(code)
			if err != nil {
				log.Println("error saving file ", err)
				resp.Message = "Could not save file."
				return
			}

			if exists {
				success = false
			}

			// Stop after 10 retries
			if retryCount2 >= 10 {
				log.Println("could not generate new code")
				resp.Message = "Could not save file."
				return
			}
		}

		// Save the file. Double check if there is a conflict.
		alreadyExists, _, err = file.SaveFile(code, saveFileReq.FileName, saveFileReq.Data)
		if alreadyExists {
			log.Println("code already in use: ", code)
			success = false
		} else if err != nil {
			log.Println("error saving file ", err)
			resp.Message = "Could not save file."
			return
		}

		// Stop after 10 retries
		if retryCount >= 10 {
			log.Println("could not save file")
			resp.Message = "Could not save file."
			return
		}
	}

	err = db.InsertNewFileInfo(code, saveFileReq.FileName, fileSize)
	if err != nil {
		log.Println("error inserting file info into database ", err)
		resp.Message = "Could not save file."
		return
	}

	resp.Success = true
	resp.DownloadLink = makeUrl(downloadFile, "/", code, "/", saveFileReq.FileName)
}

// Appends pathComponents to the end of "https://"+domain and percent encodes them
func makeUrl(pathComponents ...string) string {
	u, _ := url.Parse(protocol + domain)
	for _, component := range pathComponents {
		u.Path += component
	}
	return u.String()
}

// Does not return an error. Instead simply logs and returns if data cannot be marshalled into json.
func sendJson(writer http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("error marshalling json ", err)
		return
	}
	writer.Write(b)
}

func getUuid() string {
	return uuid.Must(uuid.NewV4()).String()
}

func sanitizeFileName(fileName string) string {
	chars := make([]rune, len(fileName), len(fileName))
	for i, c := range fileName {
		if unicode.IsLetter(c) || unicode.IsNumber(c) || c == '.' {
			chars[i] = c
		} else {
			chars[i] = '_'
		}
	}
	return string(chars)
}

func GetEndpoints() []*serverutil.Endpoint {
	return []*serverutil.Endpoint{
		{
			Path:       "/",
			HandleFunc: saveFile,
		},
	}
}
