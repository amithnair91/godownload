package lib

import (
	"fmt"
	"os"
	"sort"
	"sync"
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

	headResp, err := d.Client.Head(url)
	if err != nil {
		return err
	}

	rangeList := populateRangeList(headResp.ContentLength, concurrency, 0)

	filePartsChan := make(chan []string, 1)
	downloadErrChan := make(chan error, 1)

	filePartsChan <- []string{}
	downloadErrChan <- nil

	var wg sync.WaitGroup
	wg.Add(len(rangeList))
	for index, rangeHeader := range rangeList {
		go download(&wg, downloadErrChan, filePartsChan, filePath, fileName, index, d, url, rangeHeader)
	}
	wg.Wait()

	err = <-downloadErrChan
	if err != nil {
		return err
	}

	fileParts := <-filePartsChan

	err = d.FileUtils.MergeFiles(fileParts, filePath, fileName)
	if err != nil {
		return err
	}

	sort.Strings(fileParts)
	for _, filePartName := range fileParts {
		err = os.Remove(filePartName)
		if err != nil {
			println("unable to remove filePart: ", filePartName)
		}
	}

	return nil
}

func download(wg *sync.WaitGroup, downloadErr chan error, fileParts chan []string,
	filePath string, fileName string, index int, d *Downloader, url string,
	rangeHeader string) {
	defer wg.Done()
	err := <-downloadErr

	if err != nil {
		downloadErr <- err
		return
	}

	parts := <-fileParts

	absoluteFilePartPath := fmt.Sprintf("%s/%d-%s", filePath, index, fileName)
	//delete if filepart exists
	d.FileUtils.DeleteFile(absoluteFilePartPath)
	filePartName := fmt.Sprintf("%d-%s", index, fileName)
	_, err = d.FileUtils.CreateFileIfNotExists(filePath, filePartName)
	if err != nil {
		downloadErr <- err
		return
	}
	parts = append(parts, absoluteFilePartPath)
	response, err := d.Client.Get(url, rangeHeader)
	if err != nil {
		downloadErr <- err
		return
	}
	err = d.FileUtils.WriteToFile(response, absoluteFilePartPath)
	if err != nil {
		downloadErr <- err
		return
	}
	fileParts <- parts
	downloadErr <- err
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
