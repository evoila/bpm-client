package cmd

import (
	"fmt"
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
	"os"
)

func Download(depth string, requestBody PackageRequestBody, config *Config, openId *OpenId) {

	stat, _ := os.Stat(BuildPath([]string{"packages", requestBody.Name}))

	if stat != nil {
		fmt.Println(depth + "└─ Package '" + requestBody.Name + "' already set up in this release.")
		return
	}

	fmt.Println(depth + "├─ Downloading package: ")

	metaData := rest.GetMetaData(requestBody.Vendor, requestBody.Name, requestBody.Version, config, openId)
	permission := rest.GetDownloadPermission(config, requestBody, openId)

	if metaData == nil {
		fmt.Println(depth + "└─ Package '" + requestBody.Name + "' does not exist.")
		return
	}

	fmt.Println(metaData.String(depth))

	if permission == nil {
		fmt.Println(depth + "└─ Download permission has not been granted.")
		return
	}

	err := s3.DownloadFile(requestBody.Name, *permission)

	if err != nil {
		panic(err)
	}

	destination, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Println(depth + "├─ Download successful")
	err = UnzipPackage(requestBody.Name+".bpm", destination)

	if err != nil {
		panic(err)
	}

	defer func() {
		err = os.Remove(requestBody.Name + ".bpm")
		if err != nil {
			panic(err)
		}
	}()

	for _, dependency := range metaData.Dependencies {

		dependencyRequest := PackageRequestBody{
			Name:    dependency.Name,
			Version: dependency.Version,
			Vendor:  dependency.Vendor}
		fmt.Println(depth + "├─ Handling dependency")

		Download(depth+"│  ", dependencyRequest, config, openId)
	}

	fmt.Println(depth + "└─ Finished package: " + requestBody.Name)
}
