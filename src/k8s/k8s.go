package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mdcli/cmd"
)

var (
	inferableCmds = map[string]struct{} {
		"exec": {},
		"logs": {},
		"port-forward": {},
	}
	confirmableCmds = map[string]struct{} {
		"annotate": {},
		"delete": {},
		"patch": {},
	}
	editCmds = map[string]struct{} {
		"exec": {},
		"port-forward": {},
		"annotate": {},
		"delete": {},
		"patch": {},
	}
	resourceModifiableCmds = map[string]struct{} {
		"exec": {},
		"logs": {},
		"port-forward": {},
	}
	modifiableResources = map[string]struct{} {
		"deploy": {},
		"deployment": {},
		"statefulset": {},
		"sts": {},
		"service": {},
		"svc": {},
		"job": {},
	}
)

func isInferableCmd(cmd string) bool {
	_, ok := inferableCmds[cmd]
	return ok
}

func isConfirmableCmd(cmd string) bool {
	_, ok := confirmableCmds[cmd]
	return ok
}

func isEditCmd(cmd string) bool {
	_, ok := editCmds[cmd]
	return ok
}

func isResourceModifiableCmd(cmd string) bool {
	_, ok := resourceModifiableCmds[cmd]
	return ok
}

func isModifiableResource(resource string) bool {
	_, ok := modifiableResources[resource]
	return ok
}

func GetContextInteractive(pattern string) (string, error) {
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

func GetNamespaceInteractive(context string, pattern string) (string, error) {
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

func ParseContext(context string, interactive bool, pattern string, strict bool) (string, error) {
	if context != "" {
		return context, nil
	}

	if interactive && context == "" {
		var err error
		context, err = GetContextInteractive(pattern)
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

func ParseNamespace(namespace string, allNamespaces bool, interactive bool, context string, pattern string, strict bool) (string, bool, error) {
	if allNamespaces || namespace == "*" {
		return "", true, nil
	}

	if namespace != "" {
		return namespace, false, nil
	}

	if interactive && !allNamespaces && namespace == "" {
		var err error
		namespace, err = GetNamespaceInteractive(context, pattern)
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
