package rest

import (
	. "bytes"
	. "encoding/json"
	"github.com/Nerzal/gocloak"
	. "github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"io/ioutil"
	"log"
	. "net/http"
	"strconv"
)

func RequestPermission(data MetaData, force bool, config *Config, openId *gocloak.JWT) (*S3Permission, *MetaData) {

	client := &Client{}

	body, err := buildBody(data)

	if err != nil {
		panic(err)
	}

	request, err := NewRequest("POST", BuildPath([]string{config.Url, "upload/permission?force=" + strconv.FormatBool(force)}), NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "bearer "+openId.AccessToken)

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

func PutMetaData(url, location string, openId *gocloak.JWT, size int64) {

	path := BuildPath([]string{url, "package?location=" + location + "&size=" + strconv.FormatInt(size, 10)})
	request, err := NewRequest("PUT", path, nil)
	request.Header.Set("Content-Type", "application/json")

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
	}

	client := &Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
}

func GetMetaData(publisher, name, version string, config *Config, openId *gocloak.JWT) *MetaData {

	client := &Client{}

	request, err := NewRequest("GET", BuildPath([]string{config.Url, "package", publisher, name, version}), nil)
	request.Header.Set("Content-Type", "application/json")

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
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

func GetMetaDataListForPackageName(name string, config *Config, openId *gocloak.JWT) []MetaData {
	request, err := NewRequest("GET", BuildPath([]string{config.Url, "package?name=" + name}), nil)

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
	}
	client := &Client{}
	response, err := client.Do(request)

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

func GetDownloadPermission(config *Config, requestBody PackageRequestBody, openId *gocloak.JWT) *S3Permission {

	path := BuildPath([]string{config.Url, "download", requestBody.Publisher, requestBody.Name, requestBody.Version})
	client := &Client{}

	request, err := NewRequest("GET", path, nil)
	request.Header.Set("Content-Type", "application/json")

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
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

func CreatePublisher(config *Config, name string, jwt *gocloak.JWT) {

	path := BuildPath([]string{config.Url, "publishers?name=" + name})
	request, _ := NewRequest("POST", path, nil)
	request.Header.Set("Authorization", "bearer "+jwt.AccessToken)

	client := &Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	if response.StatusCode == 201 {
		log.Println("Publisher " + name + " created.")
		return
	}

	log.Println("Expected 200 but was " + strconv.Itoa(response.StatusCode))
}

func PublishPackage(id, accessLevel string, config *Config, openId *gocloak.JWT) bool {
	path := BuildPath([]string{config.Url, "publish", id + "?access=" + accessLevel})
	request, _ := NewRequest("PATCH", path, nil)
	request.Header.Set("Authorization", "bearer "+openId.AccessToken)

	client := &Client{}
	response, _ := client.Do(request)

	return response.StatusCode == 200
}

func GetMetaDataListByPublisher(config *Config, openId *gocloak.JWT, publisher string) (*PaginatedMetaData, int) {
	path := BuildPath([]string{config.Url, "rest", "packages", publisher})
	request, _ := NewRequest("GET", path, nil)

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
	}
	client := &Client{}
	response, _ := client.Do(request)

	if response.StatusCode == 200 {
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var paginatedMetaData PaginatedMetaData
		err = Unmarshal(responseBody, &paginatedMetaData)
		return &paginatedMetaData, response.StatusCode
	}

	return nil, response.StatusCode
}

func GetMetaDataListByPublisherAndName(config *Config, openId *gocloak.JWT, publisher, name  string) (*PaginatedMetaData, int) {
	path := BuildPath([]string{config.Url, "rest", "packages", publisher, name})
	request, _ := NewRequest("GET", path, nil)

	if openId != nil {
		request.Header.Set("Authorization", "bearer "+openId.AccessToken)
	}
	client := &Client{}
	response, _ := client.Do(request)

	if response.StatusCode == 200 {
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var paginatedMetaData PaginatedMetaData
		err = Unmarshal(responseBody, &paginatedMetaData)
		return &paginatedMetaData, response.StatusCode
	}

	return nil, response.StatusCode
}

func Login(config *Config) (*gocloak.JWT, error) {

	goCloak := gocloak.NewClient(config.KeycloakConfig.Url)

	token, err := goCloak.Login(config.KeycloakConfig.ClientID,
		"secret",
		config.KeycloakConfig.Realm,
		config.Username,
		config.Password)

	return token, err
}

func BackendLogin(config *Config, openId *gocloak.JWT) {
	path := BuildPath([]string{config.Url, "login"})
	request, _ := NewRequest("PUT", path, nil)
	request.Header.Set("Authorization", "bearer "+openId.AccessToken)
	client := &Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	if response.StatusCode == 200 {

		log.Println("You're authorized")
		return
	} else if response.StatusCode == 201 {
		log.Println("Welcome to Bosh Package Manager")
		return
	}

	log.Println("Expected 200 or 201 but was " + strconv.Itoa(response.StatusCode))
}

func buildBody(data MetaData) ([]byte, error) {
	requestBody := requestBody{
		Name:         data.Name,
		Version:      data.Version,
		Publisher:    data.Publisher,
		Files:        data.Files,
		Dependencies: data.Dependencies,
		Description:  data.Description,
		Stemcell:     data.Stemcell,
		Url:          data.Url,
	}

	return Marshal(requestBody)
}

type requestBody struct {
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Publisher    string              `json:"publisher"`
	Description  string              `json:"description"`
	Files        []string            `json:"files"`
	Dependencies []PackagesReference `json:"dependencies"`
	Stemcell     Stemcell            `json:"stemcell"`
	Url          string              `json:"url"`
}
