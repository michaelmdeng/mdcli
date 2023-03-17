package main

import (
	"os"

	"github.com/mdcli/tmux"
	"github.com/mdcli/wiki"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mdcli",
		Usage: "Custom CLI",
		Commands: []*cli.Command{
			wiki.BaseCommand(),
			tmux.BaseCommand(),
		},
	}

	app.Run(os.Args)
}
