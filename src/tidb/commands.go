package tidb

import (
	"github.com/urfave/cli/v2"
)

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidb",
		Aliases: []string{"ti", "tdb"},
		Usage:   `Commands for managing TiDB`,
		Subcommands: []*cli.Command{
		},
	}
}
