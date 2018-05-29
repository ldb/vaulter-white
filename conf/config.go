package conf

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Command     []string                     `yamle:"command"`
	Host        string                       `yaml:"host"`
	HostEnv     string                       `yaml:"hostEnv"`
	Token       string                       `yaml:"token"`
	RoleID      string                       `yaml:"roleId"`
	RoleIDEnv   string                       `yaml:"roleIdEnv"`
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

	if config.Host == "" {
		config.Host = safeLookupEnv(config.HostEnv, "VAULT_HOST")
	}

	if config.RoleID == "" {
		config.RoleID = safeLookupEnv(config.RoleIDEnv, "VAULT_ROLE_ID")
	}

	if config.SecretId == "" {
		config.SecretId = safeLookupEnv(config.SecretIdEnv, "VAULT_SECRET_ID")
	}

	return config, nil
}

func safeLookupEnv(env string, fallbackEnv string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		v = os.Getenv(fallbackEnv)
	}
	return v
}
