package mocks

import (
	"github.com/stretchr/testify/mock"
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
