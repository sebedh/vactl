package internal

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type User struct {
	Name     string   `yaml:"name"`
	Policies []string `yaml:"token_policies"`
	Method   string   `yaml:"method"`
}

func (u *User) ApplyToVault(c *Client) error {
	logical := c.VaultClient.Logical()
	path := "/auth/" + u.Method + "/users/" + u.Name

	data := make(map[string]interface{})
	if u.Method == "userpass" {
		data["password"] = GeneratePassword(10)
	}

	data["token_policies"] = u.Policies

	if _, err := logical.Write(path, data); err != nil {
		return fmt.Errorf("Could not install user: %v [%v]\n", u.Name, err)
	}

	return nil
}

type UserContainer struct {
	Type          string `yaml:"type"`
	UserContainer []User `yaml:"users"`
}

func NewUser(name string, policies []string, method string) (user *User, err error) {
	return &User{
		Name:     name,
		Policies: policies,
		Method:   method,
	}, nil
}

func (uc *UserContainer) AppendUser(user User) []User {
	uc.UserContainer = append(uc.UserContainer, user)
	return uc.UserContainer
}

func (uc *UserContainer) ImportYaml(yml []byte) error {
	if err := yaml.Unmarshal(yml, uc); err != nil {
		return fmt.Errorf("Could not unmarshal into object: %v\n", err)
	}
	return nil
}
