package rest

import (
	. "bytes"
	. "encoding/json"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"io/ioutil"
	"log"
	. "net/http"
	"strconv"
)

func PutMetaData(url string, data MetaData, force bool) (*S3Permission, *MetaData) {

	client := &Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}

	request, err := NewRequest("PUT", BuildPath([]string{url, "upload/package?force=" + strconv.FormatBool(force)}), NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode == 202 {

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var converted S3Permission
		err = Unmarshal(responseBody, &converted)

		return &converted, nil

	} else if response.StatusCode == 409 {

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var metaData MetaData
		err = Unmarshal(responseBody, &metaData)

		return nil, &metaData
	} else {
		panic("A unexpected response code: " + strconv.Itoa(response.StatusCode))
	}
}

func GetMetaData(url, vendor, name, version string) *MetaData {

	path := BuildPath([]string{url, "package", vendor, name, version})

	resp, err := Get(path)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var metaData MetaData
	err = Unmarshal(responseBody, &metaData)

	if err != nil {
		panic(err)
	}

	return &metaData

}

func GetMetaDataListForPackageName(url, name string) []MetaData {

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

func GetDownloadPermission(url string, request PackageRequestBody) *S3Permission {

	resp, err := Get(BuildPath([]string{url, "download", request.Vendor, request.Name, request.Version}))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var permission S3Permission
	err = Unmarshal(responseBody, &permission)

	if err != nil {
		panic(err)
	}

	return &permission
}

func buildBody(data MetaData) ([]byte, error) {

	requestBody := requestBody{
		Name:         data.Name,
		Version:      data.Version,
		Vendor:       data.Vendor,
		Files:        data.Files,
		Dependencies: data.Dependencies}

	return Marshal(requestBody)
}

type requestBody struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Vendor       string       `json:"vendor"`
	Files        []string     `json:"files"`
	Dependencies []Dependency `json:"dependencies"`
}
