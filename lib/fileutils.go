package lib

import (
	"os"
		"bufio"
	"io/ioutil"
	"net/http"
)

type FileUtils interface {
	CreateFileIfNotExists(filepath string) (filesize int64, err error)
	AppendContent(filepath string, content []byte) (err error)
	ConvertHTTPResponseToBytes(response *http.Response)(bytes []byte,err error)
}

type File struct{}

func (f *File) CreateFileIfNotExists(filepath string) (filesize int64, err error) {
	file, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		newFile, err := os.Create(filepath)
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
