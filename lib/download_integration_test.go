package lib_test

import (
	"testing"
	"go-downloader/lib"
	"os"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"fmt"
)

func TestDownloadSuccess(t *testing.T) {
	pwd, err := os.Getwd()

	assert.NoError(t,err)

	fileName := "Sample-Spreadsheet-10000-rows.xls"
	url := "https://www.sample-videos.com/xls/Sample-Spreadsheet-10000-rows.xls"
	filePath := filepath.FromSlash(fmt.Sprintf("%s/%s",pwd, fileName))
	file := lib.File{}
	client := lib.HTTPClient{}
	client.NewHttpClient()
	downloader := lib.Downloader{FileUtils: &file, Client: &client}

	err = downloader.DownloadFile(pwd, url)

	assert.NoError(t,err)
	assert.FileExists(t,filePath)

	err = os.Remove(filePath)

	assert.NoError(t,err)
}

