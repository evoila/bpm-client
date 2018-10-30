package cmd

import (
	"fmt"
	"github.com/evoila/BPM-Client/helpers"
	"github.com/evoila/BPM-Client/model"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "BPM-Client",
	Short: "CLI Tool to access Bosh Package Manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello Hugo")
	},
}

func init() {

	var config = helpers.ReadConfig(defaultConfigLocation)
	var pack, version string

	var uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a package to Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {

			endpoint := config.Url + ":" + config.Port

			helpers.MoveToReleaseDir()

			Upload(endpoint, pack, config.Vendor, version)
		},
	}
	uploadCmd.Flags().StringVarP(&pack, "package", "p", "", "The name of the package to upload")
	uploadCmd.MarkFlagRequired("package")
	uploadCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to upload")
	uploadCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(uploadCmd)

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a package with all dependencies from Bosh Package Manager",
		Run: func(cmd *cobra.Command, args []string) {

			endpoint := config.Url + ":" + config.Port

			helpers.MoveToReleaseDir()

			requestBody := model.PackageRequestBody{
				Name:    pack,
				Vendor:  config.Vendor,
				Version: version}

			Download(endpoint, requestBody)
		},
	}
	downloadCmd.Flags().StringVarP(&pack, "package", "p", "", "The name of the package to upload")
	downloadCmd.MarkFlagRequired("package")
	downloadCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to upload")
	downloadCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(downloadCmd)

}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const defaultConfigLocation = "config.yml"
