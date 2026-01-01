package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adammmmm/go-junos"
)

type AuthConfig struct {
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	PrivateKey string `json:"privatekey,omitempty"`
}

func readAuthJson(path string) (*junos.AuthMethod, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening auth file: %w", err)
	}
	defer file.Close()

	var cfg AuthConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing auth file: %w", err)
	}

	if cfg.Username == "" {
		return nil, fmt.Errorf("auth.json: username is required")
	}

	switch {
	case cfg.Password != "":
		return &junos.AuthMethod{
			Credentials: []string{cfg.Username, cfg.Password},
		}, nil

	case cfg.PrivateKey != "":
		return &junos.AuthMethod{
			Username:   cfg.Username,
			PrivateKey: cfg.PrivateKey,
		}, nil

	default:
		return nil, fmt.Errorf("auth.json: either password or privatekey must be set")
	}
}
