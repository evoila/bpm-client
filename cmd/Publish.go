package cmd

import (
	"fmt"
	"github.com/Nerzal/gocloak"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	. "github.com/evoila/BPM-Client/rest"
	"strings"
)

func Publish(vendor, name, version string, accessLevelInput string, config *Config, jwt *gocloak.JWT, force bool) {
	var meta = GetMetaData(vendor, name, version, config, jwt)

	if meta == nil {
		fmt.Println("Package not found. Aborting.")
		return
	}

	var accessLevel = validateAndCorrectAccessLevelInput(accessLevelInput)

	if accessLevel == nil {
		fmt.Println("Not a valid AccessType. Please enter 'vendor' or 'public'.")
		return
	}

	if force || helpers.AskUser(*meta, "", "The package "+meta.Name+" and all it's dependencies by your vendors will be published. Are you sure?") {

		if PublishPackage(meta.Id, *accessLevel, config, jwt) {
			fmt.Println("Package published.")
		} else {
			fmt.Println("Something went wrong!")
		}

	} else {
		fmt.Println("Aborting.")
	}
}

func validateAndCorrectAccessLevelInput(accessLevelInput string) *string {

	var result = strings.ToUpper(accessLevelInput)

	if result != "VENDOR" && result != "PUBLIC" {
		return nil
	}

	return &result
}
