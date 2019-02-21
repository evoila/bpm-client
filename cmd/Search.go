package cmd

import (
	"fmt"
	"github.com/evoila/BPM-Client/rest"
	"strconv"
)
import . "github.com/evoila/BPM-Client/model"

func SearchByVendor(vendor string, config *Config, openId *OpenId) {
	body, statusCode := rest.GetMetaDataListByVendor(config, openId, vendor)

	if statusCode == 200 {
		var metaData = body.Embedded.Packages

		for _, d := range metaData {
			fmt.Println(d.String2())
			fmt.Println("├─")
		}
	} else {
		fmt.Println("No Packages Found. Status: " + strconv.Itoa(statusCode))
	}
}
