package main

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Command     []string                     `yamle:"command"`
	Host        string                       `yaml:"host"`
	Token       string                       `yaml:"token"`
	RoleID      string                       `yaml:"roleId"`
	SecretId    string                       `yaml:"secretId"`
	SecretIdEnv string                       `yaml:"secretIdEnv"`
	SecretMount string                       `yaml:"secretMount"`
	SecretPaths map[string]map[string]string `yaml:"secrets"`
}

func LoadConfig(f io.Reader) (c Config, err error) {
	d := yaml.NewDecoder(f)
	config := Config{}

	err = d.Decode(&config)
	if err != nil {
		return config, err
	}

	if config.SecretId == "" {
		v, ok := os.LookupEnv(config.SecretIdEnv)
		if !ok {
			v = os.Getenv("VAULT_SECRET_ID")
		}
		config.SecretId = v
	}

	return config, nil
}
