package cmd

import (
	. "github.com/evoila/BPM-Client/bundle"
	. "github.com/evoila/BPM-Client/collections"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	. "github.com/evoila/BPM-Client/rest"
	. "github.com/evoila/BPM-Client/s3"
	"log"
	"os"
)

var set = NewStringSet()

func Upload(url, packageName, vendor, depth string) {

	if set.Get(packageName) {
		log.Println(depth + "├─  Dependency " + packageName + " already handeled.")
		return
	}

	set.Add(packageName)

	log.Println(depth + "├─ Packing: " + packageName)

	specFile := ReadSpec("./packages/" + packageName)
	pack := "./" + packageName + ".bpm"

	result := MetaData{
		Name:         packageName,
		Version:      specFile.Version,
		Vendor:       vendor,
		FilePath:     pack,
		Files:        specFile.Files,
		Dependencies: specFile.Dependencies}

	var response = PutMetaData(url, result, false)

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

		for _, dependency := range specFile.Dependencies {
			if _, err := os.Stat("./" + dependency + ".bpm"); os.IsNotExist(err) {
				log.Println(depth + "├─ Handling dependency")

				Upload(url, dependency, vendor, "|	"+depth)
			}
		}

	} else {
		log.Println(depth + "Skipping " + packageName + ". Reusing present one.")
	}

	log.Println(depth + "└─ Finished packing: " + packageName)
}
