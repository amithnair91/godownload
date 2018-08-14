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

func setupConcurrent() (int64, string, string, string, string,string, *mocks.MockClient, *mocks.MockFileUtils, http.Response, string) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	fileNamePart := fmt.Sprintf("%s-%d",fileName,0)
	absoluteFilePath := fmt.Sprintf("%s/%s", filepath, fileName)
	absoluteFilePathPart := fmt.Sprintf("%s-%d",absoluteFilePath,0)
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	content := bytes.NewBufferString("File Contents")
	httpResponse := http.Response{Body: ioutil.NopCloser(content), ContentLength: int64(content.Len())}
	return fileSize, url, filepath, fileName, absoluteFilePath,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart
}

func TestDownloadFileConcurrentFailsWhenURLIsEmpty(t *testing.T) {
	_, _, filepath, _, _,_, mockHttpClient, mockFileUtils, _, _ := setupConcurrent()
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
	fileSize, url, filepath, fileName, _,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, _ := setupConcurrent()
	expectedError := "file activity failure"
	fileNamePart := fmt.Sprintf("%s-%d",fileName,0)

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnClientHeadRequestFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _,_, mockHttpClient, mockFileUtils, _, _ := setupConcurrent()
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
	fileSize, url, filepath, fileName, _,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	expectedError := "client get request failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-6").Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnWriteToFileError(t *testing.T) {
	fileSize, url, filepath, fileName, _,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	expectedError := "unable to write to file"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-6").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, concurrency)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentFailsOnMergeFileError(t *testing.T) {
	fileSize, url, filepath, fileName, _,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	expectedError := "unable to merge files"
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(nil)
	mockFileUtils.On("MergeFiles", []string{absoluteFilePathPart}, filepath,fileName).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, 1)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileConcurrentDoesNotFailOnDeleteFileError(t *testing.T) {
	fileSize, url, filepath, fileName, _,absoluteFilePathPart, mockHttpClient, mockFileUtils, httpResponse, fileNamePart := setupConcurrent()
	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("DeleteFile", absoluteFilePathPart).Return(errors.New("could not delete file as it does not exist"))
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileNamePart).Return(fileSize, nil)
	mockHttpClient.On("Head", url).Return(&httpResponse, nil)
	mockHttpClient.On("Get", url, "0-13").Return(&httpResponse, nil)
	mockFileUtils.On("WriteToFile", &httpResponse, absoluteFilePathPart).Return(nil)
	mockFileUtils.On("MergeFiles", []string{absoluteFilePathPart}, filepath,fileName).Return(nil)

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFileConcurrent(filepath, url, 1)
	assert.NoError(t, err)
	mockFileUtils.Mock.AssertExpectations(t)
}