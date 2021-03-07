package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Json formatting for export and import
type Wrap struct {
	Data []Item `json:"data"`
}
type Item struct {
	Path  string `json:"path"`
	Pairs []Pair `json:"pairs"`
}
type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Import(path, file, ver string) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	// Check the input file exists
	if _, err := os.Stat(abs); err != nil {
		f, err := os.Create(abs)
		defer f.Close()
		if err != nil {
			return err
		}
	}

	// Read input file
	b, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	// Parse data
	var wrap Wrap
	err = json.Unmarshal(b, &wrap)
	if err != nil {
		return err
	}
	vaultAddr := viper.GetString("vaultAddr")
	vaultToken := viper.GetString("vaultToken")

	// Setup vault client
	v, err := NewVaultClient(vaultAddr, vaultToken)
	if v == nil || err != nil {
		if err != nil {
			return err
		}
		return errors.New("Unable to create vault client")
	}

	// Write each keypair to vault
	for _, item := range wrap.Data {
		data := make(map[string]string)
		for _, kv := range item.Pairs {
			data[kv.Key] = kv.Value
		}
		fmt.Printf("Writing %s\n", item.Path)
		if err := v.Write(item.Path, data, ver); err != nil {
			fmt.Printf("here %s\n", err)
		}
	}

	return nil
}
