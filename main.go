package main

import (
	"context"
	"os"

	"github.com/michaelmdeng/mdcli/k8s"
	"github.com/michaelmdeng/mdcli/rm"
	"github.com/michaelmdeng/mdcli/tidb"
	"github.com/michaelmdeng/mdcli/tmux"
	"github.com/michaelmdeng/mdcli/wiki"
	"github.com/urfave/cli/v3"
)

const (
	Version = "0.0.1"
)

var (
	tasks = []string{"k8s", "rm", "wiki", "tidb", "tmux"}
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Name:                 "mdcli",
		Usage:                "Personal CLI",
		// Authors: []*cli.Author{
		// 	{
		// 		Name:  "Michael Deng",
		// 		Email: "michaelmdeng@gmail.com",
		// 	},
		// },
		Version: Version,
		Commands: []*cli.Command{
			k8s.BaseCommand(),
			rm.BaseCommand(),
			wiki.BaseCommand(),
			tidb.BaseCommand(),
			tmux.BaseCommand(),
		},
	}

	cmd.Run(context.Background(), os.Args)
}
