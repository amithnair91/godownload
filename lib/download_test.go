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
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	expectedError := "client failure"

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)

	assert.Error(t, err)
	assert.EqualError(t, err, "client failure")
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	expectedError := "file activity failure"
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, errors.New(expectedError))
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToConvertResponseToBytes(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content")),}

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
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
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	expectedBytes := []byte("File Content")

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content")),}

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	expectedError := "unable to append to file"
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return(expectedBytes, nil)
	mockFileUtils.On("AppendContent", filepath, expectedBytes).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToGetFileNameFromURL(t *testing.T) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	expectedBytes := []byte("unable to get filename from url")
	fileName := "file.txt"

	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content")),}

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return(expectedBytes, nil)
	mockFileUtils.On("AppendContent", filepath, expectedBytes).Return(nil)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.NoError(t, err)
	mockFileUtils.Mock.AssertExpectations(t)
}
