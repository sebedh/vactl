/*
Copyright © 2020 Sebastian Edholm <sebastian.edholm@iver.se>

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
	"os"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
)

var ImportFile string
var ImportPath string
var ImportVer string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: importRun,
}

func importRun(cmd *cobra.Command, args []string) {
	if len(ImportPath) == 0 || len(ImportFile) == 0 || len(ImportVer) == 0 {
		fmt.Println("error: Must speficy, -f -v and -p")
		os.Exit(127)
	}
	if err := internal.Import(ImportPath, ImportFile, ImportVer); err != nil {
		fmt.Println("error: oh no")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().StringVarP(&ImportFile, "file", "f", "", "json file to import")
	importCmd.PersistentFlags().StringVarP(&ImportPath, "importPath", "p", "", "vault import path")
	importCmd.PersistentFlags().StringVarP(&ImportVer, "importVer", "v", "", "vault import version")
}
