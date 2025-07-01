package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type ScratchConfig struct {
	ScratchPath        string `toml:"scratch_path"`
	TmuxinatorTemplate string `toml:"tmuxinator_template"`
}

type Config struct {
	// Whether to automatically enable using cluster-admin role for non-read-only
	// commands that require it in test kubecontexts
	EnableClusterAdminForTest bool `toml:"enable_cluster_admin_for_test"`

	Scratch ScratchConfig `toml:"scratch"`
}

func NewConfig() Config {
	homeDir, err := os.UserHomeDir()
	defaultScratchPath := ""
	defaultTmuxinatorTemplate := ""

	if err == nil {
		defaultScratchPath = filepath.Join(homeDir, "Source", "scratch")
		defaultTmuxinatorTemplate = filepath.Join(homeDir, ".config", "mdcli", "scratch.yaml.template")
	}

	return Config{
		EnableClusterAdminForTest: true,
		Scratch: ScratchConfig{
			ScratchPath:        defaultScratchPath,
			TmuxinatorTemplate: defaultTmuxinatorTemplate,
		},
	}
}

func NewConfigFromToml(filePath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	config = NewConfig()

	if _, err := toml.Decode(string(data), &config); err != nil {
		return Config{}, err
	}

	if config.Scratch.TmuxinatorTemplate != "" && !filepath.IsAbs(config.Scratch.TmuxinatorTemplate) {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr == nil {
			// Handle tilde expansion explicitly
			if config.Scratch.TmuxinatorTemplate == "~" {
				config.Scratch.TmuxinatorTemplate = homeDir
			} else if strings.HasPrefix(config.Scratch.TmuxinatorTemplate, "~/") {
				config.Scratch.TmuxinatorTemplate = filepath.Join(homeDir, config.Scratch.TmuxinatorTemplate[2:])
			} else {
				absPath, absErr := filepath.Abs(config.Scratch.TmuxinatorTemplate)
				if absErr == nil {
					config.Scratch.TmuxinatorTemplate = absPath
				} else {
					// If absolute path fails, maybe log a warning but keep the original path
					fmt.Fprintf(os.Stderr, "Warning: could not make tmuxinator template path absolute: %v\n", absErr)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Warning: could not get user home directory to resolve tmuxinator template path: %v\n", homeErr)
		}
	}

	return config, nil
}

// LoadConfig loads configuration from a standard path or returns default.
func LoadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not get user home directory: %v\n", err)
		return NewConfig()
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
