package internal

import (
	"errors"
	"log"
)

type Policy struct {
	Name string
}

func NewPolicy(name string) (*Policy, error) {
	var err error
	if len(name) == 0 {
		err = errors.New("Was not given a policy name")
		log.Printf("Policy error: %v", err)
		return nil, err
	} else if name == "root" {
		return &Policy{Name: "root"}, nil
	}
	return &Policy{Name: name}, nil
}

// func GetLocalPolicies(path string) (localPolicies []Policy, dir string, err error) {
// 	f, err := os.Stat(path)
// 	if err != nil {
// 		fmt.Printf("Could not determine path as file or directory: %v", err)
// 		os.Exit(1)
// 	}
// 	dir, file := filepath.Split(path)
//
// 	// determine if path is file or dir
// 	// if not dir we should always target one file
// 	if f.IsDir() {
// 		err = filepath.Walk(dir, func(p string, info os.FileInfo, errr error) error {
// 			policyName := filepath.Base(strings.TrimSpace(strings.TrimSuffix(p, ".hcl")))
//
// 			policy, err := NewPolicy(strings.ToLower(policyName))
// 			if err != nil {
// 				fmt.Printf("Could not create policy object in code: %v", err)
// 				return err
// 			}
// 			localPolicies = append(localPolicies, *policy)
// 			return nil
// 		})
//
// 		if err != nil {
// 			fmt.Printf("Could not examine directory: %v", err)
// 			os.Exit(1)
// 		}
//
// 		// We don't want root dir name as a policy
// 		localPolicies = localPolicies[1:]
// 	} else {
// 		fileName := filepath.Base(strings.TrimSuffix(file, ".hcl"))
// 		policy, err := NewPolicy(strings.ToLower(strings.ToLower(fileName)))
// 		if err != nil {
// 			fmt.Printf("Could not make policy object from path: %v", err)
// 			os.Exit(1)
// 		}
// 		localPolicies = append(localPolicies, *policy)
// 	}
//
// 	return localPolicies, dir, nil
// }
