module github.com/michaelmdeng/mdcli/tmux

go 1.20

replace github.com/michaelmdeng/mdcli/cmd => ../cmd

require (
	github.com/michaelmdeng/mdcli/cmd v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v3 v3.0.0-beta1
)
