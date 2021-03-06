package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) NewHttpClient() {
	m.Called()
	return
}

func (m *MockClient) ResumeGet(url string, existingFileSize int64) (resp *http.Response, err error) {
	args := m.Called(url, existingFileSize)

	if args.Get(0) != nil {
		resp = args.Get(0).(*http.Response)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockClient) Head(url string) (resp *http.Response, err error) {
	args := m.Called(url)

	if args.Get(0) != nil {
		resp = args.Get(0).(*http.Response)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}

func (m *MockClient) Get(url string, rangeHeader string) (resp *http.Response, err error) {
	args := m.Called(url, rangeHeader)

	if args.Get(0) != nil {
		resp = args.Get(0).(*http.Response)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return
}
