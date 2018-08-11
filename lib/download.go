package lib

type Download interface {
	DownloadFile(filepath string, url string, client HTTPClient) (error)
}

type Downloader struct {
	Client Client
	FileUtils FileUtils
}

func (d *Downloader) DownloadFile(filePath string, url string) (error) {
	fileSize, err := d.FileUtils.CreateFileIfNotExists(filePath)
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

	err = d.FileUtils.AppendContent(filePath, bytes)
	if err != nil {
		return err
	}

	return nil
}
