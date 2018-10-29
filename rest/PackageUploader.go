package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/evoila/BPM-Client/model"
)

func PutMetaData(data model.MetaData) ResponsetBody {

	client := &http.Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("PUT", url+"upload/package", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	responsetBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var converted ResponsetBody
	err = json.Unmarshal(responsetBody, &converted)

	return converted
}

func buildBody(data model.MetaData) ([]byte, error) {

	requestBody := requestBody{
		Name:    data.Name,
		Version: data.Version,
		Vendor:  data.Vendor,
		Files:   data.Files}

	return json.Marshal(requestBody)
}

type requestBody struct {
	Name, Version, Vendor string
	Files                 []string
}

type ResponsetBody struct {
	Name, Version, Vendor, S3location string
	Files                             []string
}

const url = "http://localhost:8080/"
