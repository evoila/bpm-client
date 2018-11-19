package cmd

import (
	"fmt"
	"github.com/evoila/BPM-Client/helpers"
	. "github.com/evoila/BPM-Client/model"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "BPM-Client",
	Short: "CLI Tool to access Bosh Package Manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one command of: upload, update, download, delete, search")
	},
}

var config Config
var pack, version, vendor string
var update bool

func init() {

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a package to Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			setupConfig()
			log.Println("Begin upload.")

			if update {
				RunUpdateIfPackagePresentUploadIfNot(pack, &config)
			} else {
				CheckIfAlreadyPresentAndUpload(pack, &config)
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

			requestBody := PackageRequestBody{
				Name:    pack,
				Vendor:  vendor,
				Version: version}

			Download("", requestBody, &config)
		},
	}
	downloadCmd.Flags().StringVarP(&pack, "package", "p", "", "The name of the package")
	downloadCmd.MarkFlagRequired("package")
	downloadCmd.Flags().StringVarP(&vendor, "vendor", "v", "", "The name of the vendor")
	downloadCmd.MarkFlagRequired("package")
	downloadCmd.Flags().StringVarP(&version, "version", "s", "", "Version of the package")
	downloadCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
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
