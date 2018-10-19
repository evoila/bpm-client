package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func ZipPackage(packageName string) {

	specFile := readSpec("./packages/" + packageName)

	filesToZip, err := scanPackageFolder(packageName)
	if err != nil {
		panic(err)
	}

	for _, file := range specFile.Files {

		filesToZip = append(filesToZip, "blobs/"+file)
	}

	zipMe(filesToZip, "./"+packageName+".zip")
}

func scanPackageFolder(packageName string) ([]string, error) {

	return listFiles("./packages/" + packageName)
}

func readSpec(specLocation string) SpecFile {

	yamlFile, err := ioutil.ReadFile(specLocation + "/spec")
	if err != nil {
		panic(err)
	}

	var specFile SpecFile

	err = yaml.Unmarshal(yamlFile, &specFile)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return specFile
}

func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

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

type SpecFile struct {
	Name                string
	Files, Dependencies []string
}

func zipMe(filepaths []string, target string) error {

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(target, flags, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range filepaths {
		if err := addFileToZip(filename, zipw); err != nil {
			return fmt.Errorf("Failed to add file %s to zip: %s", filename, err)
		}
	}
	return nil

}

func addFileToZip(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)

	if err != nil {
		return fmt.Errorf("Error opening file %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {

		return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Error writing %s to zip: %s", filename, err)
	}

	return nil
}
