package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	fileFolderPath = "download/"
)

// SaveFile will save a file if it does not already exist. It will not overwrite existing files.
// Returns alreadyExists, filePath, and an error if there is one.
func SaveFile(code, fileName string, data []byte) (bool, string, error) {
	filePath := fmt.Sprintf("%s%s/%s", fileFolderPath, code, fileName)
	// Write the file if it does not exist
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		os.Mkdir(fileFolderPath+code, 0755)
		err = ioutil.WriteFile(filePath, data, 0644)
		if err != nil {
			return false, "", err
		}
	} else if err == nil {
		return true, filePath, nil
	}

	return false, "", err
}

// DeleteFile deletes the file that corresponds to code
func DeleteFile(code string) error {
	filePath := fmt.Sprintf("%s%s", fileFolderPath, code)
	err := os.RemoveAll(filePath)
	if err != nil {
		return fmt.Errorf("could not remove file for code %s: %s", code, err.Error())
	}
	return nil
}

// GetPath returns the path to the code folder which should contain a file if it exists
func GetPath(code string) string {
	return fileFolderPath + code
}
