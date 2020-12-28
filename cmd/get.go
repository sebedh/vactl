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
	"strings"

	"github.com/sebedh/vactl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:       "get",
	Short:     "A brief description of your command",
	Long:      `Get a Vault resource to stdout or yaml`,
	ValidArgs: []string{"policies", "ssh-roles", "users"},
	Run:       getRun,
}

var Out bool
var MethodUser string
var SshPath string

func getRun(cmd *cobra.Command, args []string) {
	var run_arg string
	vaultAddr := viper.GetString("vaultAddr")
	vaultToken := viper.GetString("vaultToken")

	// Create a Client
	c, err := internal.NewVaultClient(vaultAddr, vaultToken)
	if err != nil {
		fmt.Printf("Could not create Vault client: %v", err)
	}

	if len(args) > 0 {
		run_arg = strings.ToLower(args[0])
	}

	if run_arg == "policies" || run_arg == "policy" {
		if err := getPolicies(args, c); err != nil {
			fmt.Printf("Failed at getting policies: %v\n", err)
			return
		}
	} else if run_arg == "users" || run_arg == "user" {
		if err := getUsers(c); err != nil {
			fmt.Printf("Failed at getting users: %v\n", err)
			return
		}
	} else if run_arg == "ssh-roles" || run_arg == "ssh" {
		if err := getSshRoles(args, c); err != nil {
			fmt.Printf("Could not get ssh-roles: %v\n", err)
		}
	} else {
		fmt.Println("Need to give, policies, ssh-roles or users as argument.")
	}
}

func getPolicies(args []string, c *internal.Client) error {
	if len(args) > 1 {
		policy, err := c.VaultClient.Sys().GetPolicy(args[1])
		if err != nil {
			return fmt.Errorf("Could not get policy: %v\n", err)
		}
		if Out {
			fmt.Println(policy)
			return nil
		}
		if len(policy) > 0 {
			fmt.Println(args[1])
			return nil
		}
		fmt.Println("Could not find policy")
	} else {
		policiesList, err := c.VaultClient.Sys().ListPolicies()
		if err != nil {
			return fmt.Errorf("Could list policies: %v", err)
		}

		// Check if non-default policy is installed
		if len(policiesList) > 2 {
			for _, p := range policiesList {
				// We don't want to print default and root policy
				if p == "default" || p == "root" {
					continue
				} else if len(p) > 0 {
					if Out {
						policyData, err := c.VaultClient.Sys().GetPolicy(p)
						if err != nil {
							fmt.Printf("Could not get data of policy: %v, %v ", p, err)
						}
						fmt.Println(policyData)
					} else {
						fmt.Println(p)
					}
				}
			}
		} else {
			fmt.Println("No non-default policies installed")
		}
	}
	return nil
}

func getUsers(c *internal.Client) error {
	logical := c.VaultClient.Logical()
	path := "/auth/" + MethodUser + "/users"

	users, err := internal.GetList(logical, path)
	if err != nil {
		return fmt.Errorf("Could not get list of users: %v", err)
	}

	for _, u := range users {
		if !Out {
			fmt.Println(u)
		}
	}

	// We want to print
	if !Out {
		fmt.Printf("Getting users at path: %v", path)
		for _, u := range users {
			fmt.Println(u)
		}
	} else if Out {
		// Create user container
		userContainer := internal.UserContainer{Type: "users"}

		// Iterate through users
		for _, u := range users {
			// Read the user data from Vault
			userPath := path + "/" + u
			data, err := logical.Read(userPath)
			if err != nil {
				return fmt.Errorf("Cannot output: %v", err)
			}

			// Get it into yaml form
			content, err := yaml.Marshal(data.Data)

			if err != nil {
				return fmt.Errorf("Cannot Marshal into Yaml: %v", err)
			}

			// Prepare user
			user := internal.User{Name: u, Method: MethodUser}

			// Put it back to our struct
			if err := yaml.Unmarshal(content, &user); err != nil {
				return fmt.Errorf("Cannot Unmarhsal object: %v\n", err)
			}

			// Append the struct to container struct
			userContainer.AppendUser(user)
		}
		// Print and export
		if err := internal.ExportYaml(userContainer); err != nil {
			return fmt.Errorf("Cannot export yaml content: %v\n", err)
		}
	}
	return nil
}

func getSshRoles(args []string, c *internal.Client) error {
	logical := c.VaultClient.Logical()
	path := "/" + SshPath + "/roles"

	roles, err := internal.GetList(logical, path)
	if err != nil {
		return fmt.Errorf("Error getting list: %v", err)
	}

	if !Out {
		fmt.Printf("Getting ssh roles at path: %v\n", path)
		for _, r := range roles {
			fmt.Println(r)
		}
	} else if Out {
		rContainer := internal.SshRoleContainer{
			Type: "sshrole",
			Path: path,
		}
		for _, role := range roles {
			rPath := path + "/" + role

			data, err := logical.Read(rPath)

			if err != nil {
				return fmt.Errorf("Could not read the specific role: %v\nBecouse: %v", role, err)
			}

			content, err := yaml.Marshal(data.Data)
			if err != nil {
				return fmt.Errorf("Could not prase to yaml: %v", err)
			}
			m := make(map[interface{}]interface{})

			if err := yaml.Unmarshal(content, &m); err != nil {
				return fmt.Errorf("Could not unmarshal into map, %v", err)
			}

			var excluded_cidr_list []string
			var allowed_users []string
			var cidr_list []string
			var port int

			if m["excluded_cidr_list"] != nil {
				excluded_cidr_list = strings.Split(m["excluded_cidr_list"].(string), ",")
			}
			if m["allowed_users"] != nil {
				allowed_users = strings.Split(m["allowed_users"].(string), ",")
			}
			if m["cidr_list"] != nil {
				cidr_list = strings.Split(m["cidr_list"].(string), ",")
			}
			if m["port"] != nil {
				port = m["port"].(int)
			}
			_ = excluded_cidr_list

			r := internal.SshRole{
				Name:               role,
				Key_type:           m["key_type"].(string),
				Allowed_users:      allowed_users,
				Default_user:       m["default_user"].(string),
				Cidr_list:          cidr_list,
				Excluded_cidr_list: excluded_cidr_list,
				Port:               port,
			}

			rContainer.AppendSshRole(r)
		}
		// Print and export
		if err := internal.ExportYaml(rContainer); err != nil {
			return fmt.Errorf("Cannot export yaml content: %v\n", err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.PersistentFlags().BoolVarP(&Out, "out", "o", false, "output to hcl format in ./policies")
	getCmd.PersistentFlags().StringVarP(&MethodUser, "method", "m", "userpass", "The user login method & it's logical path in Vault, default is userpass")
	getCmd.PersistentFlags().StringVarP(&SshPath, "path", "p", "ssh", "Define a special path")
}
