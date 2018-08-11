package lib_test

import (
	"testing"
	"go-downloader/mocks"
	"go-downloader/lib"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestDownloadFileFailsOnClientFailure(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"
	expectedError := errors.New("client failure")
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(nil, expectedError)
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)

	assert.Error(t, err)
	assert.EqualError(t, err, "client failure")
	mockFileUtils.Mock.AssertExpectations(t)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"
	expectedError := errors.New("file activity failure")
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, expectedError)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, "file activity failure")
	mockFileUtils.Mock.AssertExpectations(t)
}
