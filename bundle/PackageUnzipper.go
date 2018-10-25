package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnzipPackage(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.Mkdir(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
		} else {

			dir := filepath.Dir(path)
			os.MkdirAll(dir, 0755)

			fmt.Println(f.Name)

			file, err := os.Create(path)

			if err != nil {
				panic(err)
			}

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(file, rc)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
