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

	// Check if the directory exists
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return fmt.Errorf("scratch directory '%s' does not exist", absScratchPath)
	}

	if interactive {
		// Interactive mode: Use fzf to select a directory
		// Command to list only directories directly under scratchPath for fzf
		// -maxdepth 1: Don't go into subdirectories
		// -mindepth 1: Don't include '.' itself
		// -type d: Only list directories
		// -printf '%f\n': Print only the basename (directory name) followed by a newline
		listDirsCmd := `find . -maxdepth 1 -mindepth 1 -type d -printf '%f\n'`

		// Prepare fzf command
		fzfCmd := exec.Command("fzf", "--ansi", "--no-preview", "--prompt", "Select Scratch Directory> ")
		fzfCmd.Dir = absScratchPath // Run the find command inside the scratch directory
		fzfCmd.Stdin = os.Stdin     // Inherit standard input
		fzfCmd.Stderr = os.Stderr    // Inherit standard error

		// Set the FZF_DEFAULT_COMMAND environment variable
		fzfCmd.Env = append(os.Environ(),
			fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", listDirsCmd),
		)

		// Capture the selected directory name (relative path/basename)
		selectedRelativeDir, err := cmd.CaptureCmd(*fzfCmd)
		if err != nil {
			// cmd.CaptureCmd returns error if fzf exits non-zero (e.g., Esc pressed)
			// or if the command fails. Check if output is empty which usually means no selection.
			if strings.TrimSpace(selectedRelativeDir) == "" {
				return fmt.Errorf("no directory selected")
			}
			return fmt.Errorf("failed to run fzf: %w", err)
		}

		// Trim whitespace and construct the full absolute path
		trimmedSelectedDir := strings.TrimSpace(selectedRelativeDir)
		fullPath := filepath.Join(absScratchPath, trimmedSelectedDir)

		// Print the full absolute path
		fmt.Println(fullPath)
	} else {
		// Non-interactive mode: List directories directly to stdout
		listCmd := exec.Command("find", ".", "-maxdepth", "1", "-mindepth", "1", "-type", "d", "-printf", "%f\n")
		listCmd.Dir = absScratchPath
		listCmd.Stdout = os.Stdout // Pipe output directly to standard out
		listCmd.Stderr = os.Stderr // Pipe errors directly to standard error
		if err := listCmd.Run(); err != nil {
			return fmt.Errorf("failed to list directories in '%s': %w", absScratchPath, err)
		}
	}

	return nil
}

// listCommand defines the 'list' subcommand for scratch.
var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List directories in the scratch path. Interactively select one if --interactive is set.",
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
