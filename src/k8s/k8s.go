package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mdcli/cmd"
)

func getContextInteractive() (string, error) {
	contextCmd := "yq eval '.contexts[].name' /Users/michael_deng/.kube/config"
	c := exec.Command("fzf", "--ansi", "--no-preview")
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", contextCmd),
	)
	context, err := cmd.CaptureCmd(*c)
	if err != nil {
		return "", err
	}
	if context == "" {
		return "", errors.New("no context selected")
	}

	return strings.TrimSpace(context), nil
}

func getNamespaceInteractive(context string) (string, error) {
	var namespaceCmd string
	if context == "" {
		namespaceCmd = "kubectl get ns -o jsonpath='{range .items[*]}{.metadata.name}{\"\\n\"}{end}'"
	} else {
		namespaceCmd = fmt.Sprintf("kubectl --context %s get ns -o jsonpath='{range .items[*]}{.metadata.name}{\"\\n\"}{end}'", context)
	}

	c := exec.Command("fzf", "--ansi", "--no-preview")
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", namespaceCmd),
	)
	namespace, err := cmd.CaptureCmd(*c)
	if err != nil {
		return "", err
	}
	if namespace == "" {
		return "", errors.New("no namespace selected")
	}

	return strings.TrimSpace(namespace), nil
}
