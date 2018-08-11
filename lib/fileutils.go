package lib

import (
	"os"
		"bufio"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
)

type FileUtils interface {
	CreateFileIfNotExists(filepath string, fileName string) (fileSize int64, err error)
	AppendContent(filepath string, content []byte) (err error)
	ConvertHTTPResponseToBytes(response *http.Response)(bytes []byte,err error)
	GetFileNameFromURL(url string)(fileName string)
}

type File struct{}

func (f *File) CreateFileIfNotExists(filePath string, fileName string) (fileSize int64, err error) {
	os.MkdirAll(filePath, os.ModePerm)
	fileLocation := fmt.Sprintf("%s/%s", filePath, fileName)
	file, err := os.Stat(fileLocation)
	if os.IsNotExist(err) {
		newFile, err := os.Create(filePath)
		if err != nil {
			return 0, err
		}
		defer newFile.Close()
	}

	return file.Size(), nil
}

func (f *File) AppendContent(filepath string, content []byte) (err error) {
	fileHandle, err := os.OpenFile(filepath, os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	writer.Write(content)
	writer.Flush()
	return nil
}

func (f *File) ConvertHTTPResponseToBytes(response *http.Response)(bytes []byte,err error) {
	return ioutil.ReadAll(response.Body)
}

func (f *File) GetFileNameFromURL(url string)(fileName string) {
		tokens := strings.Split(url, "/")
		return tokens[len(tokens)-1]
}
