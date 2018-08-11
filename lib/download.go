package lib

type Download interface {
	DownloadFile(filepath string, url string, client HTTPClient) (error)
}

type Downloader struct {
	Client Client
}

func (d *Downloader) DownloadFile(filePath string, url string) (error) {
	_, err := d.Client.Get(url, "0")
	if err != nil {
		return err
	}

	return nil
}
