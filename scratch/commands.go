package scratch

import (
	"github.com/urfave/cli/v2"
)

const (
	scratchUsage = `Manage temporary scratch directories for notes, code snippets, or experiments.
These commands help you quickly create, list, and open dated scratch directories,
optionally integrating with tmuxinator for a pre-configured terminal environment.`
)

// BaseCommand returns the base command for the scratch subcommand.
func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "scratch",
		Aliases: []string{"s"},
		Usage:   scratchUsage,
		Subcommands: []*cli.Command{
			newCommand,
			listCommand,
			tmuxCommand,
		},
	}
}
