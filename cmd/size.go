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

	"code.cloudfoundry.org/bytefmt"
	"github.com/spf13/cobra"
	"hbtrove/pkg/api"
)

var sizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		_, d, err := api.LoadTroveData()
		if err != nil {
			panic(err)
		}

		var totalSize uint64 = 0
		for _, item := range d.Items {
			for _, download := range item.Downloads {
				totalSize += uint64(download.FileSize)
			}
		}
		fmt.Printf("Total size: %v\n", bytefmt.ByteSize(totalSize))
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)
}
