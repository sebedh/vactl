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
	"log"
	"os"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// policiesCmd represents the policies command
var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Vault policies",
	Long:  `Use this to get and apply policies from Vault`,
	Run: func(cmd *cobra.Command, args []string) {
		vaultToken := viper.GetString("vaultToken")
		vaultAddr := viper.GetString("vaultAddr")

		c, err := internal.NewVaultClient(vaultAddr, vaultToken)
		if err != nil {
			log.Printf("Could not create client needed: %v", err)
		}
		policies, err := getVaultPolicies(c)
		if err != nil {
			log.Printf("Command failed at getting policies: %v", err)
		}
		fmt.Println(policies)
	},
}

func getVaultPolicies(c *internal.Client) ([]internal.Policy, error) {

	var policies []internal.Policy

	policyList, err := c.VaultClient.Sys().ListPolicies()
	if err != nil {
		return nil,
			fmt.Errorf("Could not get list of policies from Vault: %v", err)
	}
	for _, p := range policyList {
		policy, err := internal.NewPolicy(p)
		if err != nil {
			return nil, fmt.Errorf("Could not append policy to policy list: %v", err)
		}
		if (policy.Name == "root") || (policy.Name == "default") {
			continue
		}
		policies = append(policies, *policy)
	}
	return policies, nil
}

func outputPolicyToFiles(p []internal.Policy, path string, c *internal.Client) error {
	for _, p := range p {
		fileName := path + p.Name + ".hcl"
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("Could not write file: %v\nERROR: %v", fileName, err)
		}
		defer f.Close()

		// Get the data
		data, err := c.VaultClient.Sys().GetPolicy(p.Name)
		if err != nil {
			return fmt.Errorf("Could not retrieve policy for output %v\nERROR: %v", p.Name, err)
		}

		if _, err := f.WriteString(data); err != nil {
			return fmt.Errorf("Could not write to file: %v", err)
		}
	}
	return nil
}

func init() {
	getCmd.AddCommand(policiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	policiesCmd.Flags().BoolP("out", "o", false, "output to hcl format in ./policies")
}
