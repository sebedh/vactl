/*
Copyright © 2020 Sebastian Edholm <sebastian.edholm@iver.se>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0 Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FileApply string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:       "apply",
	Short:     "Used to apply policies, users and roles",
	Long:      `Used to apply policies, users and roles. Use with -f to apply files or folders`,
	ValidArgs: []string{"policies", "ssh-roles", "users"},
	Run:       applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	var filesToApply []string
	if len(FileApply) == 0 {
		fmt.Println("error: must specify -f <dir|file>")
		os.Exit(0)
	}
	f, err := os.Stat(FileApply)
	if err != nil {
		fmt.Printf("Could not stat dir or file: %v", err)
		os.Exit(1)
	}

	if f.IsDir() {
		// Do stuff when we walk
		err = filepath.Walk(FileApply, func(p string, info os.FileInfo, err error) error {
			extension := strings.ToLower(filepath.Ext(p))
			if !info.IsDir() {
				if extension == ".yml" || extension == ".yaml" || extension == ".hcl" {
					filesToApply = append(filesToApply, p)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		filesToApply = append(filesToApply, FileApply)
	}

	if err = applyFunc(filesToApply); err != nil {
		fmt.Printf("Could not apply files: %v", err)
	}
}

func applyFunc(files []string) error {
	vaultAddr := viper.GetString("vaultAddr")
	vaultToken := viper.GetString("vaultToken")
	if len(vaultAddr) <= 0 {
		return fmt.Errorf("Could not determine vaultAddress, please specify vaultAddr: in conf")
	}
	if len(vaultToken) <= 0 {
		return fmt.Errorf("You did not specify a Token in config, please specify vaultToken in conf")
	}

	// Create a client
	c, err := internal.NewVaultClient(vaultAddr, vaultToken)
	if err != nil {
		fmt.Printf("Could not establish Vault client: %v", err)
		return err
	}

	for _, f := range files {
		extension := strings.ToLower(filepath.Ext(f))
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("Could not read file: %v [%v]", f, err)
		}

		if extension == ".hcl" {
			if err = c.ApplyPolicyPath(f); err != nil {
				return err
			}
		} else {
			if err = c.ApplyDataPath(b, f); err != nil {
				return err
			}
		}

	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&FileApply, "file", "f", "", "File or folder to apply to Vault")
}
