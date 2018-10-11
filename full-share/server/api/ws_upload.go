package api

import (
	"encoding/json"
	"nick-anderssohn-website/full-share/server/db"

	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"nick-anderssohn-website/full-share/server/serverutil"

	"os"

	"unicode"

	"nick-anderssohn-website/full-share/server/httputil"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

/*
A successful save file hand shake will look like the following:
Client sends:
    {
        "FileSize": 1000,
        "FileName": "Chicken Dinner"
    }

Server replies:
	{
        "StatusCode": 200,
		"StatusMsg": "Ok",
        "Message": ""
    }

Then in a loop until file size is reached:
	Client sends file bytes.

	The server responds with
	{
        "StatusCode": 200,
		"StatusMsg": "Ok",
        "Message": ""
    }


Then the server will finish up with:
	{
        "StatusCode": 200,
		"StatusMsg": "Ok",
        "Message": "The url here"
    }

At any point in time, the server might respond with status indicating an error code and message containing an
error message.
*/

/*
UploadSetupMsg is the json message sent from the browser with information about the file
it would like to upload
*/
type UploadSetupMsg struct {
	FileSize int    `json:"FileSize"`
	FileName string `json:"FileName"`
}

/*
WsResponse contains an http status code and a message
*/
type WsResponse struct {
	StatusCode int    `json:"StatusCode"`
	StatusMsg  string `json:"StatusMsg"`
	Message    string `json:"Message"`
}

const (
	downloadFile      = "/download"
	fileFolderPath    = "download/"
	wsReadMaxDuration = time.Minute * 60
	maxWriterBufSize  = 1024 * 1024            // 1 MiB
	maxFileSize       = 1024 * 1024 * 1024 * 5 // 5 GiB
)

var upgrader = websocket.Upgrader{}

var domain string
var protocol string

func init() {
	domain = os.Getenv("WEBSITE_DOMAIN")
	protocol = os.Getenv("WEBSITE_PROTOCOL")
}

func recoverAndLog() {
	if e := recover(); e != nil {
		log.Println(e)
	}
}

func saveViaWebsocket(w http.ResponseWriter, r *http.Request) {
	defer recoverAndLog()

	// Upgrade to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade to ws conn ", err)
		return
	}
	defer conn.Close() // Close ws connection on func return

	// Handle setup messages
	fileInfo, code, errWithCode := setupUpload(conn)
	if errWithCode != nil {
		if code != "" {
			deleteCodeFromDb(code)
		}
		writeMsgToWs(conn, errWithCode.Status.Code, errWithCode.Status.Msg, "Could not setup upload.")
		return
	}

	// Format file name as download/code/fileName
	folder := fmt.Sprintf("%s%s/", fileFolderPath, code)

	// Open a file writer
	bufSize := maxWriterBufSize
	if bufSize > fileInfo.FileSize {
		bufSize = fileInfo.FileSize
	}

	// Create the processor that will be used to write the file to disk
	processor, err := newUploadProcessor(conn, fileInfo.FileSize, bufSize, folder, fileInfo.FileName)
	if err != nil {
		writeMsgToWs(conn, httputil.InternalServerError.Code, httputil.InternalServerError.Msg, "Could not write file.")
		return
	}

	defer processor.fileWriter.Complete() // flush and close or handle a messed up write process on func return
	processor.Processor.DoneChan = make(chan bool)

	bytesReceived := 0

	// Defer a check to remove bad database entries
	defer func() {
		if bytesReceived != fileInfo.FileSize {
			deleteCodeFromDb(code)
		}
	}()

	for bytesReceived < fileInfo.FileSize {
		if processor.Processor.Stopped {
			log.Println("failed to write to file ", err)
			writeMsgToWs(conn, httputil.InternalServerError.Code, httputil.InternalServerError.Msg, "Could not write file.")
			return
		}

		// Read message
		conn.SetReadDeadline(time.Now().Add(wsReadMaxDuration))
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message ", err)
			writeMsgToWs(conn, httputil.InternalServerError.Code, httputil.InternalServerError.Msg, "Could not read message.")
			return
		}

		processor.Processor.Push(msgBytes)
		writeMsgToWs(conn, httputil.Ok.Code, httputil.Ok.Msg, "")

		// Record bytes received and write to file writer
		bytesReceived += len(msgBytes)
	}

	<-processor.Processor.DoneChan

	// If we received the correct number of bytes, then it was a success.
	if bytesReceived == fileInfo.FileSize {
		downloadLink := makeUrl(downloadFile, "/", code, "/", fileInfo.FileName)
		writeMsgToWs(conn, httputil.Ok.Code, httputil.Ok.Msg, downloadLink)
	} else {
		writeMsgToWs(conn, httputil.InternalServerError.Code, httputil.InternalServerError.Msg, "Something went wrong.")
	}
}

