package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) DownloadFile(filePath string, url string) (err error) {
	args := m.Called(filePath, url)
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return
}
