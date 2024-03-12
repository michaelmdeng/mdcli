package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mdcli/cmd"
)

func noop(namespace string) string {
	return namespace
}

func tidbContext(context string) string {
	return fmt.Sprintf("m-tidb-%s-ea1-us", context)
}

func tidbNamespace(namespace string) string {
	return fmt.Sprintf("tidb-%s", namespace)
}

func getContextInteractive(pattern string) (string, error) {
	var contextCmd string
	if len(pattern) > 0 {
		contextCmd = fmt.Sprintf("yq eval '.contexts[].name' /Users/michael_deng/.kube/config | grep -e \"%s\"", pattern)
	} else {
		contextCmd = "yq eval '.contexts[].name' /Users/michael_deng/.kube/config"
	}
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

func getNamespaceInteractive(context string, pattern string) (string, error) {
	var namespaceCmd string
	if context == "" {
		namespaceCmd = "kubectl get ns -o jsonpath='{range .items[*]}{.metadata.name}{\"\\n\"}{end}'"
	} else {
		namespaceCmd = fmt.Sprintf("kubectl --context %s get ns -o jsonpath='{range .items[*]}{.metadata.name}{\"\\n\"}{end}'", context)
	}
	if len(pattern) > 0 {
		namespaceCmd = fmt.Sprintf("%s | grep %s", namespaceCmd, pattern)
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

func parseContext(context string, convertContext func(string) string, interactive bool, pattern string, strict bool) (string, error) {
	if context != "" {
		return convertContext(context), nil
	}

	if interactive && context == "" {
		var err error
		context, err = getContextInteractive("^m-tidb-")
		if strict && err != nil {
			return "", err
		} else if err != nil {
			context = ""
		}
	}

	if strict && context == "" {
		return "", errors.New("context must be specified in strict mode")
	}

	return context, nil
}

func parseNamespace(namespace string, convertNamespace func(string) string, allNamespaces bool, interactive bool, context string, pattern string, strict bool) (string, bool, error) {
	if allNamespaces || namespace == "*" {
		return "", true, nil
	}

	if namespace != "" {
		namespace = convertNamespace(namespace)
	}

	if interactive && !allNamespaces && namespace == "" {
		var err error
		namespace, err = getNamespaceInteractive(context, pattern)
		if strict && err != nil {
			return "", false, err
		} else if err != nil {
			namespace = ""
		}
	}

	if strict && namespace == "" {
		return "", false, errors.New("namespace must be specified in strict mode")
	}

	return namespace, false, nil
}
