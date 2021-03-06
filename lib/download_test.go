package lib_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/amithnair91/godownload/lib"
	"github.com/amithnair91/godownload/mocks"
	"github.com/stretchr/testify/assert"
)

func setup() (int64, string, string, string, string, *mocks.MockClient, *mocks.MockFileUtils, http.Response) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	absoluteFilePath := fmt.Sprintf("%s/%s", filepath, fileName)
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	content := bytes.NewBufferString("File Contents")
	httpResponse := http.Response{Body: ioutil.NopCloser(content), ContentLength: int64(content.Len())}
	return fileSize, url, filepath, fileName, absoluteFilePath, mockHttpClient, mockFileUtils, httpResponse
}

func TestDownloadFileFailsWhenURLIsEmpty(t *testing.T) {
	_, _, filepath, _, _, mockHttpClient, mockFileUtils, _ := setup()
	expectedError := "url cannot be empty"
	url := ""
	fileName := ""

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsOnClientGetRequestFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, _ := setup()
	expectedError := "client failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("ResumeGet", url, fileSize).Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, _ := setup()
	expectedError := "file activity failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsOnWriteToFileError(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePath, mockHttpClient, mockFileUtils, httpResponse := setup()
	expectedError := "unable to write to file"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("ResumeGet", url, fileSize).Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePath).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}
