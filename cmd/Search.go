package cmd

import (
	"fmt"
	"github.com/Nerzal/gocloak"
	"github.com/evoila/BPM-Client/rest"
	"strconv"
)
import . "github.com/evoila/BPM-Client/model"

func SearchByPublisher(publisher string, config *Config, openId *gocloak.JWT, more bool) {

	body, statusCode := rest.GetMetaDataListByPublisher(config, openId, publisher)

	if statusCode == 200 {
		if more {
			printResultMore(body.Embedded.Packages)
		} else {
			printResultLess(body.Embedded.Packages)
		}
	} else {
		fmt.Println("No Packages Found. Status: " + strconv.Itoa(statusCode))
	}
}

func SearchByPublisherAndName(publisher, name string, config *Config, openId *gocloak.JWT, more bool) {
	body, statusCode := rest.GetMetaDataListByPublisherAndName(config, openId, publisher, name)

	if statusCode == 200 {
		if more {
			printResultMore(body.Embedded.Packages)
		} else {
			printResultLess(body.Embedded.Packages)
		}
	} else {
		fmt.Println("No Packages Found. Status: " + strconv.Itoa(statusCode))
	}
}

func printResultLess(metaData []MetaData) {
	for _, d := range metaData {
		fmt.Println("├─" + d.Publisher + ":" + d.Name + ":" + d.Version)
	}
}

func printResultMore(metaData []MetaData) {

	for _, d := range metaData {
		fmt.Println(d.String2())
		fmt.Println("├─")
	}
}