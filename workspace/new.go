package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/urfave/cli/v2"
)

func newAction(cCtx *cli.Context) error {
	if cCtx.NArg() != 1 {
		return cli.Exit("exactly one argument <git-url> must be provided", 2)
	}
	gitURL := cCtx.Args().Get(0)

	var (
		name                 = cCtx.String("name")
		projectDir           = cCtx.String("workspace-dir")
		gitBranch            = cCtx.String("git-branch")
		initializeWorktree   = cCtx.Bool("initialize-worktree")
		defaultWorktreeName  = cCtx.String("default-worktree-name")
	)

	if name == "" {
		name = extractRepoName(gitURL)
	}

	if projectDir == "" {
		cfgInterface, ok := cCtx.App.Metadata["config"]
		if !ok {
			return cli.Exit("config not found in app metadata", 1)
		}

		cfg, ok := cfgInterface.(config.Config)
		if !ok {
			return cli.Exit("invalid config type", 1)
		}
		projectDir = cfg.WorkspaceDir
	}
	var err error
	projectDir, err = filepath.Abs(projectDir)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to resolve project directory path: %v", err), 1)
	}

	if gitBranch == "" {
		gitBranch = "main"
	}

	if defaultWorktreeName == "" {
		defaultWorktreeName = gitBranch
	}

	workspacePath := filepath.Join(projectDir, name)
	
	if _, err := os.Stat(workspacePath); err == nil {
		return cli.Exit(fmt.Sprintf("workspace directory already exists: %s", workspacePath), 1)
	} else if !os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("failed to check workspace directory: %v", err), 1)
	}
	
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return cli.Exit(fmt.Sprintf("failed to create workspace directory: %v", err), 1)
	}

	gitDir := filepath.Join(workspacePath, ".git")
	cmd := exec.Command("git", "clone", "--bare", gitURL, gitDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return cli.Exit(fmt.Sprintf("failed to clone repository: %s", string(output)), 1)
	}

	if initializeWorktree {
		worktreesDir := filepath.Join(workspacePath, "worktrees")
		if err := os.MkdirAll(worktreesDir, 0755); err != nil {
			return cli.Exit(fmt.Sprintf("failed to create worktrees directory: %v", err), 1)
		}

		worktreePath := filepath.Join(worktreesDir, defaultWorktreeName)
		
		cmd := exec.Command("git", "worktree", "add", worktreePath, gitBranch)
		cmd.Dir = workspacePath
		if output, err := cmd.CombinedOutput(); err != nil {
			return cli.Exit(fmt.Sprintf("failed to create worktree: %s", string(output)), 1)
		}
	}

	fmt.Printf("Workspace created at: %s\n", workspacePath)
	return nil
}

func extractRepoName(gitURL string) string {
	// Extract repo name from URL
	// Handle both SSH and HTTPS formats
	name := gitURL
	
	// Remove trailing slash if present
	if len(name) > 0 && name[len(name)-1] == '/' {
		name = name[:len(name)-1]
	}
	
	// Remove .git extension if present
	if len(name) > 4 && name[len(name)-4:] == ".git" {
		name = name[:len(name)-4]
	}
	
	// Extract the last part after the final slash
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	
	// Handle SSH format (e.g., git@github.com:owner/repo)
	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}
	
	return name
}

var newCommand = &cli.Command{
	Name:      "new",
	Usage:     "Create a new workspace with optional worktree",
	ArgsUsage: "<git-url>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Aliases: []string{"n"},
			Usage: "the name of the workspace/project, defaults to the name of the repo",
		},
		&cli.StringFlag{
			Name:  "project-dir",
			Usage: "the directory to create the workspace in",
		},
		&cli.StringFlag{
			Name:  "git-branch",
			Aliases: []string{"branch", "b"},
			Usage: "the git ref to checkout, defaults to `main`",
		},
		&cli.BoolFlag{
			Name:  "initialize-worktree",
			Usage: "whether to create a worktree for the project",
			Value: true, // Default to true
		},
		&cli.StringFlag{
			Name:  "default-worktree-name",
			Usage: "the name of the default worktree, defaults to $BRANCH",
		},
	},
	Action: newAction,
}
