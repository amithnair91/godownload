package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"bufio"
	"github.com/schollz/progressbar"
)

type FileUtils interface {
	CreateFileIfNotExists(filepath string, fileName string) (fileSize int64, err error)
	GetFileNameFromURL(url string) (fileName string, err error)
	WriteToFile(response *http.Response, filePath string) error
	MergeFiles(filePaths []string, destinationFilePath string, fileName string) error
	DeleteFile(filePath string) error
	FileExists(path string) bool
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

func (f *File) FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *File) WriteToFile(response *http.Response, filePath string) error {
	fo, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fo.Close()

	chunkSize := 1024
	buf := make([]byte, chunkSize)
	bar := progressbar.New(int(response.ContentLength))

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
		bar.Add(chunkSize)
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

func (f *File) MergeFiles(filePaths []string, destinationFilePath string, fileName string) error {
	_, err := f.CreateFileIfNotExists(destinationFilePath, fileName)
	if err != nil {
		return err
	}

	fileLocation := fmt.Sprintf("%s/%s", destinationFilePath, fileName)
	fo, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fo.Close()

	var count int

	chunkSize := 1024
	for _, filePath := range filePaths {
		data, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer data.Close()

		reader := bufio.NewReader(data)
		part := make([]byte, chunkSize)

		for {
			if count, err = reader.Read(part); err != nil {
				break
			}
			fo.Write(part[:count])
		}
		if err != io.EOF {
			return errors.New(fmt.Sprintf("%s %s: %s", "Error Reading", filePath, err.Error()))
		}
	}
	return nil
}

func (f *File) DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	return err
}
