package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/evoila/BPM-Client/model"
)

const Url = "http://localhost:8080"
const packages = "/package/6e3ccc11-630e-4a04-a8ac-f9e3d3d60f2a"
const blobPath = "/blobs/"

func GetPackages(url string) {

	path := BuildPath([]string{url, packages})
	resp, err := http.Get(path)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	var responsePackages []model.ResponseBoshPackage
	err = json.Unmarshal(body, &responsePackages)

	if err != nil {
		panic(err)
	}

	for _, responsePackage := range responsePackages {
		for _, blobId := range responsePackage.Blobs {
			var blob = getBlobMeta(url, blobId)
			DownloadBlob(url, blobId, blob)
		}
	}
}

func getBlobMeta(url, uuid string) model.BoshBlob {

	path := BuildPath([]string{url, blobPath, uuid})

	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	var responseBlob model.BoshBlob
	err = json.Unmarshal(body, &responseBlob)

	if err != nil {
		panic(err)
	}

	return responseBlob
}
