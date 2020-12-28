package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
)

func GetList(c *api.Logical, path string) (s []string, err error) {
	r, err := c.List(path)

	if err != nil || r == nil {
		return nil, fmt.Errorf("Could not return list, wrong method or no items at path: %v\n", err)
	}

	data := r.Data["keys"].([]interface{})

	s = make([]string, len(data))

	for i, v := range data {
		s[i] = fmt.Sprint(v)
	}

	return s, nil
}

func ExportYaml(data interface{}) error {
	yaml, err := yaml.Marshal(data)

	if err != nil {
		return fmt.Errorf("Could not marshal to yaml: %v", err)
	}

	// print the yaml
	fmt.Println(string(yaml))
	return nil
}

func GetLocalFile(path string) (dir string, err error) {
	f, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("Could not determine path as file or dir: %v", err)
	}

	// dir, file
	dir, _ = filepath.Split(path)

	if f.IsDir() {
		err = filepath.Walk(dir, func(p string, info os.FileInfo, errr error) error {
			fileBytes, err := ioutil.ReadFile(p)
			if err != nil {
				return fmt.Errorf("Could not open file: %v", err)
			}

			if err = yaml.Unmarshal(fileBytes, make(map[interface{}]interface{})); err != nil {
				return fmt.Errorf("Could not unmarshal, %v", err)
			}
			// Determine if user or ssh-role
			return nil
		})
	}
	return "", err
}
