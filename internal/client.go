package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
)

type Client struct {
	vaultAddr   string
	vaultToken  string
	VaultClient *vault.Client
}

func NewVaultClient(vaultAddr string, vaultToken string) (*Client, error) {

	config := &vault.Config{Address: vaultAddr}
	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("Error could not create client: %v", err)
	}

	client.SetToken(vaultToken)
	if err != nil {
		log.Printf("Error on setting token: %v", err)
	}

	return &Client{
			vaultAddr:   vaultAddr,
			vaultToken:  vaultToken,
			VaultClient: client},
		err
}

func (c *Client) ApplyDataPath(b []byte, f string) error {
	content := make(map[interface{}]interface{})
	fmt.Printf("Should apply yaml: %v\n", f)

	if err := yaml.Unmarshal(b, content); err != nil {
		return fmt.Errorf("Could not unmarshal: %v\n", err)
	}

	if content["type"] == "users" {
		uc, err := NewUserContainerFromYaml(b)

		if err != nil {
			return fmt.Errorf("Failed creating container: %v\n", err)
		}

		for _, u := range uc.UserContainer {
			if err := u.ApplyToVault(c); err != nil {
				return fmt.Errorf("Could not apply to Vault: %v\n", err)
			}
		}

	} else if content["type"] == "sshrole" {
		rc, err := NewSshContainerFromYaml(b)
		if err != nil {
			return fmt.Errorf("Failed creating container: %v\n", err)
		}
		path := content["path"].(string)

		for _, r := range rc.SshRoleContainer {
			if err := r.ApplyToVault(c, path); err != nil {
				return fmt.Errorf("Could not apply to Vault: %v\n", err)
			}
		}
	}
	return nil

}

func (c *Client) ApplyPolicyPath(path string) error {
	var reader io.Reader
	var buf bytes.Buffer
	_, fName := filepath.Split(path)

	policyName := strings.TrimSuffix(fName, ".hcl")

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open/find policy to install: %v", err)
	}

	defer file.Close()

	reader = file
	if _, err := io.Copy(&buf, reader); err != nil {
		return fmt.Errorf("Could not read policy in buffer: %v", err)
	}

	policyName = strings.TrimSpace(strings.ToLower(policyName))
	fileBuf := buf.String()

	if err := c.VaultClient.Sys().PutPolicy(policyName, fileBuf); err != nil {
		fmt.Printf("Could not apply the policy to Vault: %v", err)
	}

	fmt.Printf("Applied Policy to Vault: %v, Location: %v\n", policyName, path)

	return nil
}

func (c *Client) DeleteGivenPath(path string) error {
	logical := c.VaultClient.Logical()

	if _, err := logical.Delete(path); err != nil {
		return fmt.Errorf("Could not delete given path: %v [%v]", path, err)
	}

	return nil
}
