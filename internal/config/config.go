package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Whether to automatically enable using cluster-admin role for non-read-only
	// commands that require it in test kubecontexts
	EnableClusterAdminForTest bool
}

func NewConfig() Config {
	return Config{
		EnableClusterAdminForTest: true,
	}
}

func NewConfigFromToml(filePath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(data), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
