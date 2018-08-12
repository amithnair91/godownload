package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FileUtils interface {
	CreateFileIfNotExists(filepath string, fileName string) (fileSize int64, err error)
	AppendContent(filepath string, content []byte) (err error)
	ConvertHTTPResponseToBytes(response *http.Response) (bytes []byte, err error)
	GetFileNameFromURL(url string) (fileName string, err error)
}

type File struct{}

func (f *File) CreateFileIfNotExists(filePath string, fileName string) (fileSize int64, err error) {
	os.MkdirAll(filePath, os.ModePerm)
	fileLocation := fmt.Sprintf("%s/%s", filePath, fileName)
	file, err := os.Stat(fileLocation)
	if os.IsNotExist(err) {
		newFile, err := os.Create(fileLocation)
		defer newFile.Close()
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	return file.Size(), nil
}

func (f *File) AppendContent(filepath string, content []byte) (err error) {
	//could use O_CREATE to create file
	fileHandle, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	return binary.Write(fileHandle, binary.LittleEndian, content)
}

func (f *File) ConvertHTTPResponseToBytes(response *http.Response) (bytes []byte, err error) {
	return ioutil.ReadAll(response.Body)
}

func (f *File) GetFileNameFromURL(url string) (fileName string, err error) {
	if len(url) < 1 {
		return "", errors.New("URL cannot be empty")
	}
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1], nil
}
