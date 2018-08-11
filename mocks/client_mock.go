package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) NewHttpClient(url string, existingFileSize string) {
	m.Called(url,existingFileSize)
	return
}

func (m *MockClient) Get(url string, existingFileSize string) (resp *http.Response, err error) {
	args := m.Called(url,existingFileSize)

	if args.Get(0) != nil {
		resp = args.Get(0).(*http.Response)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

