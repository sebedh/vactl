package internal

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type SshRole struct {
	Name               string   `yaml:"name"`
	Key_type           string   `yaml:"key_type"`
	Default_user       string   `yaml:"default_user"`
	Cidr_list          []string `yaml:"cidr_list"`
	Allowed_users      []string `yaml:"allowed_users"`
	Port               int      `yaml:"port"`
	Excluded_cidr_list []string `yaml:"excluded_cidr_list"`
}

func (s *SshRole) ApplyToVault(c *Client, path string) error {
	logical := c.VaultClient.Logical()
	path = path + "/roles/" + s.Name

	data := make(map[string]interface{})
	data["key_type"] = s.Key_type
	data["default_user"] = s.Default_user
	data["allowed_users"] = strings.Join(s.Allowed_users, ",")
	data["cidr_list"] = strings.Join(s.Cidr_list, ",")
	data["excluded_cidr_list"] = strings.Join(s.Excluded_cidr_list, ",")
	data["port"] = s.Port

	if _, err := logical.Write(path, data); err != nil {
		return fmt.Errorf("Could not write role")
	}

	return nil
}

type SshRoleContainer struct {
	Type             string    `yaml:"type"`
	Path             string    `yaml:"path"`
	SshRoleContainer []SshRole `yaml:"sshroles"`
}

func NewSshContainerFromYaml(b []byte) (r *SshRoleContainer, err error) {
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, fmt.Errorf("Could not unmarshal into object: %v\n", err)
	}
	return
}

func (r *SshRoleContainer) AppendSshRole(role SshRole) []SshRole {
	r.SshRoleContainer = append(r.SshRoleContainer, role)
	return r.SshRoleContainer
}

// func (r *SshRoleContainer) ImportYaml(yml []byte) error {
// 	if err := yaml.Unmarshal(yml, r); err != nil {
// 		return fmt.Errorf("Could not unmarshal into object: %v\n", err)
// 	}
// 	return nil
// }
