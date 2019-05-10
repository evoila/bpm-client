package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/blang/semver"
	. "github.com/evoila/BPM-Client/model"
	. "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func SplitPackageReference(input string) (string, string, string, error) {
	result := strings.Split(input, ":")

	if len(result) != 3 {
		return "", "", "", errors.New("invalid input")
	}

	return result[0], result[1], result[2], nil
}

func MergeStringList(l1, l2 []string) []string {

	if len(l2) > len(l1) {
		return MergeStringList(l2, l1)
	}

	for _, s := range l2 {
		l1 = append(l1, s)
	}

	return l1
}

func MergeMetaDataList(l1, l2 []MetaData) []MetaData {

	if len(l2) > len(l1) {
		return MergeMetaDataList(l2, l1)
	}

	for _, s := range l2 {
		l1 = append(l1, s)
	}

	return l1
}

func ReadDownloadSpec() (*DownloadSpec, string) {

	yamlFile, err := ioutil.ReadFile("./packages/download.spec")
	if err != nil {
		message := "Did not find a download.spec file."

		return nil, message
	}

	var downloadSpec DownloadSpec
	err = Unmarshal(yamlFile, &downloadSpec)

	if err != nil {
		message := "The download spec is not valid."
		return nil, message
	}

	return &downloadSpec, ""
}

func ReadAndValidateSpec(packageName string) (*SpecFile, string) {

	yamlFile, err := ioutil.ReadFile("./packages/" + packageName + "/spec")
	if err != nil {
		message := "Did not find a spec file. Is '" + packageName + "' a valid package?"

		return nil, message
	}

	var specFile SpecFile
	err = Unmarshal(yamlFile, &specFile)

	if specFile.Name == "" || specFile.Version == "" || specFile.Publisher == "" {
		message := "The Specfile needs to specify package, version and vendor."

		return nil, message
	}

	if validateVersion(specFile.Version) {
		message := "The Version of the package '" + packageName + "' is not Semver conform."

		return nil, message
	}

	if specFile.Stemcell.MinorVersion != "" {
		_, err = semver.Make(specFile.Stemcell.MinorVersion)

		if validateVersion(specFile.Stemcell.MinorVersion) {
			message := "The Stemcell Minor Version of the package '" + packageName + "' is not Semver conform."

			return nil, message
		}

	}

	if specFile.Stemcell.MinorVersion != "" {
		_, err = semver.Make(specFile.Stemcell.MajorVersion)

		if validateVersion(specFile.Stemcell.MajorVersion) {
			message := "The Stemcell Major Version of the package '" + packageName + "' is not Semver conform."

			return nil, message
		}
	}

	return &specFile, ""
}

func validateVersion(version string) bool {
	_, err := semver.Make(version)

	return err != nil
}

func WriteConfig(config Config, configLocation string) {
	data, err := Marshal(config)

	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configLocation, data, 0644)

	if err != nil {
		panic(err)
	}
}

func ReadConfig(configLocation string) Config {

	yamlFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		panic(err)
	}

	var config Config

	err = Unmarshal(yamlFile, &config)

	if err != nil {
		panic(err)
	}

	return config
}

func BuildPath(path []string) string {

	return strings.Join(path, "/")
}

func MoveToReleaseDir() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(dir)
}

func AskUser(data MetaData, depth, operation, message string) bool {

	//"Update Package"

	fmt.Println(depth + "├─ " + operation)
	fmt.Println(data.String(depth))

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Println(depth + message)

	for !AcceptInput(text, depth) {
		scanner.Scan()
		text = scanner.Text()
	}

	return strings.ToLower(text) == "yes" || strings.ToLower(text) == "y"
}

func AcceptInput(text, depth string) bool {

	var acceptedAnswers = []string{"yes", "y", "no", "n"}

	for _, s := range acceptedAnswers {

		if strings.ToLower(text) == s {
			return true
		}
	}
	fmt.Print(depth + "│  Please enter yes / y or no / n	Answer: ")

	return false
}
