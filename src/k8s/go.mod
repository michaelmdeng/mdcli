module github.com/mdcli/k8s

go 1.20

replace github.com/mdcli/cmd => ../cmd

require (
	github.com/bitfield/script v0.22.0
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v3 v3.0.0-beta1
)

require (
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)
