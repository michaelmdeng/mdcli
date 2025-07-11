package scratch

import (
	"github.com/urfave/cli/v2"
)

const (
	scratchUsage = `Manage dated scratch directories for temporary notes and prototypes.`
)

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
