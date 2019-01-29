package helpers

import (
	"bufio"
	"fmt"
	. "github.com/evoila/BPM-Client/model"
	. "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

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

func ReadAndValidateSpec(packageName string) (*SpecFile, *string) {

	yamlFile, err := ioutil.ReadFile("./packages/" + packageName + "/spec")
	if err != nil {
		message := "Did not find a spec file. Is '" + packageName + "' a valid package?"

		return nil, &message
	}

	var specFile SpecFile

	err = Unmarshal(yamlFile, &specFile)

	if specFile.Name == "" || specFile.Version == "" || specFile.Vendor == ""{
		message := "The Specfile needs to specify package, version and vendor."
		return nil, &message
	}


	if err != nil {
		message := "'" + packageName + "' does not contain a valid spec file."

		return nil, &message
	}

	return &specFile, nil
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

func AskUser(data MetaData, depth, message string) bool {

	fmt.Println(depth + "├─ Update Package")
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