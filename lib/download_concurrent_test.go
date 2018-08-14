package lib_test

import (
	"testing"
	"errors"

	"go-downloader/lib"

	"github.com/stretchr/testify/assert"
	"go-downloader/mocks"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
)

const concurrency = 2

func setupConcurrent() (int64, string, string, string, string, *mocks.MockClient, *mocks.MockFileUtils, http.Response) {
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

func TestDownloadFileConcurrentFailsWhenURLIsEmpty(t *testing.T) {
	_, _, filepath, _, _, mockHttpClient, mockFileUtils, _ := setupConcurrent()
	expectedError := "url cannot be empty"
	url := ""
	fileName := ""

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, _ := setupConcurrent()
	expectedError := "file activity failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnClientHeadRequestFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, _ := setupConcurrent()
	expectedError := "client head request failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnClientGetRequestFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, httpResponse := setupConcurrent()
	expectedError := "client get request failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-6").Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnWriteToFileError(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePath, mockHttpClient, mockFileUtils, httpResponse := setup()
	expectedError := "unable to write to file"
	absoluteFilePathPart := fmt.Sprintf("%s-%d",absoluteFilePath,0)
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-6").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}
