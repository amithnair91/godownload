package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockFileUtils struct {
	mock.Mock
}

func (m *MockFileUtils) CreateFileIfNotExists(filePath string, fileName string) (filesize int64, err error) {
	args := m.Called(filePath,fileName)
	if args.Get(0) != nil {
		filesize = args.Get(0).(int64)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockFileUtils) AppendContent(filePath string, content []byte) (err error) {
	args := m.Called(filePath,content)
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return
}

func (m *MockFileUtils)ConvertHTTPResponseToBytes(response *http.Response)(bytes []byte,err error){
	args := m.Called(response)
	if args.Get(0) != nil {
		bytes = args.Get(0).([]byte)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockFileUtils) GetFileNameFromURL(url string)(fileName string, err error) {
	args := m.Called(url)
	if args.Get(0) != nil {
		fileName = args.Get(0).(string)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}