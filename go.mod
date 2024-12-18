module github.com/mdcli

go 1.21.13

require (
	github.com/mdcli/k8s v0.0.0-00010101000000-000000000000
	github.com/mdcli/rm v0.0.0-00010101000000-000000000000
	github.com/mdcli/tidb v0.0.0-00010101000000-000000000000
	github.com/mdcli/tmux v0.0.0-00010101000000-000000000000
	github.com/mdcli/wiki v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v3 v3.0.0-beta1
)

require (
	github.com/bitfield/script v0.22.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000 // indirect
	github.com/mdcli/config v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.25.0 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)

replace github.com/mdcli/config => ./src/config

replace github.com/mdcli/cmd => ./src/cmd

replace github.com/mdcli/k8s => ./src/k8s

replace github.com/mdcli/tidb => ./src/tidb

replace github.com/mdcli/tmux => ./src/tmux

replace github.com/mdcli/wiki => ./src/wiki

replace github.com/mdcli/rm => ./src/rm
