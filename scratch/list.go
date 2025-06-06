package scratch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

func listAction(cCtx *cli.Context) error {
	interactive := cCtx.Bool("interactive")

	var scratchPath string
	if cCtx.IsSet("scratch-path") {
		scratchPath = cCtx.String("scratch-path")
	} else {
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

	if scratchPath == "" {
		return fmt.Errorf("scratch path not configured. Please set 'scratch_path' in your config file or use the --scratch-path flag")
	}

	absScratchPath, err := expandPath(scratchPath)
	if err != nil {
		return fmt.Errorf("failed to resolve scratch path: %w", err)
	}

	directories, err := listScratchDirectories(absScratchPath)
	if err != nil {
		return fmt.Errorf("failed to list scratch directories: %w", err)
	}

	if interactive {
		if len(directories) == 0 {
			fmt.Println("No matching scratch directories found.")
			return nil
		}

		var dirListBuilder strings.Builder
		for _, fullPath := range directories {
			dirListBuilder.WriteString(filepath.Base(fullPath))
			dirListBuilder.WriteString("\n")
		}
		dirListString := dirListBuilder.String()

		fzfCmd := exec.Command("fzf", "--tac", "--ansi", "--no-preview", "--prompt", "Select Scratch Directory> ")
		fzfCmd.Stdin = strings.NewReader(dirListString)
		fzfCmd.Stderr = os.Stderr

		selectedBaseName, err := cmd.CaptureCmd(*fzfCmd)
		if err != nil {
			if strings.TrimSpace(selectedBaseName) == "" {
				return fmt.Errorf("no directory selected")
			}
			return fmt.Errorf("failed to run fzf: %w", err)
		}

		trimmedSelectedBaseName := strings.TrimSpace(selectedBaseName)
		fullPath := filepath.Join(absScratchPath, trimmedSelectedBaseName)

		fmt.Println(fullPath)
	} else {
		if len(directories) == 0 {
			return nil
		}

		for _, fullPath := range directories {
			fmt.Println(filepath.Base(fullPath))
		}
	}

	return nil
}

// listCommand defines the 'list' subcommand for scratch.
var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List directories matching YYYY-MM-DD-<name> in the scratch path. Interactively select one if --interactive is set.",
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
