package lib_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/amithnair91/godownload/lib"
	"github.com/stretchr/testify/assert"
)

func TestDownloadSuccess(t *testing.T) {
	pwd, err := os.Getwd()

	assert.NoError(t, err)

	fileName := "Sample-Spreadsheet-10000-rows.xls"
	url := "https://www.sample-videos.com/xls/Sample-Spreadsheet-10000-rows.xls"
	filePath := filepath.FromSlash(fmt.Sprintf("%s/%s", pwd, fileName))
	file := lib.File{}
	client := lib.HTTPClient{}
	client.NewHttpClient()
	downloader := lib.Downloader{FileUtils: &file, Client: &client}

	err = downloader.DownloadFile(pwd, url)

	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	err = os.Remove(filePath)

	assert.NoError(t, err)
}

func TestDownloadConcurrentSuccess(t *testing.T) {
	pwd, err := os.Getwd()

	assert.NoError(t, err)

	fileName := "Sample-Spreadsheet-10000-rows.xls"
	url := "https://www.sample-videos.com/xls/Sample-Spreadsheet-10000-rows.xls"
	filePath := filepath.FromSlash(fmt.Sprintf("%s/%s", pwd, fileName))
	file := lib.File{}
	client := lib.HTTPClient{}
	client.NewHttpClient()
	downloader := lib.Downloader{FileUtils: &file, Client: &client}

	var filePartNames []string
	var i int
	for i = 0; i < concurrency; i++ {
		filePartNames = append(filePartNames, fmt.Sprintf("%s-%d", fileName, i))
	}

	err = downloader.DownloadFileConcurrent(pwd, url, 6)

	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	err = os.Remove(filePath)
	assert.NoError(t, err)
}
