package scratch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath" // Add this import
	"strings"

	"github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

func listAction(cCtx *cli.Context) error {
	interactive := cCtx.Bool("interactive")

	var scratchPath string
	// Check if the flag is set
	if cCtx.IsSet("scratch-path") {
		scratchPath = cCtx.String("scratch-path")
	} else {
		// Retrieve config from App Metadata if flag is not set
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
		scratchPath = cfg.Scratch.ScratchPath // Get the scratch path from config
	}

	// Check if scratch path is determined
	if scratchPath == "" {
		return fmt.Errorf("scratch path not configured. Please set 'scratch_path' in your config file or use the --scratch-path flag")
	}

	// Ensure scratchPath is absolute for consistent output
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	// Use the utility function to list valid scratch directories (full paths)
	directories, err := listScratchDirectories(absScratchPath)
	if err != nil {
		// listScratchDirectories handles checking if the base path exists and other read errors
		return fmt.Errorf("failed to list scratch directories: %w", err)
	}

	if interactive {
		// Interactive mode: Use fzf to select a directory

		// If no directories are found, exit early
		if len(directories) == 0 {
			fmt.Println("No matching scratch directories found.")
			return nil
		}

		// Generate the list of directory basenames for fzf
		var dirListBuilder strings.Builder
		for _, fullPath := range directories {
			dirListBuilder.WriteString(filepath.Base(fullPath)) // Use basename for fzf list
			dirListBuilder.WriteString("\n")
		}
		dirListString := dirListBuilder.String()

		// Prepare fzf command
		fzfCmd := exec.Command("fzf", "--tac", "--ansi", "--no-preview", "--prompt", "Select Scratch Directory> ")
		// No need to set fzfCmd.Dir as we provide absolute paths later
		fzfCmd.Stdin = strings.NewReader(dirListString) // Pipe the directory list to fzf
		fzfCmd.Stderr = os.Stderr                      // Inherit standard error

		// Capture the selected directory name (basename)
		selectedBaseName, err := cmd.CaptureCmd(*fzfCmd)
		if err != nil {
			// cmd.CaptureCmd returns error if fzf exits non-zero (e.g., Esc pressed)
			// or if the command fails. Check if output is empty which usually means no selection.
			if strings.TrimSpace(selectedBaseName) == "" {
				// User likely pressed Esc or Ctrl+C in fzf
				return fmt.Errorf("no directory selected")
			}
			return fmt.Errorf("failed to run fzf: %w", err)
		}

		// Trim whitespace and construct the full absolute path using the base scratch path
		trimmedSelectedBaseName := strings.TrimSpace(selectedBaseName)
		fullPath := filepath.Join(absScratchPath, trimmedSelectedBaseName)

		// Print the full absolute path
		fmt.Println(fullPath)
	} else {
		// Non-interactive mode: List directory basenames directly to stdout

		if len(directories) == 0 {
			// Print nothing if no directories found, consistent with `ls` behavior
			return nil
		}

		for _, fullPath := range directories {
			fmt.Println(filepath.Base(fullPath)) // Print only the directory basename
		}
	}

	return nil
}

// listCommand defines the 'list' subcommand for scratch.
var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List directories matching YYYY-MM-DD-<name> in the scratch path. Interactively select one if --interactive is set.", // Updated usage
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "scratch-path",
			Usage: "Override the scratch path from config",
		},
		&cli.BoolFlag{
			Name:    "interactive",
			Aliases: []string{"i"},
			Usage:   "Use interactive fuzzy finder (fzf) to select a directory",
			Value:   true,
		},
	},
	Action: listAction,
}
