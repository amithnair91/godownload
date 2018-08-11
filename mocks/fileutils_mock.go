package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockFileUtils struct {
	mock.Mock
}

func (m *MockFileUtils) CreateFileIfNotExists(filepath string) (filesize int64, err error) {
	args := m.Called(filepath)
	if args.Get(0) != nil {
		filesize = args.Get(0).(int64)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockFileUtils) AppendContent(filepath string, content []byte) (err error) {
	args := m.Called(filepath,content)
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
