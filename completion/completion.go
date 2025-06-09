package completion

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const zshAutocompleteScript = `_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

compdef _cli_zsh_autocomplete mdcli
`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:  "completion",
		Usage: "Generate shell completions",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "zsh",
				Usage:   "Output zsh completion script",
				Aliases: []string{"z"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("zsh") {
				// Print the embedded script directly
				fmt.Print(zshAutocompleteScript)
				return nil
			}
			return cli.Exit("specify shell for completion (e.g., --zsh)", 1)
		},
	}
}
