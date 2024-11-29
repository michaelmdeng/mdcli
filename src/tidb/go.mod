module github.com/mdcli/tidb

go 1.21.11

replace github.com/mdcli/cmd => ../cmd

replace github.com/mdcli/config => ../config

replace github.com/mdcli/k8s => ../k8s

require (
	github.com/bitfield/script v0.22.0
	github.com/fatih/color v1.18.0
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000
	github.com/mdcli/config v0.0.0-00010101000000-000000000000
	github.com/mdcli/k8s v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.25.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.25.0 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)
