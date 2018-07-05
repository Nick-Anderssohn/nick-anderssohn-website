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
		os.Mkdir(fileFolderPath + code, 0755)
		ioutil.WriteFile(filePath, data, 0644)
		err = nil
	} else if err == nil {
		return true, filePath, err
	}
	return
}

func DeleteFile(code string) {
	filePath := fmt.Sprintf("%s%s", fileFolderPath, code)
	if err := os.RemoveAll(filePath); err != nil {
		panic("Could not remove file for code " + code)
	}
}

// returns the path to the code folder which should contain a file if it exists
func GetPath(code string) string {
	return fileFolderPath + code
}