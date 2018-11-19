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

	switch response.StatusCode {

	case 202:
		{
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			var converted S3Permission
			err = Unmarshal(responseBody, &converted)

			return &converted, nil
		}
	case 409:
		{
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			var metaData MetaData
			err = Unmarshal(responseBody, &metaData)

			return nil, &metaData
		}
	case 401:
		{
			return nil, nil
		}
	default:
		{
			panic("A unexpected response code: " + strconv.Itoa(response.StatusCode))
		}
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

	client := &Client{}

	request, err := NewRequest("GET", BuildPath([]string{config.Url, "package", vendor, name, version}), nil)
	request.Header.Set("Content-Type", "application/json")

	if config.Username != "" && config.Password != "" {
		request.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

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

func GetDownloadPermission(config *Config, requestBody PackageRequestBody) *S3Permission {

	path := BuildPath([]string{config.Url, "download", requestBody.Vendor, requestBody.Name, requestBody.Version})
	client := &Client{}

	request, err := NewRequest("GET", path, nil)
	request.Header.Set("Content-Type", "application/json")

	if config.Username != "" && config.Password != "" {
		request.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(request)

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
