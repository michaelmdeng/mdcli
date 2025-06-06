package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

// generateTmuxinatorConfig creates a temporary tmuxinator config file.
func generateTmuxinatorConfig(templatePath, projectName, projectRoot string) (string, error) {
	if templatePath == "" {
		return "", fmt.Errorf("tmuxinator template path is not configured")
	}
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("tmuxinator template file not found at '%s'", templatePath)
	} else if err != nil {
		return "", fmt.Errorf("failed to check tmuxinator template file '%s': %w", templatePath, err)
	}

	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read tmuxinator template file '%s': %w", templatePath, err)
	}

	// Replace placeholders
	content := string(templateContent)
	content = strings.ReplaceAll(content, "{{PROJECT_NAME}}", projectName)
	content = strings.ReplaceAll(content, "{{PROJECT_ROOT}}", projectRoot)

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("tmuxinator-%s-*.yaml", projectName))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary tmuxinator config file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write to temporary tmuxinator config file: %w", err)
	}

	return tmpFile.Name(), nil
}

// runTmuxinator starts a tmuxinator session using the generated config file.
func runTmuxinator(configPath string) error {
	// Use --project-config to specify the temporary config file path
	err := cmd.RunCommand("tmuxinator", "start", "--project-config", configPath)
	if err != nil {
		return fmt.Errorf("failed to start tmuxinator session: %w", err)
	}
	return nil
}

// tmuxAction implements the logic for the 'scratch tmux' command.
func tmuxAction(cCtx *cli.Context) error {
	if cCtx.NArg() != 1 {
		return cli.Exit("exactly one argument <name> must be provided", 2)
	}
	name := cCtx.Args().Get(0)

	var scratchPath string
	if cCtx.IsSet("scratch-path") {
		scratchPath = cCtx.String("scratch-path")
	} else {
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
	if scratchPath == "" {
		return cli.Exit("scratch path not configured", 1)
	}
	absScratchPath, err := expandPath(scratchPath)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("scratch directory '%s' does not exist", absScratchPath), 1)
	} else if err != nil {
		return cli.Exit(fmt.Sprintf("failed to check scratch directory '%s': %v", absScratchPath, err), 1)
	}

	var tmuxinatorTemplate string
	if cCtx.IsSet("tmuxinator-template") {
		tmuxinatorTemplate = cCtx.String("tmuxinator-template")
	} else {
		cfgInterface, ok := cCtx.App.Metadata["config"]
		if !ok {
			return cli.Exit("configuration not found in application metadata", 1)
		}
		cfg, ok := cfgInterface.(config.Config)
		if !ok {
			return cli.Exit("invalid configuration type in application metadata", 1)
		}
		tmuxinatorTemplate = cfg.Scratch.TmuxinatorTemplate
	}

	if tmuxinatorTemplate != "" && !filepath.IsAbs(tmuxinatorTemplate) {
		absTemplatePath, err := filepath.Abs(tmuxinatorTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not make tmuxinator template path absolute: %v\n", err)
		} else {
			tmuxinatorTemplate = absTemplatePath
		}
	}
	if tmuxinatorTemplate == "" {
		return cli.Exit("tmuxinator template path not configured", 1)
	}

	targetDir, err := findScratchDirectory(absScratchPath, name)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error finding scratch directory: %v", err), 1)
	}

	if targetDir == "" {
		createReadme := cCtx.Bool("create-readme")
		newDirPath, err := createScratchDirectory(absScratchPath, name, createReadme)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		fmt.Printf("Created scratch directory: %s\n", newDirPath)
		targetDir = newDirPath
	}

	projectName := filepath.Base(targetDir)

	tmpConfigPath, err := generateTmuxinatorConfig(tmuxinatorTemplate, projectName, targetDir)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to generate tmuxinator config: %v", err), 1)
	}
	defer os.Remove(tmpConfigPath)

	fmt.Printf("Starting tmuxinator session '%s' for project '%s'...\n", projectName, targetDir)
	if err := runTmuxinator(tmpConfigPath); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}

// tmuxCommand defines the 'tmux' subcommand for scratch.
var tmuxCommand = &cli.Command{
	Name:      "tmux",
	Aliases:   []string{"t"},
	Usage:     "Find or create a scratch directory and launch tmuxinator",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "scratch-path",
			Usage: "Override the scratch path from config",
		},
		&cli.StringFlag{
			Name:  "tmuxinator-template",
			Usage: "Override the tmuxinator template path from config",
		},
		&cli.BoolFlag{
			Name:    "create-readme",
			Aliases: []string{"r"},
			Usage:   "Create an empty README.md if a new directory is created",
			Value:   true, // Default to true
		},
	},
	Action: tmuxAction,
}
