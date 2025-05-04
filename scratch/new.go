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

	// Ensure scratchPath is absolute
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to get absolute path for scratch directory '%s': %v", scratchPath, err), 1)
	}

	// Check if the base scratch directory exists
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("scratch directory '%s' does not exist", absScratchPath), 1)
	} else if err != nil {
		// Handle other potential stat errors
		return cli.Exit(fmt.Sprintf("failed to check scratch directory '%s': %v", absScratchPath, err), 1)
	}

	// Check for existing directory with the same name suffix
	suffixToCheck := "-" + name
	entries, err := os.ReadDir(absScratchPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to read scratch directory '%s': %v", absScratchPath, err), 1)
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), suffixToCheck) {
			return cli.Exit(fmt.Sprintf("directory with name suffix '%s' already exists: %s", suffixToCheck, filepath.Join(absScratchPath, entry.Name())), 1)
		}
	}

	// Format the new directory name
	today := time.Now().Format("2006-01-02")
	newDirName := fmt.Sprintf("%s-%s", today, name)
	newDirPath := filepath.Join(absScratchPath, newDirName)

	// Create the new directory
	if err := os.Mkdir(newDirPath, 0755); err != nil {
		return cli.Exit(fmt.Sprintf("failed to create directory '%s': %v", newDirPath, err), 1)
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
