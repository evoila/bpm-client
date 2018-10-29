package execute

import (
	. "github.com/evoila/BPM-Client/bundle"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
	"os"
)

func Upload(packageName, vendor, version string) {

	result := ZipPackage(packageName, version, vendor, "")

	for _, r := range result {
		response := rest.PutMetaData(r)

		err := s3.UploadFile(r.FilePath, response)
		if err != nil {
			panic(err)
		}

		err = os.Remove(r.FilePath)
		if err != nil {
			panic(err)
		}
	}
}
