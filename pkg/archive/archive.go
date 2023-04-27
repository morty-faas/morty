package archive

import (
	"archive/zip"
	"os"
	"path"
)

// Unzip a source archive into the destination folder
func Unzip(src, dst string) error {
	read, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}

		open, err := file.Open()
		if err != nil {
			return err
		}

		name := path.Join(dst, file.Name)

		if err := os.MkdirAll(path.Dir(name), 0755); err != nil {
			return err
		}

		create, err := os.Create(name)
		if err != nil {
			return err
		}

		if _, err := create.ReadFrom(open); err != nil {
			return err
		}
		create.Close()
	}
	return nil
}
