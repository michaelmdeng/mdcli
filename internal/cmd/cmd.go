package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func GetConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	numAttempts := 3
	for i := 0; i < numAttempts; i++ {
		fmt.Fprintf(os.Stderr, "%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
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
