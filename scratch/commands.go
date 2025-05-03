package scratch

import (
	"github.com/urfave/cli/v2"
)

const (
	scratchUsage = `Scratchpad commands`
)

// BaseCommand returns the base command for the scratch subcommand.
func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:        "scratch",
		Aliases:     []string{"s"},
		Usage:       scratchUsage,
		Subcommands: []*cli.Command{
			newCommand,
			listCommand,
		},
	}
}
