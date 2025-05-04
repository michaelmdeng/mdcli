package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

// findScratchDirectory searches for a directory matching the name within the scratch path.
// It prioritizes exact matches (YYYY-MM-DD-name) and then suffix matches (-name).
// Returns the full path of the found directory or an empty string if not found.
// Returns an error if multiple suffix matches are found.
func findScratchDirectory(scratchPath, name string) (string, error) {
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	entries, err := os.ReadDir(absScratchPath)
	if err != nil {
		return "", fmt.Errorf("failed to read scratch directory '%s': %w", absScratchPath, err)
	}

	var exactMatch string
	var suffixMatches []string
	suffixToCheck := "-" + name
	datePrefix := time.Now().Format("2006-01-02")
	exactName := fmt.Sprintf("%s-%s", datePrefix, name)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		fullPath := filepath.Join(absScratchPath, entryName)

		// Check for exact match (including today's date)
		if entryName == exactName {
			exactMatch = fullPath
			break // Exact match found, no need to check further
		}

		// Check for suffix match
		if strings.HasSuffix(entryName, suffixToCheck) {
			suffixMatches = append(suffixMatches, fullPath)
		}
	}

	if exactMatch != "" {
		return exactMatch, nil
	}

	if len(suffixMatches) == 1 {
		return suffixMatches[0], nil
	}

	if len(suffixMatches) > 1 {
		// Sort by name to get the most recent if dates are prefixes
		sort.Strings(suffixMatches)
		// Return the last one (most recent date assuming YYYY-MM-DD prefix)
		return suffixMatches[len(suffixMatches)-1], nil
		// Alternatively, error out:
		// return "", fmt.Errorf("multiple directories found matching suffix '%s': %v", suffixToCheck, suffixMatches)
	}

	return "", nil // No match found
}

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
	defer tmpFile.Close() // Close the file handle

	if _, err := tmpFile.WriteString(content); err != nil {
		// Attempt to remove the partially written file on error
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

	// --- Determine Scratch Path ---
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
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to get absolute path for scratch directory '%s': %v", scratchPath, err), 1)
	}
	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("scratch directory '%s' does not exist", absScratchPath), 1)
	} else if err != nil {
		return cli.Exit(fmt.Sprintf("failed to check scratch directory '%s': %v", absScratchPath, err), 1)
	}

	// --- Determine Tmuxinator Template Path ---
	var tmuxinatorTemplate string
	if cCtx.IsSet("tmuxinator-template") {
		tmuxinatorTemplate = cCtx.String("tmuxinator-template")
	} else {
		cfgInterface, ok := cCtx.App.Metadata["config"]
		if !ok {
			return cli.Exit("configuration not found in application metadata", 1) // Should not happen if scratch path worked
		}
		cfg, ok := cfgInterface.(config.Config)
		if !ok {
			return cli.Exit("invalid configuration type in application metadata", 1) // Should not happen
		}
		tmuxinatorTemplate = cfg.Scratch.TmuxinatorTemplate
	}
	// Ensure template path is absolute if provided
	if tmuxinatorTemplate != "" && !filepath.IsAbs(tmuxinatorTemplate) {
		absTemplatePath, err := filepath.Abs(tmuxinatorTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not make tmuxinator template path absolute: %v\n", err)
			// Proceed with the potentially relative path, hoping tmuxinator can find it or it fails later
		} else {
			tmuxinatorTemplate = absTemplatePath
		}
	}
	if tmuxinatorTemplate == "" {
		return cli.Exit("tmuxinator template path not configured", 1)
	}


	// --- Find or Create Directory ---
	targetDir, err := findScratchDirectory(absScratchPath, name)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error finding scratch directory: %v", err), 1)
	}

	projectName := ""
	if targetDir == "" {
		// Not found, create it
		today := time.Now().Format("2006-01-02")
		newDirName := fmt.Sprintf("%s-%s", today, name)
		newDirPath := filepath.Join(absScratchPath, newDirName)

		// Check again for existing directory with the exact new name (race condition mitigation)
		if _, err := os.Stat(newDirPath); !os.IsNotExist(err) {
			// Directory already exists (or other error), use it instead of failing Mkdir
			if err == nil {
				targetDir = newDirPath
				projectName = newDirName
				fmt.Fprintf(os.Stderr, "Warning: directory '%s' already existed, using it.\n", newDirPath)
			} else {
				return cli.Exit(fmt.Sprintf("failed to check existing directory '%s': %v", newDirPath, err), 1)
			}
		} else {
			// Proceed with creation
			if err := os.Mkdir(newDirPath, 0755); err != nil {
				return cli.Exit(fmt.Sprintf("failed to create directory '%s': %v", newDirPath, err), 1)
			}
			fmt.Printf("Created scratch directory: %s\n", newDirPath)
			targetDir = newDirPath
			projectName = newDirName
		}
	}

	if projectName == "" {
		projectName = filepath.Base(targetDir)
	}

	// --- Generate Tmuxinator Config ---
	tmpConfigPath, err := generateTmuxinatorConfig(tmuxinatorTemplate, projectName, targetDir)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to generate tmuxinator config: %v", err), 1)
	}
	// Ensure temporary file is removed even if tmuxinator fails
	defer os.Remove(tmpConfigPath)

	// --- Run Tmuxinator ---
	fmt.Printf("Starting tmuxinator session '%s' for project '%s'...\n", projectName, targetDir)
	if err := runTmuxinator(tmpConfigPath); err != nil {
		// runTmuxinator already formats the error
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
	},
	Action: tmuxAction,
}
