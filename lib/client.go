package lib

import (
	"net/http"
	"fmt"
)

type Client interface {
	NewHttpClient(url string, existingFileSize string)
	Get(url string, existingFileSize string) (resp *http.Response, err error)
}

type HTTPClient struct {
	client *http.Client
}

func (c *HTTPClient) NewHttpClient(url string, existingFileSize string) () {
	c.client = &http.Client{}
}

func (c *HTTPClient) Get(url string, existingFileSize string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	addRangeHeader(req,existingFileSize)
	return c.client.Do(req)
}

func addRangeHeader(req *http.Request, rangeFrom string) (*http.Request) {
	req.Header.Set("Range", fmt.Sprintf("bytes=%s-", rangeFrom))
	return req
}

