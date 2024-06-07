package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mdcli/cmd"
)

var (
	inferableCmds = map[string] struct{}{
		"exec": {},
		"logs": {},
		"port-forward": {},
	}
	editableCmds = map[string] struct{}{
		"annotate": {},
		"delete": {},
		"patch": {},
	}
)

func isInferableCmd(cmd string) bool {
	_, ok := inferableCmds[cmd]
	return ok
}

func isEditableCmd(cmd string) bool {
	_, ok := editableCmds[cmd]
	return ok
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

func ParseContext(context string, interactive bool, pattern string, strict bool) (string, error) {
	if context != "" {
		if ctx, ok := ContextsByAlias[context]; ok {
			return ctx, nil
		} else {
			return context, nil
		}
	}

	if interactive && context == "" {
		var err error
		context, err = getContextInteractive(pattern)
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
		ns, inferred := inferNamespace(context, namespace)
		if inferred {
			return ns, false, nil
		} else {
			return namespace, false, nil
		}
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

var NamespaceVariableAliases = []string{"%n", "%ns"}
var ContextVariableAliases = []string{"%c", "%ctx"}
var TidbClusterVariableAliases = []string{"%tc", "%t"}
var AZVariableAliases = []string{"%z", "%az"}

func substituteAliases(args []string, context string, namespace string) []string {
	for i, arg := range args {
		for _, alias := range NamespaceVariableAliases {
			arg = strings.ReplaceAll(arg, alias, namespace)
		}
		for _, alias := range ContextVariableAliases {
			arg = strings.ReplaceAll(arg, alias, context)
		}
		for _, alias := range TidbClusterVariableAliases {
			tc := strings.TrimPrefix(namespace, "tidb-")
			arg = strings.ReplaceAll(arg, alias, tc)
		}
		for _, alias := range AZVariableAliases {
			// parse zone from context
			// ex. m-tidb-test-<zone>-ea1-us
			pattern := regexp.MustCompile(`m-tidb-[a-z]+-([a-z])-ea1-us`)
			matches := pattern.FindStringSubmatch(context)
			if len(matches) >= 2 {
				zone := matches[1]
				if zone == "c" {
					zone = "e"
				}
				az := fmt.Sprintf("us-east-1%s", zone)
				arg = strings.ReplaceAll(arg, alias, az)
			}
		}
		args[i] = arg
	}
	return args
}

func BuildKubectlArgs(context string, namespace string, allNamespaces bool, assumeClusterAdmin bool, args []string) ([]string, bool) {
	parsedArgs := substituteAliases(args, context, namespace)

	output := make([]string, 0)
	if context != "" {
		output = append(output, "--context", context)
	}

	if namespace != "" {
		output = append(output, "-n", namespace)
	}

	var edit bool
	if isEditableCmd(args[0]) {
		edit = true
	}

	output = append(output, parsedArgs...)

	if allNamespaces {
		output = append(output, "--all-namespaces")
	}

	if assumeClusterAdmin {
		output = append(output, "--as=compute:cluster-admin")
	}

	return output, edit
}

func BuildK9sArgs(context string, namespace string, allNamespaces bool, args []string) ([]string, error) {
	parsedArgs := substituteAliases(args, context, namespace)

	output := make([]string, 0)
	if context != "" {
		output = append(output, "--context", context)
	}

	if namespace != "" {
		output = append(output, "-n", namespace)
	}

	var trailArg string
	if len(parsedArgs) == 2 && parsedArgs[0] == "get" {
		trailArg = parsedArgs[1]
		output = append(output, "-c", trailArg)
	} else if len(parsedArgs) == 1 {
		trailArg = parsedArgs[0]
		output = append(output, "-c", trailArg)
	} else if len(parsedArgs) == 0 {
		// do nothing
	} else {
		return nil, errors.New(fmt.Sprintf("too many arguments provided to k9s: %s", strings.Join(args, " ")))
	}

	if allNamespaces {
		output = append(output, "--all-namespaces")
	}

	return output, nil
}
