package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OAuth struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURI  string `yaml:"redirect_uri"`
	} `yaml:"oauth"`
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate required fields
	if config.OAuth.ClientID == "" {
		return nil, fmt.Errorf("oauth.client_id is required")
	}
	if config.OAuth.ClientSecret == "" {
		return nil, fmt.Errorf("oauth.client_secret is required")
	}
	if config.OAuth.RedirectURI == "" {
		return nil, fmt.Errorf("oauth.redirect_uri is required")
	}

	// Set defaults for server config
	if config.Server.Port == 0 {
		config.Server.Port = 3000
	}
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	return &config, nil
}
