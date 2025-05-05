package scratch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

// newAction implements the logic for the 'scratch new' command.
func newAction(cCtx *cli.Context) error {
	if cCtx.NArg() != 1 {
		return cli.Exit("exactly one argument <name> must be provided", 2)
	}
	name := cCtx.Args().Get(0)

	var scratchPath string
	// Check if the flag is set
	if cCtx.IsSet("scratch-path") {
		scratchPath = cCtx.String("scratch-path")
	} else {
		// Retrieve config from App Metadata if flag is not set
		cfgInterface, ok := cCtx.App.Metadata["config"]
		if !ok {
			return cli.Exit("configuration not found in application metadata", 1)
		}
		cfg, ok := cfgInterface.(config.Config)
		if !ok {
			return cli.Exit("invalid configuration type in application metadata", 1)
		}
		scratchPath = cfg.Scratch.ScratchPath
	}

	// Check if scratch path is determined
	if scratchPath == "" {
		return cli.Exit("scratch path not configured. Please set 'scratch_path' in your config file or use the --scratch-path flag", 1)
	}

	// Ensure scratchPath is absolute and expand ~
	absScratchPath, err := expandPath(scratchPath) // Use the new helper
	if err != nil {
		// expandPath provides a formatted error
		return cli.Exit(err.Error(), 1)
	}

	// Check if the base scratch directory exists
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("scratch directory '%s' does not exist", absScratchPath), 1)
	} else if err != nil {
		// Handle other potential stat errors
		return cli.Exit(fmt.Sprintf("failed to check scratch directory '%s': %v", absScratchPath, err), 1)
	}

	// Check for existing directory matching the name (exact or suffix)
	existingPath, err := findScratchDirectory(absScratchPath, name)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error checking for existing directory '%s': %v", name, err), 1)
	}
	if existingPath != "" {
		// A directory matching the name (exact or suffix) already exists.
		// findScratchDirectory returns the full path.
		return cli.Exit(fmt.Sprintf("directory matching name '%s' already exists: %s", name, existingPath), 1)
	}

	// Create the new directory using the utility function
	createReadme := cCtx.Bool("create-readme")
	newDirPath, err := createScratchDirectory(absScratchPath, name, createReadme)
	if err != nil {
		// createScratchDirectory already provides a descriptive error
		return cli.Exit(err.Error(), 1)
	}

	// Print the full path of the created directory
	fmt.Println(newDirPath)

	return nil
}

// newCommand defines the 'new' subcommand for scratch.
var newCommand = &cli.Command{
	Name:      "new",
	Usage:     "Create a new dated directory in the scratch path",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "scratch-path",
			Usage: "Override the scratch path from config",
		},
		&cli.BoolFlag{
			Name:    "create-readme",
			Aliases: []string{"r"},
			Usage:   "Create an empty README.md in the new directory",
			Value:   true, // Default to true
		},
	},
	Action: newAction,
}
