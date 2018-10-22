package main

import (
	"os"

	"github.com/evoila/BPM-Client/rest/bundle"
)

func main() {
	moveToReleaseDir()

	bundle.ZipPackage(packageName)
}

func moveToReleaseDir() {
	/*	dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
	*/
	dir := "/home/johannes/workspace/osb-bosh-kafka"
	os.Chdir(dir)
}

const packageName = "kafka-smoke-test"
const url = "http://localhost:8080"
