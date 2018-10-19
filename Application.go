package main

import (
	"github.com/evoila/BPM-Client/rest"
)

func main() {

	//	rest.DownloadBlob(url, "83204f44-fee1-4ecf-b973-d22c86fdeb62", "openjdk.tar.gz")
	rest.GetPackages(url)
}

const url = "http://localhost:8080"
