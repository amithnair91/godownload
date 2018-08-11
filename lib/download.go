package lib

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

	err = d.FileUtils.AppendContent(filePath, bytes)
	if err != nil {
		return err
	}

	return nil
}


//func downloadFile(filepath string, url string) (err error) {
//
//	// Create the file
//	out, err := os.Create(filepath)
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//
//	// Get the data
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	// Check server response
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("bad status: %s", resp.Status)
//	}
//
//	// Writer the body to file
//	_, err = io.Copy(out, resp.Body)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func downloadFromUrl(url string) {
//	tokens := strings.Split(url, "/")
//	fileName := tokens[len(tokens)-1]
//	fmt.Println("Downloading", url, "to", fileName)
//
//	// TODO: check file existence first with io.IsExist
//	output, err := os.Create(fileName)
//	if err != nil {
//		fmt.Println("Error while creating", fileName, "-", err)
//		return
//	}
//	defer output.Close()
//
//	response, err := http.Get(url)
//	if err != nil {
//		fmt.Println("Error while downloading", url, "-", err)
//		return
//	}
//	defer response.Body.Close()
//
//	n, err := io.Copy(output, response.Body)
//	if err != nil {
//		fmt.Println("Error while downloading", url, "-", err)
//		return
//	}
//
//	fmt.Println(n, "bytes downloaded.")
//}
