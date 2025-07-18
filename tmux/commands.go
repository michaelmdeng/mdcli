package tmux

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const tmuxUsage = `Commands for manipulating custom tmux layouts.`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "tmux",
		Aliases: []string{"tm", "tx"},
		Usage:   tmuxUsage,
		Subcommands: []*cli.Command{
			layoutCommand(),
			shellCommand(),
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

			var currSession string
			var currWindow string
			aggregateErrors := false
			aggregatedErrors := []error{}
			if len(windows) > 1 {
				aggregateErrors = true
				currSession, currWindow, err = currentWindow()
				if err != nil {
					return nil
				}
			}

			for _, window := range windows {
				err = selectWindow(session, window)
				if err != nil {
					if aggregateErrors {
						aggregatedErrors = append(aggregatedErrors, err)
					} else {
						return err
					}
				}

				err := setDefaultLayout(session, window)
				if err != nil {
					if aggregateErrors {
						aggregatedErrors = append(aggregatedErrors, err)
					} else {
						return err
					}
				}
			}

			if len(currWindow) > 0 {
				err = selectWindow(currSession, currWindow)
				if err != nil {
					return err
				}
			}

			if aggregateErrors && len(aggregatedErrors) > 0 {
				return fmt.Errorf("encountered errors in setting layout: %v", aggregatedErrors)
			}

			return nil
		},
	}
}

const switchUsage = `Switch to the corresponding pane in a "window-based" layout`

func switchCommand() *cli.Command {
	return &cli.Command{
		Name:    "switch",
		Aliases: []string{"sw"},
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

const shellUsage = `Get the fully-qualified target for the current window's shell pane`

func shellCommand() *cli.Command {
	return &cli.Command{
		Name:    "shell",
		Aliases: []string{"sh"},
		Usage:   shellUsage,
		Action: func(cCtx *cli.Context) error {
			session, window, err := currentWindow()
			if err != nil {
				return err
			}

			fmt.Printf("%s:%s.1", session, window)
			return nil
		},
	}
}
