package rest

import (
	. "bytes"
	. "encoding/json"
	"io/ioutil"
	"log"
	. "net/http"
	"strconv"

	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
)

func RequestPermission(data MetaData, force bool, config *Config) (*S3Permission, *MetaData) {

	client := &Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}

	request, err := NewRequest("POST", BuildPath([]string{config.Url, "upload/permission?force=" + strconv.FormatBool(force)}), NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.Username, config.Password)

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
	} else if response.StatusCode == 401 {
		return nil, nil

	} else {
		panic("A unexpected response code: " + strconv.Itoa(response.StatusCode))
	}
}

func PutMetaData(url, location string) {

	path := BuildPath([]string{url, "package?location=" + location})
	request, err := NewRequest("PUT", path, nil)

	request.Header.Set("Content-Type", "application/json")

	client := &Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
}

func GetMetaData(vendor, name, version string, config *Config) *MetaData {

	path := BuildPath([]string{config.Url, "package", vendor, name, version})

	resp, err := Get(path)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		responseBody, _ := ioutil.ReadAll(resp.Body)

		var metaData MetaData
		err = Unmarshal(responseBody, &metaData)

		if err != nil {
			panic(err)
		}

		return &metaData
	} else {
		return nil
	}
}

func GetMetaDataListForPackageName(name string, config *Config) []MetaData {

	path := BuildPath([]string{config.Url, "package?name=" + name})

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

func GetDownloadPermission(config *Config, request PackageRequestBody) *S3Permission {

	resp, err := Get(BuildPath([]string{config.Url, "download", request.Vendor, request.Name, request.Version}))

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
