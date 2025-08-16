package workspace

import (
	"github.com/urfave/cli/v2"
)

const (
	workspaceUsage = `Manage workspaces via git worktrees.`
)

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "workspace",
		Aliases: []string{"ws"},
		Usage:   workspaceUsage,
		Subcommands: []*cli.Command{
			newCommand,
		},
	}
}
