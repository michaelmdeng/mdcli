package main

import (
	"os"

	"github.com/mdcli/k8s"
	"github.com/mdcli/rm"
	"github.com/mdcli/tidb"
	"github.com/mdcli/tmux"
	"github.com/mdcli/wiki"
	"github.com/urfave/cli/v2"
)

const (
	Version = "0.0.1"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "mdcli",
		Usage:                "Personal CLI",
		Authors: []*cli.Author{
			{
				Name:  "Michael Deng",
				Email: "michaelmdeng@gmail.com",
			},
		},
		Version: Version,
		Commands: []*cli.Command{
			k8s.BaseCommand(),
			rm.BaseCommand(),
			wiki.BaseCommand(),
			tidb.BaseCommand(),
			tmux.BaseCommand(),
		},
	}

	app.Run(os.Args)
}
