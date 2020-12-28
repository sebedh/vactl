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
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
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
	var run_args string
	if len(FileApply) == 0 {
		fmt.Println("error: must specify -f <dir|file>")
		os.Exit(0)
	}

	if len(args[0]) > 0 {
		run_args = strings.ToLower(args[0])
	}

	if run_args == "policies" || run_args == "policy" {
		policiesToCommit, dir, err := internal.GetLocalPolicies(FileApply)

		if err != nil {
			fmt.Printf("Could not get local policies: %v", err)
		}

		if err := applyToVault(policiesToCommit, dir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if run_args == "user" || run_args == "users" {
		usersToCommit, dir, err := internal.GetLocalUsers(FileApply)

		if err != nil {
			fmt.Printf("Could not get local users yaml: %v", err)
		}

		if err := applyToVault(usersToCommit, dir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func applyToVault(o interface{}, dir string) error {
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

	// What should we apply?
	t := reflect.TypeOf(o)

	// It's a policy
	if t == reflect.TypeOf([]internal.Policy{}) {
		err := applyPolicy(c, o, dir)
		if err != nil {
			return fmt.Errorf("Could not apply Policy: %v", err)
		}
	}
	return nil
}

func applyPolicy(c *internal.Client, o interface{}, dir string) error {
	s := reflect.ValueOf(o)
	var reader io.Reader
	var buf bytes.Buffer

	for i := 0; i < s.Len(); i++ {
		policyName := s.Index(i).FieldByName("Name").String()
		path := dir + policyName + ".hcl"
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Could not open/find policy to install: %v", err)
		}

		defer file.Close()

		reader = file
		if _, err := io.Copy(&buf, reader); err != nil {
			return fmt.Errorf("Could not reat policy in buffer: %v", err)
		}

		policyName = strings.TrimSpace(strings.ToLower(policyName))
		fileBuf := buf.String()

		if err := c.VaultClient.Sys().PutPolicy(policyName, fileBuf); err != nil {
			fmt.Printf("Could not apply the policy to Vault: %v", err)
		}
		fmt.Printf("Applied Policy to Vault: %v", dir+policyName)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&FileApply, "file", "f", "", "File or folder to apply to Vault")
}
