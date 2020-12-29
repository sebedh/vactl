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
	"strings"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes, user, ssh-role or policy object in Vault",
	Long: `Use delete to delete resources in Vault. You must specify what type you want to delete,
A user, ssh-role or policy. Ssh-role is simply typed as ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		if !(len(args) > 1) {
			fmt.Println("Need arguments, e.g, vactl delete policy <policy_name>")
			os.Exit(0)
		}
		run_args := strings.ToLower(args[0])
		vaultAddr := viper.GetString("vaultAddr")
		vaultToken := viper.GetString("vaultToken")
		client, err := internal.NewVaultClient(vaultAddr, vaultToken)
		if err != nil {
			fmt.Printf("Could not connect to Vault with Client: %v\n", err)
		}

		if run_args == "policy" {
			policy := strings.ToLower(args[1])
			if err := client.VaultClient.Sys().DeletePolicy(policy); err != nil {
				fmt.Printf("Could not delete policy: %v\n", err)
				os.Exit(1)
			}
		} else if run_args == "user" {
			user := strings.ToLower(args[1])
			path := "/auth/" + MethodUser + "/users/" + user
			if err := client.DeleteGivenPath(path); err != nil {
				fmt.Printf("Could not delete user: %v\n", err)
			}
		} else if run_args == "ssh" || run_args == "ssh-role" {
			sshRole := strings.ToLower(args[1])
			path := SshPath + "/roles/" + sshRole
			if err := client.DeleteGivenPath(path); err != nil {
				fmt.Printf("Could not delete ssh-role: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringVarP(&SshPath, "path", "p", "ssh", "Define a special path")
	deleteCmd.PersistentFlags().StringVarP(&MethodUser, "method", "m", "userpass", "The user login method & it's logical path in Vault, default is userpass")

}
