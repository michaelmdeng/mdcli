module github.com/mdcli

go 1.20

require (
	github.com/mdcli/wiki v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.25.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/mdcli/cmd v0.0.0-00010101000000-000000000000 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
)

replace github.com/mdcli/wiki => ./src/wiki

replace github.com/mdcli/cmd => ./src/cmd
