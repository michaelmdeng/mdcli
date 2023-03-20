package tmux

import (
	"github.com/urfave/cli/v2"
)

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "tmux",
		Aliases: []string{"t"},
		Usage:   "Commands for tmux",
		Subcommands: []*cli.Command{
			layoutCommand(),
			switchCommand(),
			panesCommand(),
			windowsCommand(),
			toggleCommand(),
		},
	}
}

func layoutCommand() *cli.Command {
	return &cli.Command{
		Name:    "layout",
		Aliases: []string{"l"},
		Usage:   "Set the default pane layout",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "session",
			Aliases: []string{"s"},
			Value: "",
			Usage: "`SESSION` to set the default layout for",
		    },
		    &cli.StringFlag{
			Name:  "window",
			Aliases: []string{"w"},
			Value: "",
			Usage: "`WINDOW` to set the default layout for",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			session := cCtx.String("session")
			window := cCtx.String("window")

			var windows []string
			var err error
			if window == "" {
				windows, err = listWindows(session)
				if err != nil {
					return err
				}
			} else {
				windows = []string{window}
			}

			for _, window := range windows {
				err := setDefaultLayout(session, window)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func switchCommand() *cli.Command {
	return &cli.Command{
		Name:    "switch",
		Aliases: []string{"s"},
		Usage:   "Switch to the corresponding pane",
		Action: func(cCtx *cli.Context) error {
			return switchExtraPane()
		},
	}
}

func toggleCommand() *cli.Command {
	return &cli.Command{
		Name:    "toggle",
		Aliases: []string{"t"},
		Usage:   "Toggle the current window layout",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "session",
			Aliases: []string{"s"},
			Value: "",
			Usage: "`SESSION` to toggle the window layout for",
		    },
		    &cli.StringFlag{
			Name:  "window",
			Aliases: []string{"w"},
			Value: "",
			Usage: "`WINDOW` to toggle the window layout for",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			session := cCtx.String("session")
			isWindow, err := isWindowBased(session)
			if err != nil {
				return err
			}

			window := cCtx.String("window")
			if isWindow {
				return setPaneWindowLayout(session, window)
			}

			return setWindowWindowLayout(session, window)
		},
	}
}

func panesCommand() *cli.Command {
	return &cli.Command{
		Name:    "panes",
		Aliases: []string{"p"},
		Usage:   "Change the current window layout to pane-based",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "session",
			Aliases: []string{"s"},
			Value: "",
			Usage: "`SESSION` to change the window layout for",
		    },
		    &cli.StringFlag{
			Name:  "window",
			Aliases: []string{"w"},
			Value: "",
			Usage: "`WINDOW` to change the window layout for",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			return setPaneWindowLayout(cCtx.String("session"), cCtx.String("window"))
		},
	}
}

func windowsCommand() *cli.Command {
	return &cli.Command{
		Name:    "windows",
		Aliases: []string{"w"},
		Usage:   "Change the current window layout to window-based",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "session",
			Aliases: []string{"s"},
			Value: "",
			Usage: "`SESSION` to change the window layout for",
		    },
		    &cli.StringFlag{
			Name:  "window",
			Aliases: []string{"w"},
			Value: "",
			Usage: "`WINDOW` to change the window layout for",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			return setWindowWindowLayout(cCtx.String("session"), cCtx.String("window"))
		},
	}
}

