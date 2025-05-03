package scratch

import (
	"fmt"
	"os"

	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

func listAction(cCtx *cli.Context) error {
	// Retrieve config from App Metadata
	cfgInterface, ok := cCtx.App.Metadata["config"]
	if !ok {
		// This should ideally not happen if main.go sets it up correctly
		return fmt.Errorf("configuration not found in application metadata")
	}
	cfg, ok := cfgInterface.(config.Config)
	if !ok {
		// This indicates a programming error (wrong type stored)
		return fmt.Errorf("invalid configuration type in application metadata")
	}

	scratchCfg := cfg.Scratch // Get the scratch specific config

	// Check if scratch path is configured
	if scratchCfg.ScratchPath == "" {
		return fmt.Errorf("scratch path not configured. Please set 'scratch_path' in your config file")
	}

	// Check if the directory exists
	if _, err := os.Stat(scratchCfg.ScratchPath); os.IsNotExist(err) {
		return fmt.Errorf("scratch directory '%s' does not exist", scratchCfg.ScratchPath)
	}

	// List files using script package
	fmt.Printf("Listing files in %s:\n", color.CyanString(scratchCfg.ScratchPath))
	_, err := script.ListFiles(scratchCfg.ScratchPath).Stdout()
	if err != nil {
		return fmt.Errorf("failed to list files in scratch directory: %w", err)
	}

	return nil
}

// listCommand defines the 'list' subcommand for scratch.
var listCommand = &cli.Command{
	Name:   "list",
	Aliases: []string{"ls"},
	Usage:  "List files in the scratch directory",
	Action: listAction,
}
