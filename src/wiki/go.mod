module github.com/mdcli/wiki

go 1.20

replace github.com/mdcli/cmd => ../cmd

require (
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.25.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
)
