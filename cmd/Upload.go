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
	"log"
	"os"
	"strings"
)

var set = NewStringSet()

func CheckIfAlreadyPresentAndUpload(url, packageName, vendor string) {

	metaData := GetMetaDataForPackageName(url, packageName)

	if len(metaData) < 1 || askOperatorForProcedure(metaData) {
		upload(url, packageName, vendor, "")
	}
}

func upload(url, packageName, vendor, depth string) {

	if set.Get(packageName) {
		log.Println(depth + "└─  Dependency " + packageName + " already handled")
		return
	}

	set.Add(packageName)

	log.Println(depth + "├─ Packing: " + packageName)

	specFile := ReadSpec("./packages/" + packageName)
	pack := "./" + packageName + ".bpm"

	dependencies := readDependencies(specFile, vendor)

	result := MetaData{
		Name:         packageName,
		Version:      specFile.Version,
		Vendor:       vendor,
		FilePath:     pack,
		Files:        specFile.Files,
		Dependencies: dependencies}

	var response, _ = PutMetaData(url, result, false)

	if response != nil {

		filesToZip, err := ScanPackageFolder(packageName)
		if err != nil {
			panic(err)
		}

		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./blobs/"))
		filesToZip = MergeStringList(filesToZip, ScanFolderAndFilter(specFile.Files, "./src/"))

		ZipMe(filesToZip, pack)

		log.Println(depth + "├─ Upload Package")

		err = UploadFile(result.FilePath, *response)
		if err != nil {
			panic(err)
		}

		log.Println(depth + "├─ Successfully uploaded")

		err = os.Remove(result.FilePath)
		if err != nil {
			panic(err)
		}

		for _, dependency := range result.Dependencies {
			log.Println(depth + "├─ Handling dependency")

			upload(url, dependency.Name, dependency.Vendor, "|	"+depth)
		}

		log.Println(depth + "└─ Finished packing: " + packageName)

	} else {
		log.Println(depth + "└─ Skipped. Already present. Use update if you want to replace it")
	}
}

func askOperatorForProcedure(data []MetaData) bool {

	fmt.Println("Found fhe following packages with similar or same content")

	for _, d := range data {
		fmt.Println(d.String())
	}

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	for !acceptInput(text) {
		fmt.Print("Do you want to upload your version anyway? ")
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func acceptInput(text string) bool {

	var acceptedAnswers = []string{"yes", "y", "no", "n"}

	for _, s := range acceptedAnswers {

		if strings.ToLower(text) == s {
			return true
		}
	}

	fmt.Println("Please enter yes / y or no / n")
	return false
}

func readDependencies(specFile SpecFile, vendor string) []Dependency {

	var dependencies []Dependency

	for _, d := range specFile.Dependencies {
		dependencySpec := ReadSpec("./packages/" + d)

		dependencies = append(dependencies, Dependency{
			Name:    dependencySpec.Name,
			Version: dependencySpec.Version,
			Vendor:  vendor}, )
	}
	return dependencies
}
