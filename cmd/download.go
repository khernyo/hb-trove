/*
Copyright © 2019 Szabolcs Berecz <szabolcs.berecz@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hbtrove/pkg/api"
	"hbtrove/pkg/downloader"
)

var downloadCmd = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		if directory == "" {
			panic("No directory is set")
		}

		config, err := loadConfig()
		if err != nil {
			panic(err)
		}

		err = download(config, directory, checkContents, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

var dryRun bool

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&directory, "directory", "d", "", "")
	downloadCmd.Flags().BoolVarP(&checkContents, "check-contents", "c", false, "")
	downloadCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "")
}

func download(config *downloader.Config, dir string, checkContents bool, dryRun bool) error {
	fmt.Printf("Downloading into dir: %v\n", directory)

	jsons, td, err := api.LoadTroveData()
	if err != nil {
		return err
	}
	err = saveJsons(jsons, dir)
	if err != nil {
		return err
	}

	err = downloader.Download(config, td, dir, checkContents, dryRun)
	return err
}

func loadConfig() (*downloader.Config, error) {
	return downloader.NewConfigFromFile("config.toml")
}
