package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	. "os"
	. "path/filepath"
	"strings"

	"github.com/evoila/BPM-Client/helpers"
)

func ScanPackageFolder(packageName string) ([]string, error) {

	return listFiles("./packages/" + packageName)
}

func ScanFolderAndFilter(files []string, folder string) []string {

	var matches []string

	for _, file := range files {
		curMatches, err := Glob(folder + file)

		if err != nil {
			panic(err)
		}

		for _, match := range curMatches {

			info, err := Stat(match)
			if err != nil {
				panic(err)
			}

			if info.IsDir() {
				content, err := listFiles(match)
				if err != nil {
					panic(err)
				}
				matches = helpers.MergeStringList(matches, content)
			} else {
				normalized := strings.TrimPrefix(match, "./")
				matches = append(matches, normalized)
			}
		}
	}

	return matches
}

func listFiles(root string) ([]string, error) {
	var files []string

	err := Walk(root, func(path string, info FileInfo, err error) error {

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ZipMe(filePath []string, target string) (int64, error) {

	file, err := Create(target)

	if err != nil {
		return 0, fmt.Errorf("failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	for _, filename := range filePath {
		if err := addFileToZip(filename, zipWriter); err != nil {
			return 0, fmt.Errorf("failed to add file %s to zip: %s", filename, err)
		}
	}
	info, err := file.Stat()

	if err != nil {
		return 0, fmt.Errorf("failed to read file info: %s", err)
	}

	return info.Size(), nil
}

func addFileToZip(filename string, zipw *zip.Writer) error {
	file, err := Open(filename)

	if err != nil {
		return fmt.Errorf("error opening file %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {

		return fmt.Errorf("error adding file; '%s' to zip : %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("error writing %s to zip: %s", filename, err)
	}

	return nil
}
