package internal

import (
	"bytes"
	"encoding/base64"
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
		nil
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

func (c *Client) Write(path string, data map[string]string, ver string) error {
	logical := c.VaultClient.Logical()
	body := make(map[string]interface{})

	// Decode the base64 values
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		body[k] = string(b)
	}

	var err error

	if ver == "2" {
		splitted := strings.Split(path, "/")
		splitted = append(splitted, "")
		copy(splitted[2:], splitted[1:])
		splitted[1] = "data"
		newPath := strings.Join(splitted, "/")

		d := make(map[string]interface{})
		d["data"] = body
		_, err = logical.Write(newPath, d)
	} else {
		_, err = logical.Write(path, body)
	}

	return err
}

// This is too read secrets
func (c *Client) Read(path string) map[string]interface{} {
	out := make(map[string]interface{})

	s, err := c.VaultClient.Logical().Read(path)
	if err != nil {
		fmt.Printf("Error reading secrets, err=%v", err)
		return nil
	}

	// Encode all k,v pairs
	if s == nil || s.Data == nil {
		fmt.Printf("No data to read at path, %s\n", path)
		return out
	}
	for k, v := range s.Data {
		switch t := v.(type) {
		case string:
			out[k] = base64.StdEncoding.EncodeToString([]byte(t))
		case map[string]interface{}:
			if k == "data" {
				for x, y := range t {
					if z, ok := y.(string); ok {
						out[x] = base64.StdEncoding.EncodeToString([]byte(z))
					}
				}
			}
		default:
			fmt.Printf("error reading value at %s, key=%s, type=%T\n", path, k, v)
		}
	}

	return out
}

// List the keys at at given vault path. This has only been tested on the generic backend.
// It will return nil if something goes wrong.
func (c *Client) List(path string) []string {
	secret, err := c.VaultClient.Logical().List(path)
	if secret == nil {
		return nil
	}
	if err != nil {
		fmt.Printf("Unable to read path %q, err=%v\n", path, err)
		return nil
	}

	r, ok := secret.Data["keys"].([]interface{})
	if ok {
		out := make([]string, len(r))
		for i := range r {
			out[i] = r[i].(string)
		}
		return out
	}
	return nil
}
