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
