package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/michaelmdeng/mdcli/completion"
	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/michaelmdeng/mdcli/k8s"
	"github.com/michaelmdeng/mdcli/rm"
	"github.com/michaelmdeng/mdcli/scratch"
	"github.com/michaelmdeng/mdcli/tidb"
	"github.com/michaelmdeng/mdcli/split"
	"github.com/michaelmdeng/mdcli/tmux"
	"github.com/michaelmdeng/mdcli/wiki"
	"github.com/urfave/cli/v2"
)

const (
	Version = "0.0.1"
)

func CreateApp(cfg config.Config) cli.App {
	return cli.App{
		Metadata: map[string]interface{}{
			"config": cfg,

		},
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
			scratch.BaseCommand(),
			completion.BaseCommand(),
			split.BaseCommand(),
		},
	}
}

func main() {
	color.NoColor = false

	cfg := config.LoadConfig()
	app := CreateApp(cfg)
	app.Run(os.Args)
}
