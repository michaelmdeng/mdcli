package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type ScratchConfig struct {
	ScratchPath string `toml:"scratch_path"`
}

type Config struct {
	// Whether to automatically enable using cluster-admin role for non-read-only
	// commands that require it in test kubecontexts
	EnableClusterAdminForTest bool `toml:"enable_cluster_admin_for_test"`

	Scratch ScratchConfig `toml:"scratch"`
}

func NewConfig() Config {
	return Config{
		EnableClusterAdminForTest: true,
		Scratch:                   ScratchConfig{}, // Initialize with zero values (empty path)
	}
}

func NewConfigFromToml(filePath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	// Set default values before decoding
	config = NewConfig()

	if _, err := toml.Decode(string(data), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
