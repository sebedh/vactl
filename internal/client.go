package internal

import (
	"log"

	vault "github.com/hashicorp/vault/api"
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
