package cmd

import (
	"fmt"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "BPM-Client",
	Short: "CLI Tool to access Bosh Package Manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one command of: upload, update, download, delete, search, publish or create-vendor")
	},
}

var config Config
var pack, version, vendor, accessLevel string
var update, force bool

func init() {

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a package to Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			log.Println("Begin upload.")

			openId := rest.Login(&config)

			if openId == nil {
				log.Println("└─ Unauthorized. Upload canceled.")
			}

			if update {
				RunUpdateIfPackagePresentUploadIfNot(pack, &config, openId)
			} else {
				CheckIfAlreadyPresentAndUpload(pack, &config, openId)
			}

			log.Println("Finished upload.")
		},
	}
	uploadCmd.Flags().StringVarP(&pack, "package", "p", "", "The name of the package to upload")
	uploadCmd.MarkFlagRequired("package")
	uploadCmd.Flags().BoolVar(&update, "update", false, "Set if you want tp update packages.")

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a package with all dependencies from Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()

			var openId *OpenId

			if config.Username != "" && config.Password != "" {
				openId = rest.Login(&config)
			}

			requestBody := PackageRequestBody{
				Name:    pack,
				Vendor:  vendor,
				Version: version}

			Download("", requestBody, &config, openId)
		},
	}
	downloadCmd.Flags().StringVarP(&pack, "package", "p", "", "The name of the package")
	downloadCmd.MarkFlagRequired("package")
	downloadCmd.Flags().StringVarP(&vendor, "vendor", "v", "", "The name of the vendor")
	downloadCmd.MarkFlagRequired("vendor")
	downloadCmd.Flags().StringVarP(&version, "version", "s", "", "Version of the package")
	downloadCmd.MarkFlagRequired("version")

	var loginTest = &cobra.Command{
		Use:   "login",
		Short: "test credentials",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			if config.Username == "" && config.Password == "" {
				log.Println("Please set your username and password in the config file and reference it via path variable")
				return
			}

			log.Println("Testing login for " + config.Username)

			openId := rest.Login(&config)

			if openId != nil {
				rest.AuthTest(&config, openId)
			} else {
				log.Println("login failed.")
			}
		},
	}

	var createVendor = &cobra.Command{
		Use:   "create-vendor",
		Short: "creates a new vendor and adds you to it as a member",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			if config.Username == "" && config.Password == "" {
				log.Println("Please set your username and password in the config file and reference it via path variable")
				return
			}

			log.Println("Testing login for " + config.Username)

			openId := rest.Login(&config)

			if openId != nil {
				rest.CreateVendor(&config, vendor, openId)
			} else {
				log.Println("login failed.")
			}
		},
	}
	createVendor.Flags().StringVarP(&vendor, "vendor", "v", "", "The name of the vendor")

	var publishPackage = &cobra.Command{
		Use:   "publish-package",
		Short: "publish a package you own",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			if config.Username == "" && config.Password == "" {
				log.Println("Please set your username and password in the config file and reference it via path variable")
				return
			}
			openId := rest.Login(&config)

			if openId != nil {
				Publish(vendor, pack, version, accessLevel, &config, openId, force)
			} else {
				log.Println("login failed.")
			}
		},
	}
	publishPackage.Flags().StringVarP(&vendor, "vendor", "v", "", "The name of the vendor")
	publishPackage.Flags().StringVarP(&pack, "package", "p", "", "The name of the package")
	publishPackage.MarkFlagRequired("package")
	publishPackage.Flags().StringVarP(&version, "version", "s", "", "Version of the package")
	publishPackage.MarkFlagRequired("version")
	publishPackage.Flags().StringVarP(&accessLevel, "access-level", "a", "", "The desired access level. Either vendor or public")
	publishPackage.MarkFlagRequired("access-level")
	publishPackage.Flags().BoolVarP(&force, "force", "f", false, "Set this flag to skip all prompts")

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(loginTest)
	rootCmd.AddCommand(createVendor)
	rootCmd.AddCommand(publishPackage)

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
