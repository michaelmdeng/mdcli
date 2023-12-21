package main

import (
	"os"

	"github.com/mdcli/k8s"
	"github.com/mdcli/tmux"
	"github.com/mdcli/wiki"
	"github.com/urfave/cli/v2"
)

const (
	Version = "0.0.1"
)

func main() {
	app := &cli.App{
		Name:  "mdcli",
		Usage: "mdeng personal CLI",
		Authors: []*cli.Author{
			{
				Name:  "mdeng",
				Email: "michaelmdeng@gmail.com",
			},
		},
		Version: Version,
		Commands: []*cli.Command{
			k8s.BaseCommand(),
			wiki.BaseCommand(),
			tmux.BaseCommand(),
		},
	}

	app.Run(os.Args)
}
