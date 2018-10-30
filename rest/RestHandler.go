package rest

import (
	. "bytes"
	. "encoding/json"
	"fmt"
	. "github.com/evoila/BPM-Client/helpers"
	"io/ioutil"
	"log"
	. "net/http"

	. "github.com/evoila/BPM-Client/model"
)

func PutMetaData(data MetaData) S3Permission {

	client := &Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}
	request, err := NewRequest("PUT", BuildPath([]string{url + "upload/package"}), NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var converted S3Permission
	err = Unmarshal(responseBody, &converted)

	return converted
}

func GetMetaDataForPackageName(name string) []MetaData {

	path := BuildPath([]string{url, "package?name=" + name})

	response, err := Get(path)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var metaData []MetaData

	err = Unmarshal(responseBody, &metaData)

	return metaData
}

func GetDownloadPermission(request PackageRequestBody) S3Permission {

	resp, err := Get(BuildPath([]string{url, request.Vendor, request.Name, request.Version}))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var permission S3Permission
	err = Unmarshal(responseBody, &permission)

	if err != nil {
		panic(err)
	}

	return permission
}

func buildBody(data MetaData) ([]byte, error) {

	requestBody := requestBody{
		Name:    data.Name,
		Version: data.Version,
		Vendor:  data.Vendor,
		Files:   data.Files}

	return Marshal(requestBody)
}

type requestBody struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Vendor  string   `json:"vendor"`
	Files   []string `json:"files"`
}

const url = "http://localhost:8080"
