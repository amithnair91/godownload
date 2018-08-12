package lib

import (
	"fmt"
)

type Download interface {
	DownloadFile(filepath string, url string, client HTTPClient) error
}

type Downloader struct {
	Client    Client
	FileUtils FileUtils
}

func (d *Downloader) DownloadFile(filePath string, url string) error {
	fileName, err := d.FileUtils.GetFileNameFromURL(url)
	if err != nil {
		return err
	}
	absoluteFilePath := fmt.Sprintf("%s/%s", filePath, fileName)

	fileSize, err := d.FileUtils.CreateFileIfNotExists(filePath, fileName)
	if err != nil {
		return err
	}

	response, err := d.Client.Get(url, fileSize)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = d.FileUtils.WriteToFile(response, absoluteFilePath)
	if err != nil {
		return err
	}

	return nil
}
