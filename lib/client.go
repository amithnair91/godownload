package lib

import (
	"fmt"
	"net/http"
)

type Client interface {
	NewHttpClient()
	ResumeGet(url string, existingFileSize int64) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	Get(url string, rangeHeader string) (resp *http.Response, err error)
}

type HTTPClient struct {
	client *http.Client
}

func (c *HTTPClient) NewHttpClient() {
	c.client = &http.Client{}
}

func (c *HTTPClient) ResumeGet(url string, existingFileSize int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	addResumeRangeHeader(req, existingFileSize)
	return c.client.Do(req)
}

func (c *HTTPClient) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	return c.client.Do(req)
}

func (c *HTTPClient) Get(url string, rangeHeader string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	addRangeHeaders(req, rangeHeader)
	return c.client.Do(req)
}

func addRangeHeaders(req *http.Request, rangeHeader string) {
	req.Header.Set("Range", fmt.Sprintf("bytes=%s", rangeHeader))
}

func addResumeRangeHeader(req *http.Request, rangeFrom int64) {
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", rangeFrom))
}
