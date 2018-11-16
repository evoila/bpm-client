package cmd

import (
	"bufio"
	"fmt"
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/collections"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	. "github.com/evoila/BPM-Client/rest"
	. "github.com/evoila/BPM-Client/s3"
	"os"
	"strings"
)

var set = NewStringSet()

func RunUpdateIfPackagePresentUploadIfNot(packageName string, config *Config) {

	specFile, errMessage := ReadSpec(packageName)

	if errMessage != nil {
		fmt.Println(errMessage)
		return
	}

	metaData := GetMetaData(specFile.Vendor, packageName, specFile.Version, config)

	if metaData == nil {
		fmt.Println("Specified package not found. Uploading instead of updating.")

		CheckIfAlreadyPresentAndUpload(packageName, config)
	} else {
		upload(packageName, specFile.Vendor, "", true, config)
	}
}

func CheckIfAlreadyPresentAndUpload(packageName string, config *Config) {

	specFile, errMessage := ReadSpec(packageName)

	if errMessage != nil {
		fmt.Println(errMessage)
		return
	}

	metaData := GetMetaDataListForPackageName(packageName, config)

	if len(metaData) < 1 || askOperatorForProcedure(metaData) {
		upload(packageName, specFile.Vendor, "", false, config)
	}
}

func upload(packageName, vendor, depth string, update bool, config *Config) {

	if set.Get(packageName) {
		fmt.Println(depth + "└─  Dependency " + packageName + " already handled")
		return
	}

	set.Add(packageName)

	fmt.Println(depth + "├─ Packing: " + packageName)

	specFile, errMessage := ReadSpec(packageName)

	if errMessage != nil {
		fmt.Println(depth + "└─  " + *errMessage)
		return
	}

	pack := "./" + packageName + ".bpm"

	dependencies, errMessage := readDependencies(*specFile, vendor)

	if errMessage != nil {
		fmt.Println(depth + "└─  Error in Dependency: " + *errMessage)
	}

	result := MetaData{
		Name:         packageName,
		Version:      specFile.Version,
		Vendor:       vendor,
		FilePath:     pack,
		Files:        specFile.Files,
		Dependencies: dependencies}

	var permission, oldMeta = RequestPermission(result, false, config)

	if update && oldMeta != nil && askUser(*oldMeta, depth) {
		permission, _ = RequestPermission(result, true, config)
	}

	if permission != nil {

		filesToZip, err := ScanPackageFolder(packageName)
		if err != nil {
			panic(err)
		}

		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./blobs/"))
		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./src/"))

		err = ZipMe(filesToZip, pack)

		defer func() {
			err = os.Remove(result.FilePath)
			if err != nil {
				panic(err)
			}
		}()

		if err != nil {
			panic(err)
		}

		fmt.Println(depth + "├─ Upload Package")

		err = UploadFile(result.FilePath, *permission)
		if err != nil {
			panic(err)
		}

		PutMetaData(config.Url, permission.S3location)
		fmt.Println(depth + "├─ Successfully uploaded")

		for _, dependency := range result.Dependencies {
			fmt.Println(depth + "├─ Handling dependency")

			upload(dependency.Name, dependency.Vendor, "│  "+depth, update, config)
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

func askUser(data MetaData, depth string) bool {

	fmt.Println(depth + "├─ Update Package")
	fmt.Println(data.String(depth))

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Println(depth + "│  This will alter the package with all dependencies. Are you sure? ")

	for !acceptInput(text, depth) {
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func askOperatorForProcedure(data []MetaData) bool {

	fmt.Println("│  Found fhe following packages with similar or same content")

	for _, d := range data {
		fmt.Println(d.String(""))
	}

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	fmt.Println("│  Do you want to upload your version anyway? ")
	for !acceptInput(text, "") {
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func acceptInput(text, depth string) bool {

	var acceptedAnswers = []string{"yes", "y", "no", "n"}

	for _, s := range acceptedAnswers {

		if strings.ToLower(text) == s {
			return true
		}
	}
	fmt.Print(depth + "│  Please enter yes / y or no / n	Answer: ")

	return false
}

func readDependencies(specFile SpecFile, vendor string) ([]Dependency, *string) {

	var dependencies []Dependency

	for _, d := range specFile.Dependencies {
		dependencySpec, errMessage := ReadSpec(d)

		if errMessage != nil {
			return nil, errMessage
		}

		dependencies = append(dependencies, Dependency{
			Name:    dependencySpec.Name,
			Version: dependencySpec.Version,
			Vendor:  vendor}, )
	}

	return dependencies, nil
}
