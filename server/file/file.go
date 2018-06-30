package file

import (
	"io/ioutil"
	"os"
	"fmt"
)

const (
	fileFolderPath = "download/"
)

func SaveFile(code, fileName string, data []byte) (fileAlreadyExists bool, filePath string, err error) {
	filePath = fmt.Sprintf("%s%s/%s", fileFolderPath, code, fileName)
	// Write the file if it does not exist
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(fileFolderPath + code, 0666)
		ioutil.WriteFile(filePath, data, 0666)
		err = nil
	} else if err == nil {
		return true, filePath, err
	}
	return
}