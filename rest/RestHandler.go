package rest

import (
	. "bytes"
	. "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	. "net/http"

	. "github.com/evoila/BPM-Client/model"
)

func PutMetaData(data MetaData) UploadPermission {

	client := &Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}
	request, err := NewRequest("PUT", url+"upload/package", NewBuffer(body))
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

	var converted UploadPermission
	err = Unmarshal(responseBody, &converted)

	return converted
}

func GetMetaData(request packageRequestBody) ResponseBody {

	resp, err := Get(BuildPath([]string{url, request.Vendor, request.Name, request.Version}))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var responsePackage ResponseBody
	err = Unmarshal(responseBody, &responsePackage)

	if err != nil {
		panic(err)
	}

	return responsePackage
}

type packageRequestBody struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor"`
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

const url = "http://localhost:8080/"
