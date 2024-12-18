package main

import (
	"context"
	"os"

	"github.com/mdcli/k8s"
	"github.com/mdcli/rm"
	"github.com/mdcli/tidb"
	"github.com/mdcli/tmux"
	"github.com/mdcli/wiki"
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
