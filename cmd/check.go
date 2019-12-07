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
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
	"hbtrove/pkg/api"
	"hbtrove/pkg/checker"
)

var checkCmd = &cobra.Command{
	Use: "check",
	Run: func(cmd *cobra.Command, args []string) {
		if directory == "" {
			panic("No directory is set")
		}

		err := check(directory, checkContents)
		if err != nil {
			panic(err)
		}
	},
}

var directory string
var checkContents bool

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&directory, "directory", "d", "", "")
	checkCmd.Flags().BoolVarP(&checkContents, "check-contents", "c", false, "")
}

func check(dir string, checkContents bool) error {
	fmt.Printf("Checking data in dir: %v\n", directory)

	jsons, td, err := api.LoadTroveData()
	if err != nil {
		return err
	}
	err = saveJsons(jsons, dir)
	if err != nil {
		return err
	}

	results := checker.Check(td, dir, checkContents)
	summary := map[checker.DownloadStatus]int{}
	for _, result := range results {
		if result.Status != checker.Same {
			fmt.Printf("%v %v %v %v\n", result.Status, result.Platform, result.Method, result.Product.HumanName)
		}
		summary[result.Status] += 1
	}

	println()
	println("Summary:")
	for status, count := range summary {
		println(status, count)
	}
	return nil
}

func saveJsons(jsons [][]byte, dir string) error {
	for i, json := range jsons {
		err := saveJson(dir, i, json)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveJson(dir string, i int, json []byte) error {
	f, err := os.Create(path.Join(dir, fmt.Sprintf("data-chunk-%03d.json", i)))
	if err != nil {
		return nil
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(json))
	return err
}
