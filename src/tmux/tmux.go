package tmux

import (
	"fmt"
	"strings"

	"github.com/mdcli/cmd"
)

func currentWindow() (string, string, error) {
	output, err := cmd.CaptureCommand(
		"tmux", "display-message",
		"-p", "#S:#W",
	)
	if err != nil {
		return "", "", err
	}

	splitOutput := strings.Split(strings.TrimSpace(output), ":")
	return splitOutput[0], splitOutput[1], nil

}

func listPanes(session string, window string) ([]string, error) {
	output, err := cmd.CaptureCommand(
		"tmux", "list-panes",
		"-t", fmt.Sprintf("%v:%v", session, window),
		"-F", "#P",
	)
	if err != nil {
		return []string{}, err
	}

	panes := make([]string, 0)
	for _, pane := range strings.Split(output, "\n") {
		if pane != "" {
			panes = append(panes, pane)
		}
	}
	return panes, nil
}

func listWindows(session string) ([]string, error) {
	output, err := cmd.CaptureCommand(
		"tmux", "list-windows",
		"-t", session,
		"-F", "#W",
	)
	if err != nil {
		return []string{}, err
	}

	windows := make([]string, 0)
	for _, window := range strings.Split(output, "\n") {
		if window != "" {
			windows = append(windows, window)
		}
	}
	return windows, nil
}

func newWindow(session string, window string) error {
	return cmd.RunCommand(
		"tmux", "new-window", "-d",
		"-t", session,
		"-n", window,
	)
}

func selectLayout(session string, window string, layout string) error {
	return cmd.RunCommand(
		"tmux", "select-layout",
		"-t", fmt.Sprintf("%v:%v", session, window),
		layout,
	)
}

func movePane(
	sessionFrom string, windowFrom string,
	sessionTo string, windowTo string,
	pane string,
) error {
	return cmd.RunCommand(
		"tmux", "move-pane", "-d",
		"-s", fmt.Sprintf("%v:%v.%v", sessionFrom, windowFrom, pane),
		"-t", fmt.Sprintf("%v:%v", sessionTo, windowTo),
	)
}

func killPane(session string, window string, pane string) error {
	return cmd.RunCommand(
		"tmux", "kill-pane",
		"-t", fmt.Sprintf("%v:%v.%v", session, window, pane),
	)
}

func isPaneBased(session string) (bool, error) {
	windows, err := listWindows(session)
	if err != nil {
		return false, err
	}

	for _, window := range windows {
		if strings.Contains(window, "-extra") {
			return false, nil
		}
	}

	return true, nil

}

func isWindowBased(session string) (bool, error) {
	windows, err := listWindows(session)
	if err != nil {
		return false, err
	}

	for _, window := range windows {
		if isExtraWindow(window) {
			return true, nil
		}
	}

	return false, nil

}

func isExtraWindow(window string) bool {
	return strings.Contains(window, "-extra")
}

func extraWindowName(window string) string {
	return fmt.Sprintf("%v-extra", window)
}

func mainWindowName(window string) string {
	return strings.Replace(window, "-extra", "", -1)
}

func setWindowWindowLayout(session string, window string) error {
	windows, err := listWindows(session)
	if err != nil {
		return err
	}

	mainWindows := make(map[string]bool)
	extraWindows := make(map[string]bool)
	for _, w := range windows {
		if isExtraWindow(w) {
			extraWindows[w] = true
		} else {
			mainWindows[w] = true
		}
	}

	for mainWindow, _ := range mainWindows {
		if window != "" && mainWindow != window {
			continue
		}

		panes, err := listPanes(session, mainWindow)
		if err != nil {
			return err
		}

		extraPanes := panes[1:]
		if len(extraPanes) == 0 {
			continue
		}

		createdExtraWindow := false
		extraWindow := extraWindowName(mainWindow)
		if _, ok := extraWindows[extraWindow]; !ok {
			err := newWindow(session, extraWindow)
			if err != nil {
				return err
			}
			createdExtraWindow = true
		}

		// Panes will be renumbered after move-pane, so we need to move
		// one-at-a-time and re-check panes after every move.
		// This is preferable to iterating in reverse-order so that we
		// maintain the proper pane order in the extra window.
		for {
			err := movePane(
				session, mainWindow,
				session, extraWindow,
				extraPanes[0])
			if err != nil {
				return err
			}

			panes, err := listPanes(session, mainWindow)
			if err != nil {
				return err
			}

			extraPanes = panes[1:]
			if len(extraPanes) == 0 {
				break
			}
		}

		if createdExtraWindow {
			err := killPane(session, extraWindow, "0")
			if err != nil {
				return err
			}
		}

		err = selectLayout(session, extraWindow, "tiled")
		if err != nil {
			return err
		}
	}

	return nil
}

func setPaneWindowLayout(session string, window string) error {
	windows, err := listWindows(session)
	if err != nil {
		return err
	}

	mainWindows := make(map[string]bool)
	extraWindows := make(map[string]bool)
	for _, w := range windows {
		if isExtraWindow(w) {
			extraWindows[w] = true
		} else {
			mainWindows[w] = true
		}
	}

	for extraWindow, _ := range extraWindows {
		mainWindow := mainWindowName(extraWindow)

		for {
			extraPanes, err := listPanes(session, extraWindow)
			if err != nil {
				break
			}

			if len(extraPanes) == 0 {
				break
			}

			err = movePane(
				session, extraWindow,
				session, mainWindow,
				extraPanes[0],
			)
			if err != nil {
				return err
			}
		}

		err = selectLayout(session, mainWindow, "main-vertical")
		if err != nil {
			return err
		}
	}

	err = setDefaultLayout(session, window)
	if err != nil {
		return err
	}

	return nil
}

func switchExtraPane() error {
	session, window, err := currentWindow()
	if err != nil {
		return err
	}

	var switchWindow string
	if isExtraWindow(window) {
		switchWindow = mainWindowName(window)
	} else {
		switchWindow = extraWindowName(window)
	}

	return cmd.RunCommand(
		"tmux", "select-window",
		"-t", fmt.Sprintf("%v:%v", session, switchWindow),
	)
}
