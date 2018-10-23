package main

import (
	"os"

	"github.com/evoila/BPM-Client/bundle"
)

func main() {
	moveToReleaseDir()

	result := bundle.ZipPackage(packageName, "1.0", "Myself", "")

	for _, r := range result {
		error := bundle.UnzipPackage(r.FilePath, "./test")
		if error != nil {
			panic(error)
		}

	}
}

func moveToReleaseDir() {
	/*dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}*/

	os.Chdir(dir)
}

const dir = "/home/johannes/workspace/osb-bosh-kafka"
const packageName = "kafka-smoke-test"
const url = "http://localhost:8080"
