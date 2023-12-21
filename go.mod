module github.com/mdcli

go 1.20

require (
	github.com/mdcli/k8s v0.0.0-00010101000000-000000000000
	github.com/mdcli/tmux v0.0.0-00010101000000-000000000000
	github.com/mdcli/wiki v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.25.0
)

require (
	github.com/bitfield/script v0.22.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)

replace github.com/mdcli/cmd => ./src/cmd

replace github.com/mdcli/k8s => ./src/k8s

replace github.com/mdcli/wiki => ./src/wiki

replace github.com/mdcli/tmux => ./src/tmux
