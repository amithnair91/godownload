package lib_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-downloader/lib"
	"go-downloader/mocks"
	"io/ioutil"
	"net/http"
	"testing"
)

func setup(httpResponseStr string) (int64, string, string, string, string, []byte, *mocks.MockClient, *mocks.MockFileUtils, http.Response) {
	fileSize := int64(0)
	url := "www.someurl.com/file.txt"
	filepath := "filepath"
	fileName := "file.txt"
	absoluteFilePath := fmt.Sprintf("%s/%s", filepath, fileName)
	expectedBytes := []byte(httpResponseStr)
	mockHttpClient := &mocks.MockClient{}
	mockFileUtils := &mocks.MockFileUtils{}
	httpResponse := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("File Content"))}
	return fileSize, url, filepath, fileName, absoluteFilePath, expectedBytes, mockHttpClient, mockFileUtils, httpResponse
}

func TestDownloadFileFailsWhenURLIsEmpty(t *testing.T) {
	_, _, filepath, _, _, _, mockHttpClient, mockFileUtils, _ := setup("File Content")
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

func TestDownloadFileFailsOnClientFailure(t *testing.T) {
	fileSize, url, filepath, fileName, _, _, mockHttpClient, mockFileUtils, _ := setup("")
	expectedError := "client failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(nil, errors.New(expectedError))
	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockHttpClient.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToCreateFile(t *testing.T) {
	fileSize, url, filepath, fileName, _, _, mockHttpClient, mockFileUtils, _ := setup("")
	expectedError := "file activity failure"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToConvertResponseToBytes(t *testing.T) {
	fileSize, url, filepath, fileName, _, _, mockHttpClient, mockFileUtils, httpResponse := setup("File Content")
	expectedError := "unable to unmarshall request"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return([]byte{}, errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}

func TestDownloadFileFailsWhenUnableToAppendToFile(t *testing.T) {
	fileSize, url, filepath, fileName, absoluteFilePath, expectedBytes, mockHttpClient, mockFileUtils, httpResponse := setup("File Content")
	expectedError := "unable to append to file"

	mockFileUtils.On("GetFileNameFromURL", url).Return(fileName, nil)
	mockFileUtils.On("CreateFileIfNotExists", filepath, fileName).Return(fileSize, nil)
	mockHttpClient.On("Get", url, fileSize).Return(&httpResponse, nil)
	mockFileUtils.On("ConvertHTTPResponseToBytes", &httpResponse).Return(expectedBytes, nil)
	mockFileUtils.On("AppendContent", absoluteFilePath, expectedBytes).Return(errors.New(expectedError))

	downloader := lib.Downloader{Client: mockHttpClient, FileUtils: mockFileUtils}

	err := downloader.DownloadFile(filepath, url)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockFileUtils.Mock.AssertExpectations(t)
}
