package tmux

import (
	"github.com/urfave/cli/v2"
)

const tmuxUsage = `Provides commands for manipulating custom tmux layouts.

tmux provides the built-in main-vertical layout, which lays out one large pane
on the left and the rest of the panes tiled vertically on the right. However,
this layout divides horizontal space evenly between the left and right side.
These commands provide a default pane layout that uses a 2:1 ratio for the left
and right side, which allows an editor to lay out two full-width columns in the
main left pane, essentially allowing three columns.

While this default pane layout is suitable for wide screens, narrow screens
often don't have enough horizontal real estate to support this layout. These
commands support switching between a "pane-based" window layout, where all
panes are layed out in a single window according to the default pane layout,
and a "window-based" window layout, where the rest of the panes on the right
side are moved to a dedicated window. This allows sessions to easily switch
between narrow and wide screens.`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "tmux",
		Aliases: []string{"t"},
		Usage:   tmuxUsage,
		Subcommands: []*cli.Command{
			layoutCommand(),
			switchCommand(),
			panesCommand(),
			windowsCommand(),
			toggleCommand(),
		},
	}
}

const layoutUsage = `Set the default pane layout`

func layoutCommand() *cli.Command {
	return &cli.Command{
		Name:    "layout",
		Aliases: []string{"l"},
		Usage:   layoutUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   "",
				Usage:   "`SESSION` to set the default layout for. Defaults to the current session",
			},
			&cli.StringFlag{
				Name:    "window",
				Aliases: []string{"w"},
				Value:   "",
				Usage:   "`WINDOW` to set the default layout for. Defaults to all windows in the session.",
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

const switchUsage = `Switch to the corresponding pane in a "window-based" layout`

func switchCommand() *cli.Command {
	return &cli.Command{
		Name:    "switch",
		Aliases: []string{"s"},
		Usage:   switchUsage,
		Action: func(cCtx *cli.Context) error {
			return switchExtraPane()
		},
	}
}

const toggleUsage = `Toggle the current window layout`

func toggleCommand() *cli.Command {
	return &cli.Command{
		Name:    "toggle",
		Aliases: []string{"t"},
		Usage:   toggleUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   "",
				Usage:   "`SESSION` to toggle the window layout for. Defaults to the current session",
			},
			&cli.StringFlag{
				Name:    "window",
				Aliases: []string{"w"},
				Value:   "",
				Usage:   "`WINDOW` to toggle the window layout for. Defaults to all windows in the session.",
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

const panesUsage = `Switch the current window layout to pane-based`

func panesCommand() *cli.Command {
	return &cli.Command{
		Name:    "panes",
		Aliases: []string{"p"},
		Usage:   panesUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   "",
				Usage:   "`SESSION` to change the window layout for. Defaults to the current session.",
			},
			&cli.StringFlag{
				Name:    "window",
				Aliases: []string{"w"},
				Value:   "",
				Usage:   "`WINDOW` to change the window layout for. Defaults to all windows in the session.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return setPaneWindowLayout(cCtx.String("session"), cCtx.String("window"))
		},
	}
}

const windowsUsage = `Switch the current window layout to window-based`

func windowsCommand() *cli.Command {
	return &cli.Command{
		Name:    "windows",
		Aliases: []string{"w"},
		Usage:   windowsUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   "",
				Usage:   "`SESSION` to change the window layout for. Defaults to the current session.",
			},
			&cli.StringFlag{
				Name:    "window",
				Aliases: []string{"w"},
				Value:   "",
				Usage:   "`WINDOW` to change the window layout for. Defaults to all windows in the session.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return setWindowWindowLayout(cCtx.String("session"), cCtx.String("window"))
		},
	}
}
