package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings" // Add this import

	"github.com/BurntSushi/toml"
)

type ScratchConfig struct {
	ScratchPath        string `toml:"scratch_path"`
	TmuxinatorTemplate string `toml:"tmuxinator_template"` // Add this line
}

type Config struct {
	// Whether to automatically enable using cluster-admin role for non-read-only
	// commands that require it in test kubecontexts
	EnableClusterAdminForTest bool `toml:"enable_cluster_admin_for_test"`

	Scratch ScratchConfig `toml:"scratch"`
}

func NewConfig() Config {
	homeDir, err := os.UserHomeDir()
	defaultScratchPath := ""         // Default to empty if home dir cannot be found
	defaultTmuxinatorTemplate := "" // Default to empty if home dir cannot be found

	if err == nil {
		defaultScratchPath = filepath.Join(homeDir, "Source", "scratch")
		// Set default template path relative to home directory
		defaultTmuxinatorTemplate = filepath.Join(homeDir, ".config", "mdcli", "scratch.yaml.template") // Add this line
	}

	return Config{
		EnableClusterAdminForTest: true,
		Scratch: ScratchConfig{
			ScratchPath:        defaultScratchPath,
			TmuxinatorTemplate: defaultTmuxinatorTemplate, // Add this line
		},
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

	// Ensure TmuxinatorTemplate path is absolute if it was loaded from config
	// If it's still the default, it's already absolute or empty.
	// If it was loaded, it might be relative to the config file or user's CWD.
	// Best practice is often to resolve relative paths based on the config file's location,
	// but for simplicity here, we'll resolve based on CWD if not absolute.
	// However, since the default is already absolute, we only need to handle
	// the case where the user provided a non-absolute path in the config file.
	// Let's make it absolute based on the user's home directory if it's not already absolute.
	if config.Scratch.TmuxinatorTemplate != "" && !filepath.IsAbs(config.Scratch.TmuxinatorTemplate) { // Add this block
		homeDir, homeErr := os.UserHomeDir()
		if homeErr == nil {
			// Handle tilde expansion explicitly
			if config.Scratch.TmuxinatorTemplate == "~" {
				config.Scratch.TmuxinatorTemplate = homeDir
			} else if strings.HasPrefix(config.Scratch.TmuxinatorTemplate, "~/") {
				config.Scratch.TmuxinatorTemplate = filepath.Join(homeDir, config.Scratch.TmuxinatorTemplate[2:])
			} else {
				// If it's not using tilde, assume it's relative to home dir for consistency
				// Or alternatively, make it relative to the config file directory.
				// Let's stick to making it absolute based on home dir for now.
				// config.Scratch.TmuxinatorTemplate = filepath.Join(homeDir, config.Scratch.TmuxinatorTemplate)
				// Correction: Let's make it absolute based on the current working directory if not absolute and not using tilde.
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
	} // End added block

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
