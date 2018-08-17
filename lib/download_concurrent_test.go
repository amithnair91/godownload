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

const concurrency = 1

func setupConcurrent() (int64, string, string, string, string, *mocks.MockClient, *mocks.MockFileUtils, http.Response, string) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	fileNamePart := fmt.Sprintf("%d-%s", 0, fileName)
	absoluteFilePathPart := fmt.Sprintf("%s/%d-%s", filepath, 0, fileName)
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	content := bytes.NewBufferString("File Contents")
	httpResponse := http.Response{Body: ioutil.NopCloser(content), ContentLength: int64(content.Len())}
	return fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart
}

func TestDownloadFileConcurrentFailsWhenURLIsEmpty(t *testing.T) {
	_, _, filepath, _, _, mockHttpClient, mockFileUtils, _, _ := setupConcurrent()
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
	fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	createFileError := errors.New("file activity failure")
	expectedError := fmt.Errorf("unable to download filepart %v", createFileError)

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, createFileError)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnClientHeadRequestFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _, mockHttpClient, mockFileUtils, _, _ := setupConcurrent()
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
	fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	clientError := errors.New("client get request failure")
	expectedError := fmt.Errorf("unable to download filepart %v", clientError)

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(nil, clientError)
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnWriteToFileError(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	writeToFileError := errors.New("unable to write to file")
	expectedError := fmt.Errorf("unable to download filepart %v", writeToFileError)

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(writeToFileError)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnMergeFileError(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	expectedError := "unable to merge files"
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(nil)
	mockFileUtils.On("MergeFiles", []string{absoluteFilePathPart}, filepath, fileName).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, 1)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentDoesNotFailOnDeleteFileError(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(nil)
	mockFileUtils.On("MergeFiles", []string{absoluteFilePathPart}, filepath, fileName).Return(nil)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, 1)
	assert.NoError(t, err)
	mockFileUtils.Mock.AssertExpectations(t)
}
