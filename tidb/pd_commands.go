package tidb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

func BasePdCommand() *cli.Command {
	return &cli.Command{
		Name:  "pd",
		Usage: `Commands for handling PD on K8s`,
		Subcommands: []*cli.Command{
			pdTsoCommand(),
		},
	}
}

func pdTsoCommand() *cli.Command {
	return &cli.Command{
		Name:      "tso",
		Usage:     "Perform TSO/timestamp conversion. Provide either a TSO (int) or a timestamp (RFC3339 string).",
		ArgsUsage: "<tso_or_timestamp>",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				return cli.Exit("exactly one argument (tso or timestamp) must be provided", 1)
			}
			input := cCtx.Args().Get(0)

			// Try parsing as TSO (integer)
			tso, err := strconv.ParseInt(input, 10, 64)
			if err == nil {
				// Successfully parsed as TSO
				t := time.Unix(int64((tso/1000)>>18), 0).UTC()
				ts := t.Format(time.RFC3339)
				fmt.Println(ts)
				return nil
			}

			// Try parsing as timestamp (RFC3339 string)
			t, err := time.Parse(time.RFC3339, input)
			if err == nil {
				// Successfully parsed as timestamp
				tsUnix := t.Unix()
				tsoResult := (tsUnix << 18) * 1000
				fmt.Println(tsoResult)
				return nil
			} else {
				// If both parsing attempts fail, return an error.
				return cli.Exit(fmt.Sprintf("input '%s' is not a valid TSO (integer) or timestamp (RFC3339): %s", input, err), 1)
			}
		},
	}
}