/*
Reads the setup message from the websocket and returns file information, code, error.
*/
func setupUpload(conn *websocket.Conn) (*UploadSetupMsg, string, *httputil.ErrWithStatus) {
	// Receive the first message
	conn.SetReadDeadline(time.Now().Add(wsReadMaxDuration))
	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		return nil, "", &httputil.ErrWithStatus{Status: httputil.InternalServerError, Err: err}
	}

	// Unmarshal
	var setupMsg UploadSetupMsg
	err = json.Unmarshal(msgBytes, &setupMsg)
	if err != nil {
		return nil, "", &httputil.ErrWithStatus{Status: httputil.BadRequest, Err: err}
	}

	if setupMsg.FileSize > maxFileSize {
		return nil, "", &httputil.ErrWithStatus{
			Status: httputil.BadRequest,
			Err:    fmt.Errorf("file too large. max allowed file size: %d bytes. provided file size: %d bytes", maxFileSize, setupMsg.FileSize),
		}
	}

	setupMsg.FileName = sanitizeFileName(setupMsg.FileName)

	// Get a code for it
	code, err := getUuid()
	if err != nil {
		return nil, "", &httputil.ErrWithStatus{Status: httputil.InternalServerError, Err: err}
	}

	// Insert a record into the database
	err = db.InsertNewFileInfo(code, setupMsg.FileName, setupMsg.FileSize)
	if err != nil {
		return nil, code, &httputil.ErrWithStatus{Status: httputil.InternalServerError, Err: err}
	}

	// Write success and return
	writeMsgToWs(conn, httputil.Ok.Code, httputil.Ok.Msg, "")
	return &setupMsg, code, nil
}

func writeMsgToWs(conn *websocket.Conn, statusCode int, statusMsg, message string) error {
	return conn.WriteJSON(&WsResponse{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
		Message:    message,
	})
}

func getUuid() (code string, err error) {
	// Grab a code until an unused one is found
	for retryCount, success := 0, false; !success; retryCount++ {
		code = uuid.NewV4().String()
		success = true
		exists, err := db.FileInfoExists(code)
		if err != nil {
			success = false
		}

		if exists {
			success = false
		}

		// Stop after 10 retries
		if retryCount >= 10 {
			return "", fmt.Errorf("could not generate new code")
		}
	}
	return
}

func deleteCodeFromDb(code string) {
	err := db.DeleteDbEntryFromCode(code)
	if err != nil {
		log.Println("failed to delete db entry for code ", code)
	}
}

// GetEndpoints returns all the endpoints for the full share server.
func GetEndpoints() []*serverutil.Endpoint {
	return []*serverutil.Endpoint{
		{
			Path:       "/",
			HandleFunc: saveViaWebsocket,
		},
	}
}

// Appends pathComponents to the end of "https://"+domain and percent encodes them
func makeUrl(pathComponents ...string) string {
	u, _ := url.Parse(protocol + domain)
	for _, component := range pathComponents {
		u.Path += component
	}
	return u.String()
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
