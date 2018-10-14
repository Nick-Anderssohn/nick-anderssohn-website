package file

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Nick-Anderssohn/sherlog"
)

// FileWriter is meant to write files in parts instead of all at once.
type FileWriter struct {
	TargetFileSize int
	currentSize    int
	f              *os.File
	writer         *bufio.Writer
	filePath       string
}

/*
NewFileWriter returns a new file writer. If the file already exists, or it fails to make a new one,
an error will be returned.
*/
func NewFileWriter(targetFileSize, writerBufSize int, folder, fileName string) (*FileWriter, error) {
	if folder[len(folder)-1] != '/' {
		folder += "/"
	}
	filePath := fmt.Sprintf("%s%s", folder, fileName)
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return nil, sherlog.AsError(fmt.Sprintf("%s already exists", filePath))
	}
	os.MkdirAll(folder, 0755)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		TargetFileSize: targetFileSize,
		f:              f,
		writer:         bufio.NewWriterSize(f, writerBufSize),
		filePath:       filePath,
	}, nil
}

/*
Writes to the buffered writer.
*/
func (fw *FileWriter) Write(bytes []byte) (int, error) {
	bytesWritten, err := fw.writer.Write(bytes)
	fw.currentSize += bytesWritten
	return bytesWritten, sherlog.AsError(err)
}

/*
Complete flushes the buffer and closes the writer. If the written file's size does not equal the target
file size, then the corrupt file will be deleted.
*/
func (fw *FileWriter) Complete() {
	fw.writer.Flush()
	fw.f.Close()

	if fw.currentSize != fw.TargetFileSize {
		os.Remove(fw.filePath)
	}
}
