package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	homeDir, err := os.UserHomeDir()
	defaultScratchPath := "" // Default to empty if home dir cannot be found
	if err == nil {
		defaultScratchPath = filepath.Join(homeDir, "Source", "scratch")
	}

	return Config{
		EnableClusterAdminForTest: true,
		Scratch:                   ScratchConfig{ScratchPath: defaultScratchPath},
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

// LoadConfig loads configuration from a standard path or returns default.
func LoadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not get user home directory: %v\n", err)
		return NewConfig() // Return default if home dir fails
	}
	configPath := filepath.Join(home, ".config", "mdcli", "config.toml")

	cfg, err := NewConfigFromToml(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			// Only warn if it's not a "file not found" error
			fmt.Fprintf(os.Stderr, "Warning: failed to load config from %s: %v\n", configPath, err)
		}
		// Return default config if file doesn't exist or fails to parse
		return NewConfig()
	}
	return cfg
}
