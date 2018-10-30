package helpers

import (
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

func ReadSpec(specLocation string) SpecFile {

	yamlFile, err := ioutil.ReadFile(specLocation + "/spec")
	if err != nil {
		panic(err)
	}

	var specFile SpecFile

	err = Unmarshal(yamlFile, &specFile)

	if err != nil {
		panic(err)
	}

	return specFile
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
