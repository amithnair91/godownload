package lib_test

import (
	"testing"
	"go-downloader/mocks"
	"go-downloader/lib"
	"github.com/stretchr/testify/assert"
	"errors"
	"net/http"
	"io/ioutil"
	"bytes"
)

func TestDownloadFileFailsOnClientFailure(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"
	expectedError := "client failure"

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)

	assert.Error(t, err)
	assert.EqualError(t, err, "client failure")
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"
	expectedError := "file activity failure"
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}


func TestDownloadFileFailsWhenUnableToConvertResponseToBytes(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content")),}

	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	expectedError := "unable to unmarshall request"
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return([]byte{}, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}


func TestDownloadFileFailsWhenUnableToAppendToFile(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com"
	filepath := "filepath"
	expectedBytes:= []byte("File Content")

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content")),}

	mockFileUtils.On("CreateFileIfNotExists", filepath).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	expectedError := "unable to append to file"
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return(expectedBytes, nil)
	mockFileUtils.On("AppendContent", filepath,expectedBytes).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}
