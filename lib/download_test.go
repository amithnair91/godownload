package lib_test

import (
	"testing"
	"go-downloader/mocks"
	"go-downloader/lib"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestDownloadFileFailsOnClientFailure(t *testing.T) {
	fileSize := "0"
	url := "www.someurl.com"
	expectedError := errors.New("client failure")
	mockHttpClient := &mocks.MockClient{}
	mockHttpClient.On("Get",url,fileSize).Return(nil, expectedError)
	downloader := lib.Downloader{Client:mockHttpClient,}

	err := downloader.DownloadFile("",url)

	assert.Error(t,err)
	assert.EqualError(t,err,"client failure")
}
