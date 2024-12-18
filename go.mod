module github.com/michaelmdeng/mdcli

go 1.21.13

require (
	github.com/michaelmdeng/mdcli/k8s v0.0.0-00010101000000-000000000000
	github.com/michaelmdeng/mdcli/rm v0.0.0-00010101000000-000000000000
	github.com/michaelmdeng/mdcli/tidb v0.0.0-00010101000000-000000000000
	github.com/michaelmdeng/mdcli/tmux v0.0.0-00010101000000-000000000000
	github.com/michaelmdeng/mdcli/wiki v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v3 v3.0.0-beta1
)

require (
	github.com/bitfield/script v0.22.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/michaelmdeng/mdcli/cmd v0.0.0-00010101000000-000000000000 // indirect
	github.com/michaelmdeng/mdcli/config v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.25.0 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)

replace github.com/michaelmdeng/mdcli/config => ./src/config

replace github.com/michaelmdeng/mdcli/cmd => ./src/cmd

replace github.com/michaelmdeng/mdcli/k8s => ./src/k8s

replace github.com/michaelmdeng/mdcli/tidb => ./src/tidb

replace github.com/michaelmdeng/mdcli/tmux => ./src/tmux

replace github.com/michaelmdeng/mdcli/wiki => ./src/wiki

replace github.com/michaelmdeng/mdcli/rm => ./src/rm
