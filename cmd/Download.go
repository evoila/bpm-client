package cmd

import (
	"fmt"
	. "github.com/Nerzal/gocloak"
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
	"log"
	"os"
	"strconv"
)

func DownloadPackageWithDependencies(depth string, requestBody PackagesReference, config *Config, jwt *JWT) {

	stat, _ := os.Stat(BuildPath([]string{"packages", requestBody.Name}))

	if stat != nil {
		fmt.Println(depth + "└─ Package '" + requestBody.Name + "' already set up in this release.")
		return
	}

	metaData := rest.GetMetaData(requestBody, config, jwt)
	permission := rest.GetDownloadPermission(config, requestBody, jwt)

	if metaData == nil {
		fmt.Println(depth + "└─ Package '" + requestBody.Name + "' does not exist.")
		return
	}

	fmt.Println(metaData.String(depth))

	if permission == nil {
		fmt.Println(depth + "└─ DownloadPackageWithDependencies permission has not been granted.")
		return
	}

	fmt.Println(depth + "├─ Downloading " + strconv.FormatInt(metaData.Size/1000000, 10) + "MB")

	err := s3.DownloadFile(requestBody.Name, depth+"├─"+
		"  ", *permission)

	if err != nil {
		log.Fatal(err)
	}

	destination, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	err = UnzipPackage(requestBody.Name+".bpm", destination)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = os.Remove(requestBody.Name + ".bpm")
		if err != nil {
			log.Fatal(err)
		}
	}()

	for _, dependency := range metaData.Dependencies {
		fmt.Println(depth + "├─ Handling dependency")

		DownloadPackageWithDependencies(depth+"│  ", dependency, config, jwt)
	}

	fmt.Println(depth + "└─ Finished package: " + requestBody.Name)
}

func DownloadBySpec(spec DownloadSpec, config *Config, jwt *JWT) {

	for _, packageReference := range spec.Packages {
		fmt.Println("┌─ Beginning with:")

		DownloadPackageWithDependencies("", packageReference, config, jwt)
	}
}
