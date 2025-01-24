package tidb

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func BasePdCommand() *cli.Command {
	return &cli.Command{
		Name:    "pd",
		Usage:   `Commands for handling PD on K8s`,
		Subcommands: []*cli.Command{
			pdTsoCommand(),
		},
	}
}

func pdTsoCommand() *cli.Command {
	return &cli.Command{
		Name:  "tso",
		Usage: "Perform TSO/timestamp conversion",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "tso",
				Value:   0,
				Usage:   "TSO to convert",
			},
			&cli.StringFlag{
				Name:    "ts",
				Value:   "",
				Usage:   "timestamp to convert",
			},
		},
		Action: func(cCtx *cli.Context) error {
			tso := cCtx.Int("tso")
			ts := cCtx.String("ts")
			if tso == 0 && ts == "" {
				return errors.New("one of tso or ts must be provided")
			}

			if tso != 0 && ts != "" {
				return errors.New("only one of tso or ts can be provided")
			}

			if tso != 0 {
				t := time.Unix(int64((tso / 1000) >> 18), 0).UTC()
				ts := t.Format(time.RFC3339)
				fmt.Println(ts)
			} else {
				t, err := time.Parse(time.RFC3339, ts)
				if err != nil {
					return err
				}
				ts := t.Unix()
				tso := (ts << 18) * 1000
				fmt.Println(tso)
			}

			return nil
		},
	}
}

