package api

import (
	"nick-anderssohn-website/full-share/server/concurrency"
	"nick-anderssohn-website/full-share/server/file"
	"nick-anderssohn-website/full-share/server/slog"

	"github.com/Nick-Anderssohn/sherlog"
	"github.com/gorilla/websocket"
)

type uploadProcessor struct {
	Processor      *concurrency.Processor
	conn           *websocket.Conn
	fileWriter     *file.FileWriter
	targetSize     int
	bytesProcessed int
}

func newUploadProcessor(conn *websocket.Conn, fileSize, bufSize int, folder, fileName string) (*uploadProcessor, error) {
	fileWriter, err := file.NewFileWriter(fileSize, bufSize, folder, fileName)
	if err != nil {
		return nil, sherlog.PrependMsg(err, "could not create upload processor")
	}

	up := &uploadProcessor{
		conn:       conn,
		fileWriter: fileWriter,
		targetSize: fileSize,
	}
	up.Processor = concurrency.NewProcessor(up.process)
	return up, nil
}

func (up *uploadProcessor) process(val interface{}) {
	fileSlice := val.([]byte)
	bytesWritten, err := up.fileWriter.Write(fileSlice)

	if err != nil {
		slog.Logger.Error(err)
		up.Processor.Stop()
		writeMsgToWs(up.conn, 501, "", "Could not write file.")
		return
	}

	up.bytesProcessed += bytesWritten

	if up.bytesProcessed >= up.targetSize {
		go up.Processor.Stop()
	}
}
