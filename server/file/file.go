package file

import (
	"io/ioutil"
	"os"
)

const (
	imageFolderPath = "data/"
	fileExt         = ".zlib"
)

func SaveFile(code string, data []byte) (fileAlreadyExists bool, err error) {
	filePath := imageFolderPath + code + fileExt
	// Write the file if it does not exist
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		ioutil.WriteFile(filePath, data, 0666)
		err = nil
	} else if err == nil {
		return true, err
	}
	return
}
