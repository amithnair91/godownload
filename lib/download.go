package lib

import "fmt"

type Download interface {
	DownloadFile(filepath string, url string, client HTTPClient) (error)
}

type Downloader struct {
	Client Client
	FileUtils FileUtils
}

func (d *Downloader) DownloadFile(filePath string, url string) (error) {
	fileName:= d.FileUtils.GetFileNameFromURL(url)

	fileSize, err := d.FileUtils.CreateFileIfNotExists(filePath,fileName)
	if err != nil {
		return err
	}

	response, err := d.Client.Get(url, fileSize)
	if err != nil {
		return err
	}

	bytes, err := d.FileUtils.ConvertHTTPResponseToBytes(response)
	if err != nil {
		return err
	}

	absoluteFilePath := fmt.Sprintf("%s/%s", filePath, fileName)
	err = d.FileUtils.AppendContent(absoluteFilePath, bytes)
	if err != nil {
		return err
	}

	return nil
}

