package main

import (
	"fmt"
	"godownload/lib"
)

func main() {
	file := lib.File{}
	client := lib.HTTPClient{}
	client.NewHttpClient()
	downloader := lib.Downloader{FileUtils: &file, Client: &client}
	println("Start Download of File")

	err := downloader.DownloadFileConcurrent("./", "http://dynamodb-local.s3-website-us-west-2.amazonaws.com/dynamodb_local_2016-05-17.zip", 7)

	println(fmt.Sprintf("%v", err))
	println("Finished Downloading File")
}
