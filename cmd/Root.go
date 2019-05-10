package cmd

import (
	"fmt"
	"github.com/Nerzal/gocloak"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var config Config
var pack, version, publisher, accessLevel string
var update, force bool

var rootCmd = &cobra.Command{
	Use:   "BPM-Client",
	Short: "CLI Tool to access Bosh Package Manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one command of: upload," +
			" download, delete, create-release, search, publish or create-publisher")
	},
}

func init() {
	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a package to Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Begin upload.")
			setupConfig()
			jwt, err := SetUsernamePasswordIfNewAndPerformLogin(&config)

			if err != nil {
				log.Println("└─ Unauthorized. Upload canceled.")
				return
			}

			if update {
				RunUpdateIfPackagePresentUploadIfNot(pack, &config, jwt)
			} else {
				CheckIfAlreadyPresentAndUpload(pack, &config, jwt)
			}

			log.Println("Finished upload.")
		},
	}
	uploadCmd.Flags().StringVarP(&pack, "package", "n", "", "The name of the package to upload")
	_ = uploadCmd.MarkFlagRequired("package")
	uploadCmd.Flags().BoolVar(&update,
		"update",
		false,
		"Set if you want tp update packages.")

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a package with all its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()

			var jwt *gocloak.JWT
			var err error

			if config.Username != "" && config.Password != "" {
				jwt, err = rest.Login(&config)
			}

			if err != nil {
				log.Println("└─ Unauthorized. Download canceled.")
			}

			var requestBody PackagesReference

			if pack == "" || publisher == "" || version == "" {
				publisher, pack, version, err = helpers.SplitPackageReference(args[0])

				if err != nil {
					log.Print("Invalid input.")
				}
			}

			requestBody = PackagesReference{
				Name:      pack,
				Publisher: publisher,
				Version:   version}

			DownloadPackageWithDependencies("", requestBody, &config, jwt)
		},
	}
	downloadCmd.Flags().StringVarP(&pack, "package", "n", "", "The name of the package")
	//	_ = downloadCmd.MarkFlagRequired("package")
	downloadCmd.Flags().StringVarP(&publisher, "publisher", "p", "", "The name of the publisher")
	//	_ = downloadCmd.MarkFlagRequired("publisher")
	downloadCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package")
	//	_ = downloadCmd.MarkFlagRequired("version")

	var createRelease = &cobra.Command{
		Use:   "create-release",
		Short: "DownloadPackageWithDependencies a package with all dependencies from Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			downloadSpec, errMessage := helpers.ReadDownloadSpec()

			log.Println("Downloading packages based on download.spec")

			if errMessage != "" {
				log.Print(errMessage)
			}

			var jwt *gocloak.JWT
			var err error

			if config.Username != "" && config.Password != "" {
				jwt, err = rest.Login(&config)
			}

			if err != nil {
				log.Println("└─ Invalid Credentials.")
			}

			DownloadBySpec(*downloadSpec, &config, jwt)
		},
	}

	var loginTest = &cobra.Command{
		Use:   "login",
		Short: "test credentials",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()

			log.Println("Testing login for " + config.Username)
			jwt, err := SetUsernamePasswordIfNewAndPerformLogin(&config)

			if err == nil {
				rest.BackendLogin(&config, jwt)
			} else {
				log.Println("login failed.")
			}
		},
	}

	var createPublisher = &cobra.Command{
		Use:   "create-publisher",
		Short: "creates a new publisher and adds you to it as a member",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			if publisher == "" {
				log.Println("Please specify a name for the new publisher.")
				return
			}
			jwt, err := SetUsernamePasswordIfNewAndPerformLogin(&config)

			if err == nil {
				rest.CreatePublisher(&config, publisher, jwt)
			} else {
				log.Println("login failed.")
			}
		},
	}
	createPublisher.Flags().StringVarP(&publisher, "publisher", "p", "", "The name of the publisher")
	var publishPackage = &cobra.Command{
		Use:   "publish-package",
		Short: "publish a package you own",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			if config.Username == "" && config.Password == "" {
				log.Println("Please set your username and password in the " +
					"config file and reference it via path variable")
				return
			}
			openId, err := rest.Login(&config)

			if err == nil {

				var requestBody PackagesReference

				if pack == "" || publisher == "" || version == "" {
					publisher, pack, version, err = helpers.SplitPackageReference(args[0])

					if err != nil {
						log.Print("Invalid input.")
					}
				}

				requestBody = PackagesReference{
					Name:      pack,
					Publisher: publisher,
					Version:   version}
				Publish(requestBody, accessLevel, &config, openId, force)
			} else {
				log.Println("login failed.")
			}
		},
	}
	publishPackage.Flags().StringVarP(&publisher, "publisher", "p", "", "The name of the publisher")
	publishPackage.Flags().StringVarP(&pack, "package", "n", "", "The name of the package")
	publishPackage.Flags().StringVarP(&version, "version", "v", "", "Version of the package")
	publishPackage.Flags().StringVarP(&accessLevel, "access-level", "a", "", "The desired access level. Either publisher or public")
	_ = publishPackage.MarkFlagRequired("access-level")
	publishPackage.Flags().BoolVarP(&force, "force", "f", false, "Set this flag to skip all prompts")

	var search = &cobra.Command{
		Use:   "search",
		Short: "search packages by a given publisher",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			var jwt *gocloak.JWT
			var err error

			if config.Username != "" && config.Password != "" {
				jwt, err = rest.Login(&config)
			}
			if err != nil {
				log.Println("login failed.")
			}

			if pack != "" {
				SearchByPublisherAndName(publisher, pack, &config, jwt)
			} else {
				SearchByPublisher(publisher, &config, jwt)
			}
		},
	}
	search.Flags().StringVarP(&pack, "name", "n", "", "The name of the publisher")
	search.Flags().StringVarP(&publisher, "publisher", "p", "", "The name of the publisher")
	_ = search.MarkFlagRequired("publisher")

	var register = &cobra.Command{
		Use:   "register",
		Short: "register a new user",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			Register(&config)
		},
	}

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(loginTest)
	rootCmd.AddCommand(createPublisher)
	rootCmd.AddCommand(publishPackage)
	rootCmd.AddCommand(search)
	rootCmd.AddCommand(register)
	rootCmd.AddCommand(createRelease)
}

func setupConfig() {
	configLocation := os.Getenv("BOSH_PACKAGE_MANAGER_CONFIG")
	config = helpers.ReadConfig(configLocation)
	helpers.MoveToReleaseDir()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
