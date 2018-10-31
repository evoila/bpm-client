package rest

import (
	"bufio"
	. "bytes"
	. "encoding/json"
	"fmt"
	. "github.com/evoila/BPM-Client/helpers"
	"io/ioutil"
	"log"
	. "net/http"
	"os"
	"strconv"
	"strings"

	. "github.com/evoila/BPM-Client/model"
)

func PutMetaData(url string, data MetaData, force bool) *S3Permission {

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

		return &converted

	} else if response.StatusCode == 409 {

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var metaData []MetaData
		err = Unmarshal(responseBody, &metaData)

		fmt.Println("At least one package named " + data.Name + " already exists:")

		if askOperatorForProcedure(metaData) {
			return PutMetaData(url, data, true)
		}
	}

	return nil
}

func askOperatorForProcedure(data []MetaData) bool {

	for _, d := range data {
		fmt.Println(d.String())
	}

	scanner := bufio.NewScanner(os.Stdin)
	var text string

	for !acceptInput(text) {
		fmt.Print("Do you want to upload your own package anyway? ")
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func acceptInput(text string) bool {

	var acceptedAnswers = []string{"yes", "y", "no", "n"}

	for _, s := range acceptedAnswers {

		if strings.ToLower(text) == s {
			return true
		}
	}

	fmt.Println("Please enter yes / y or no / No")
	return false
}

func GetMetaDataForPackageName(url, name string) []MetaData {

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

func GetDownloadPermission(url string, request PackageRequestBody) S3Permission {

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
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Vendor       string   `json:"vendor"`
	Files        []string `json:"files"`
	Dependencies []string `json:"dependencies"`
}
