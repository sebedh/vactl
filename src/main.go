package main

import (
	"log"

	vault "github.com/hashicorp/vault/api"
)

var (
	vaultAddr  = "http://127.0.0.1:8200"
	vaultToken = "s.4ZSrKgLPuAqhZui4SxPVRfuy"
)

func main() {
	config := &vault.Config{Address: vaultAddr}
	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("Error could not create client: %v", err)
	}

	client.SetToken(vaultToken)
	if err != nil {
		log.Printf("Error on setting token: %v", err)
	}
}
