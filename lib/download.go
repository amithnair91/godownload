package lib

import (
	"fmt"
)

type Download interface {
	DownloadFile(filepath string, url string) error
	DownloadFileConcurrent(filepath string, url string, concurrency int64) error
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

	response, err := d.Client.ResumeGet(url, fileSize)
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

func (d *Downloader) DownloadFileConcurrent(filePath string, url string, concurrency int64) error {
	fileName, err := d.FileUtils.GetFileNameFromURL(url)
	if err != nil {
		return err
	}
	//absoluteFilePath := fmt.Sprintf("%s/%s", filePath, fileName)

	fileSize, err := d.FileUtils.CreateFileIfNotExists(filePath, fileName)
	if err != nil {
		return err
	}

	headResp, err := d.Client.Head(url)
	if err != nil {
		return err
	}

	rangeList := populateRangeList(headResp.ContentLength, concurrency, fileSize)

	println(fmt.Sprintf("%v",rangeList))


	for _, rangeHeader := range rangeList {
		resp, err := d.Client.Get(url, rangeHeader)
		if err != nil {
			return err
		}
		println(rangeHeader, resp)
	}

	return nil
}

func populateRangeList(contentLength int64, concurrency int64, fileSize int64) []string {
	rangeLimit := contentLength / concurrency
	remainder := contentLength % concurrency
	var rangeList []string
	var i int64
	var previousRange = fileSize
	for i = 0; i < concurrency; i++ {
		nextRange := previousRange + rangeLimit

		byteRange := fmt.Sprintf("%d-%d", previousRange, nextRange)
		rangeList = append(rangeList, byteRange)
		previousRange = previousRange + rangeLimit
	}
	if remainder > 0 {
		finalRange := previousRange + remainder
		finalByteRange := fmt.Sprintf("%d-%d", previousRange, finalRange)
		rangeList = append(rangeList, finalByteRange)
	}
	return rangeList
}
