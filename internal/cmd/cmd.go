package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CaptureCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CaptureCmd(cmd exec.Cmd) (string, error) {
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func RunCommandDiscardOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExitError creates cli.Exit errors, extracting exit code from exec.ExitError if possible.
// Moved from tidb/commands.go
func ExitError(err error) error {
	if err == nil {
		return nil
	}
	exitCode := 1 // Default exit code for generic errors
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}
	// Use fmt.Sprintf to ensure we pass a string message to cli.Exit
	return cli.Exit(fmt.Sprintf("%v", err), exitCode)
}

func GetConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	numAttempts := 3
	for range numAttempts {
		fmt.Fprintf(os.Stderr, "%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}

	return false
}

func IsPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true
	}

	fi, err = os.Stdout.Stat()
	if err != nil {
		panic(err)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true
	}

	return false
}
