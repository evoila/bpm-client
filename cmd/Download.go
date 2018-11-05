package cmd

import (
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
	"github.com/pkg/errors"
	"log"
	"os"
)

func Download(url, depth string, requestBody PackageRequestBody) {

	log.Println(depth + "├─  Downloading package: " + requestBody.Name)

	var permission = rest.GetDownloadPermission(url, requestBody)

	if permission == nil {
		log.Println(depth + "└─  Package has not been found")
		return
	}

	s3.DownloadFile(requestBody.Name, *permission)

	var destination, err = os.Getwd()

	if err != nil {
		panic(err)
	}

	log.Println(depth + "├─  Setting it up")
	UnzipPackage(requestBody.Name+".bpm", destination)

	err = os.Remove(requestBody.Name + ".bpm")
	if err != nil {
		panic(err)
	}

	var specFile = ReadSpec(BuildPath([]string{"packages", requestBody.Name}))

	for _, dependency := range specFile.Dependencies {

		_, err := os.Stat(BuildPath([]string{"packages", dependency}))

		if os.IsNotExist(err) {
			meta, err := findDependency(url, dependency)

			if err != nil {
				panic(err)
			}

			dependencyRequest := PackageRequestBody{
				Name:    dependency,
				Version: meta.Version,
				Vendor:  meta.Vendor}
			log.Println(depth + "├─  Handling dependency")

			Download(url, "|	", dependencyRequest)
		}
	}

	log.Println(depth + "└─ Finished package: " + requestBody.Name)
}

func findDependency(url string, dependency string) (*MetaData, error) {
	var possibilities = rest.GetMetaDataForPackageName(url, dependency)

	for _, metaData := range possibilities {
		if metaData.Name == dependency {

			return &metaData, nil
		}
	}

	return nil, errors.New("did not find a matching package")
}
