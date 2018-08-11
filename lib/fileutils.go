package lib

import "os"

type FileUtils interface{
	CreateFileIfNotExists(filepath string) (filesize int64, err error)
}

type File struct{}

func (f *File)CreateFileIfNotExists(filepath string) (filesize int64, err error) {
	file, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		newFile, err := os.Create(filepath)
		if err != nil {
			return 0, err
		}
		defer newFile.Close()
	}

	return file.Size(), nil
}
