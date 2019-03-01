package cmd

import (
	"bufio"
	"fmt"
	"github.com/Nerzal/gocloak"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func Register(config *Config) {
	fmt.Println("Please visit the following page to register")
	path := helpers.BuildPath([]string{config.KeycloakConfig.Url,
		"auth/realms",
		config.KeycloakConfig.Realm,
		"protocol/openid-connect/registrations?client_id=" +
			config.KeycloakConfig.ClientID + "&response_type=code&redirect_uri=http://localhost:4200/packages",
	})
	fmt.Println(path)
}

func SetUsernamePasswordIfNewAndPerformLogin(config *Config) (*gocloak.JWT, error) {
	var jwt *gocloak.JWT
	var err error

	if config.Username == "" || config.Password == "" {
		var username = enterUsername()
		var password = enterPassword()
		config.Username = username
		config.Password = password
		jwt, err = rest.Login(config)

		if err != nil {
			configLocation := os.Getenv("BOSH_PACKAGE_MANAGER_CONFIG")
			helpers.WriteConfig(config, configLocation)
		}
	} else {

		jwt, err = rest.Login(config)
	}

	return jwt, err
}

func enterUsername() string {
	fmt.Println("Enter your Username")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()

		if input != "" {
			return input
		}
	}
}

func enterPassword() string {
	fmt.Println("Enter your password")
	password, err := terminal.ReadPassword(0)

	if err != nil {
		panic(err)
	}
	return string(password)
}
