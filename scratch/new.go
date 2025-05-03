package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

// newAction implements the logic for the 'scratch new' command.
func newAction(cCtx *cli.Context) error {
	if cCtx.NArg() != 1 {
		return fmt.Errorf("exactly one argument <name> must be provided")
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
			return fmt.Errorf("configuration not found in application metadata")
		}
		cfg, ok := cfgInterface.(config.Config)
		if !ok {
			return fmt.Errorf("invalid configuration type in application metadata")
		}
		scratchPath = cfg.Scratch.ScratchPath
	}

	// Check if scratch path is determined
	if scratchPath == "" {
		return fmt.Errorf("scratch path not configured. Please set 'scratch_path' in your config file or use the --scratch-path flag")
	}

	// Ensure scratchPath is absolute
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	// Check if the base scratch directory exists
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return fmt.Errorf("scratch directory '%s' does not exist", absScratchPath)
	}

	// Check for existing directory with the same name suffix
	suffixToCheck := "-" + name
	entries, err := os.ReadDir(absScratchPath)
	if err != nil {
		return fmt.Errorf("failed to read scratch directory '%s': %w", absScratchPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), suffixToCheck) {
			return fmt.Errorf("directory with name suffix '%s' already exists: %s", suffixToCheck, filepath.Join(absScratchPath, entry.Name()))
		}
	}

	// Format the new directory name
	today := time.Now().Format("2006-01-02")
	newDirName := fmt.Sprintf("%s-%s", today, name)
	newDirPath := filepath.Join(absScratchPath, newDirName)

	// Create the new directory
	if err := os.Mkdir(newDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", newDirPath, err)
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
	},
	Action: newAction,
}
