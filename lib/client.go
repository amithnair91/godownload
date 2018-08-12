package lib

import (
	"fmt"
	"net/http"
)

type Client interface {
	NewHttpClient()
	Get(url string, existingFileSize int64) (resp *http.Response, err error)
}

type HTTPClient struct {
	client *http.Client
}

func (c *HTTPClient) NewHttpClient() {
	c.client = &http.Client{}
}

func (c *HTTPClient) Get(url string, existingFileSize int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	addRangeHeader(req, existingFileSize)
	return c.client.Do(req)
}

func addRangeHeader(req *http.Request, rangeFrom int64) *http.Request {
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", rangeFrom))
	return req
}
