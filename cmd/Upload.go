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
	"os"
	"strconv"
	"strings"
)

var set = NewStringSet()

func RunUpdateIfPackagePresentUploadIfNot(packageName string, config *Config, jwt *gocloak.JWT) {

	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != nil {
		fmt.Println(errMessage)
		return
	}

	metaData := GetMetaData(specFile.Vendor, packageName, specFile.Version, config, jwt)

	if metaData == nil {
		fmt.Println("Specified package not found. Uploading instead of updating.")

		CheckIfAlreadyPresentAndUpload(packageName, config, jwt)
	} else {
		upload(packageName, specFile.Vendor, "", true, config, jwt)
	}
}

func CheckIfAlreadyPresentAndUpload(packageName string, config *Config, jwt *gocloak.JWT) {

	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != nil {
		fmt.Println(errMessage)
		return
	}

	metaData := GetMetaDataListForPackageName(packageName, config, jwt)

	if len(metaData) < 1 || askOperatorForProcedure(metaData) {
		upload(packageName, specFile.Vendor, "", false, config, jwt)
	}
}

func upload(packageName, vendor, depth string, update bool, config *Config, openId *gocloak.JWT) {

	if set.Get(packageName) {
		fmt.Println(depth + "└─  PackagesReference " + packageName + " already handled")
		return
	}

	set.Add(packageName)

	fmt.Println(depth + "├─ Packing: " + packageName)

	specFile, errMessage := ReadAndValidateSpec(packageName)

	if errMessage != nil {
		fmt.Println(depth + "└─  " + *errMessage)
		return
	}

	pack := "./" + packageName + ".bpm"

	dependencies, errMessage := readDependencies(*specFile)

	if errMessage != nil {
		fmt.Println(depth + "└─  Error in PackagesReference: " + *errMessage)
	}

	result := MetaData{
		Name:         packageName,
		Version:      specFile.Version,
		Vendor:       vendor,
		FilePath:     pack,
		Files:        specFile.Files,
		Dependencies: dependencies,
		Description:  specFile.Description,
		Stemcell:     specFile.Stemcell,
		Url:          specFile.Url}

	var permission, oldMeta = RequestPermission(result, false, config, openId)

	if update && oldMeta != nil && AskUser(*oldMeta, depth, "│  This will alter the package with all dependencies. Are you sure? ") {
		permission, _ = RequestPermission(result, true, config, openId)
	}

	if permission != nil {

		filesToZip, err := ScanPackageFolder(packageName)
		if err != nil {
			panic(err)
		}

		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./blobs/"))
		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./src/"))
		size, err := ZipMe(filesToZip, pack)

		defer func() {
			err = os.Remove(result.FilePath)
			if err != nil {
				panic(err)
			}
		}()

		if err != nil {
			panic(err)
		}

		fmt.Println(depth + "├─ Upload package. Size " + strconv.FormatInt(size/1000000, 10) + "MB")

		err = UploadFile(result.FilePath, depth+"├─", *permission)
		if err != nil {
			panic(err)
		}

		PutMetaData(config.Url, permission.S3location, openId, size)
		fmt.Println(depth + "├─ Successfully uploaded")

		for _, dependency := range result.Dependencies {
			fmt.Println(depth + "├─ Handling dependency")

			upload(dependency.Name, dependency.Vendor, "│  "+depth, update, config, openId)
		}

		fmt.Println(depth + "└─ Finished packing: " + packageName)

	} else {
		if oldMeta != nil {
			fmt.Println(depth + "└─ Skipped. Already present. Use update if you want to replace it")
		} else {
			fmt.Println(depth + "└─ Skipped. Not a Member of the Vendor: " + result.Vendor)
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

func readDependencies(specFile SpecFile) ([]PackagesReference, *string) {

	var dependencies []PackagesReference

	for _, d := range specFile.Dependencies {
		dependencySpec, errMessage := ReadAndValidateSpec(d)

		if errMessage != nil {
			return nil, errMessage
		}

		dependencies = append(dependencies, PackagesReference{
			Name:    dependencySpec.Name,
			Version: dependencySpec.Version,
			Vendor:  dependencySpec.Vendor}, )
	}

	return dependencies, nil
}
