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

	td, err := api.LoadTroveData()
	if err != nil {
		return err
	}

	results := checker.Check(td, dir, checkContents)
	for _, r1 := range results {
		for _, r2 := range r1.Results {
			for _, r3 := range r2.Results {
				if r3.Status != checker.Same {
					fmt.Printf("%v %v %v %v\n", r3.Status, r2.Platform, r3.Method, r1.Product.HumanName)
				}
			}
		}

	}
	return nil
}
