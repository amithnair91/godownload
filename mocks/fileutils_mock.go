package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockFileUtils struct {
	mock.Mock
}

func (m *MockFileUtils) CreateFileIfNotExists(filePath string, fileName string) (filesize int64, err error) {
	args := m.Called(filePath, fileName)
	if args.Get(0) != nil {
		filesize = args.Get(0).(int64)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockFileUtils) GetFileNameFromURL(url string) (fileName string, err error) {
	args := m.Called(url)
	if args.Get(0) != nil {
		fileName = args.Get(0).(string)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockFileUtils) WriteToFile(response *http.Response, filePath string) (err error) {
	args := m.Called(response, filePath)
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return
}

func (m *MockFileUtils) MergeFiles(filePaths []string, destinationFilePath string, fileName string) (err error) {
	args := m.Called(filePaths, destinationFilePath, fileName)
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return
}

func (m *MockFileUtils) DeleteFile(filePath string) (err error) {
	args := m.Called(filePath)
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return
}

func (m *MockFileUtils) FileExists(path string) (exists bool) {
	args := m.Called(path)
	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}
	return
}
