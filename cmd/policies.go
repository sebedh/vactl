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
	"path/filepath"
	"strings"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
)

// policiesCmd represents the policies command
var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Vault policies",
	Long:  `Represent policies in Vault`,
	Run: func(cmd *cobra.Command, args []string) {
		localPolicies := getLocalPolicies(FileApply)
		fmt.Println(*localPolicies)
	},
}

func getLocalPolicies(path string) *[]internal.Policy {
	var localPolicies []internal.Policy
	f, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Could not determine path as file or directory: %v", err)
		os.Exit(1)
	}

	// determine if path is file or dir
	// we if not dir we should always target one file
	if f.IsDir() {
		err = filepath.Walk(path, func(p string, info os.FileInfo, errr error) error {
			policyName := filepath.Base(strings.TrimSpace(strings.TrimSuffix(p, ".hcl")))

			policy, err := internal.NewPolicy(strings.ToLower(policyName))
			if err != nil {
				fmt.Printf("Could not create policy object in code: %v", err)
				return err
			}
			localPolicies = append(localPolicies, *policy)
			return nil
		})

		if err != nil {
			fmt.Printf("Could not examine directory: %v", err)
			os.Exit(1)
		}

		// We don't want root dir name as a policy
		localPolicies = localPolicies[1:]
	} else {
		fileName := filepath.Base(strings.TrimSuffix(path, ".hcl"))
		policy, err := internal.NewPolicy(strings.ToLower(strings.ToLower(fileName)))
		if err != nil {
			fmt.Printf("Could not make policy object from path: %v", err)
			os.Exit(1)
		}
		localPolicies = append(localPolicies, *policy)
	}

	return &localPolicies
}

func init() {
	getCmd.AddCommand(policiesCmd)
	applyCmd.AddCommand(policiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
