package cmd

import (
	. "github.com/evoila/BPM-Client/bundle"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
	"os"
)

func Upload(url, packageName, vendor string) {

	result := ZipPackage(packageName, vendor, "")

	for _, r := range result {

		var response = rest.PutMetaData(url, r, false)

		if response != nil {
			err := s3.UploadFile(r.FilePath, *response)
			if err != nil {
				panic(err)
			}
		}

		err := os.Remove(r.FilePath)
		if err != nil {
			panic(err)
		}
	}
}
