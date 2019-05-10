package cmd

import (
	"fmt"
	"github.com/Nerzal/gocloak"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	. "github.com/evoila/BPM-Client/rest"
	"strings"
)

func Publish(packageReference PackagesReference, accessLevelInput string, config *Config, jwt *gocloak.JWT, force bool) {
	var meta = GetMetaData(packageReference, config, jwt)

	if meta == nil {
		fmt.Println("Package not found. Aborting.")
		return
	}

	var accessLevel = validateAndCorrectAccessLevelInput(accessLevelInput)

	if accessLevel == nil {
		fmt.Println("Not a valid AccessType. Please enter 'publisher' or 'public'.")
		return
	}

	if force || helpers.AskUser(*meta, "", "Publish Package", "├─ The package "+meta.Name+" and all it's dependencies by your vendors will be published. Are you sure?") {

		if PublishPackage(meta.Id, *accessLevel, config, jwt) {
			fmt.Println("Package '" + packageReference.String() + "' published.")
		} else {
			fmt.Println("Something went wrong!")
		}
	} else {
		fmt.Println("Aborting.")
	}
}

func validateAndCorrectAccessLevelInput(accessLevelInput string) *string {

	var result = strings.ToUpper(accessLevelInput)

	if result != "PUBLISHER" && result != "PUBLIC" {
		return nil
	}

	return &result
}
