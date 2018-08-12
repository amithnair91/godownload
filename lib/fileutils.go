package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type FileUtils interface {
	CreateFileIfNotExists(filepath string, fileName string) (fileSize int64, err error)
	GetFileNameFromURL(url string) (fileName string, err error)
	WriteToFile(response *http.Response, filePath string) error
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

func (f *File) WriteToFile(response *http.Response, filePath string) error {
	fo, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fo.Close()
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := response.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) GetFileNameFromURL(url string) (fileName string, err error) {
	if len(url) < 1 {
		return "", errors.New("URL cannot be empty")
	}
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1], nil
}
