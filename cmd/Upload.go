package cmd

import (
	"bufio"
	"fmt"
	"github.com/Nerzal/gocloak"
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/collections"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	. "github.com/evoila/BPM-Client/rest"
	. "github.com/evoila/BPM-Client/s3"
	"log"
	"os"
	"strconv"
	"strings"
)

var set = NewStringSet()

func RunUpdateIfPackagePresentUploadIfNot(packageName string, config *Config, jwt *gocloak.JWT) {
	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != "" {
		fmt.Println(errMessage)
		return
	}

	packageReference := PackagesReference{
		Publisher: specFile.Publisher,
		Name:      packageName,
		Version:   specFile.Version}
	metaData := GetMetaData(packageReference, config, jwt)

	if metaData == nil {
		fmt.Println("Specified package not found. Uploading instead of updating.")

		CheckIfAlreadyPresentAndUpload(packageName, config, jwt)
	} else {
		upload(packageName, specFile.Publisher, "", true, config, jwt)
	}
}

func CheckIfAlreadyPresentAndUpload(packageName string, config *Config, jwt *gocloak.JWT) {

	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != "" {
		fmt.Println(errMessage)
		return
	}

	metaData := GetMetaDataListForPackageName(packageName, config, jwt)

	if len(metaData) < 1 || askOperatorForProcedure(metaData) {
		upload(packageName, specFile.Publisher, "", false, config, jwt)
	}
}

func upload(packageName, publisher, depth string, update bool, config *Config, openId *gocloak.JWT) {
	if set.Get(packageName) {
		fmt.Println(depth + "└─  Dependency " + packageName + " already handled")
		return
	}

	set.Add(packageName)
	fmt.Println(depth + "├─ Packing: " + packageName)
	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != "" {
		fmt.Println(depth + "└─  " + errMessage)
		return
	}

	pack := "./" + packageName + ".bpm"
	dependencies, errMessage := readDependencies(*specFile)

	if errMessage != "" {
		fmt.Println(depth + "└─  Error in Dependency: " + errMessage)
	}

	result := MetaData{
		Name:         packageName,
		Version:      specFile.Version,
		Publisher:    publisher,
		FilePath:     pack,
		Files:        specFile.Files,
		Dependencies: dependencies,
		Description:  specFile.Description,
		Stemcell:     specFile.Stemcell,
		Url:          specFile.Url}

	var permission, oldMeta = RequestPermission(result, false, config, openId)

	if update && oldMeta != nil && AskUser(*oldMeta, depth, "Update Package", "│  This will alter the access to the package with all dependencies. Are you sure? ") {
		permission, _ = RequestPermission(result, true, config, openId)
	}

	if permission != nil {

		filesToZip, err := ScanPackageFolder(packageName)
		if err != nil {
			log.Fatal(err)
		}

		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./blobs/"))
		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./src/"))
		size, err := ZipMe(filesToZip, pack)

		defer func() {
			err = os.Remove(result.FilePath)
			if err != nil {
				log.Fatal(err)
			}
		}()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(depth + "├─ Upload package. Size " + strconv.FormatInt(size/1000000, 10) + "MB")
		err = UploadFile(result.FilePath, depth+"├─", *permission)

		if err != nil {
			log.Fatal(err)
		}

		PutMetaData(config.Url, permission.S3location, openId, size)
		fmt.Println(depth + "├─ Successfully uploaded")

		for _, dependency := range result.Dependencies {
			fmt.Println(depth + "├─ Handling dependency")

			upload(dependency.Name, dependency.Publisher, "│  "+depth, update, config, openId)
		}

		fmt.Println(depth + "└─ Finished packing: " + packageName)

	} else {
		if oldMeta != nil {
			fmt.Println(depth + "└─ Skipped. Already present. Use update if you want to replace it.")
		} else {
			fmt.Println(depth + "└─ Skipped. Not a member of the publisher: " + result.Publisher)
		}
	}
}

func askOperatorForProcedure(data []MetaData) bool {
	fmt.Println("│  Found fhe following packages with similar or same content")

	for _, d := range data {
		fmt.Println(d.String(""))
	}

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	fmt.Println("│  Do you want to upload your version anyway? ")
	for !AcceptInput(text, "") {
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func readDependencies(specFile SpecFile) ([]PackagesReference, string) {
	var dependencies []PackagesReference

	for _, d := range specFile.Dependencies {
		dependencySpec, errMessage := ReadAndValidateSpec(d)

		if errMessage != "" {
			return nil, errMessage
		}

		dependencies = append(dependencies, PackagesReference{
			Name:      dependencySpec.Name,
			Version:   dependencySpec.Version,
			Publisher: dependencySpec.Publisher}, )
	}

	return dependencies, ""
}
