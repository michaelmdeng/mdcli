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

// Order matters, parse longer aliases w/ common prefix first
var NamespaceVariableAliases = []string{"%ns", "%n"}
var ContextVariableAliases = []string{"%ctx", "%c"}
var TidbClusterVariableAliases = []string{"%tc", "%t"}
var AZVariableAliases = []string{"%z", "%az"}
var AppVariableAliases = []string{"%app", "%ap"}

func generateContextAlias(context string, namespace string) string {
	return context
}

func generateNamespaceAlias(context string, namespace string) string {
	return namespace
}

func generateTidbClusterAlias(context string, namespace string) string {
	return strings.TrimPrefix(namespace, "tidb-")
}

func generateAppAlias(context string, namespace string) string {
	return strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(namespace, "tidb-"), "-test"), "-stg"), "-prod")
}

func generateAZAlias(context string, namespace string) (string, error) {
	// parse zone from context
	// ex. m-tidb-test-<zone>-ea1-us
	pattern := regexp.MustCompile(`m-tidb-[a-z]+-([a-z])-ea1-us`)
	matches := pattern.FindStringSubmatch(context)
	if len(matches) >= 2 {
		zone := matches[1]
		if zone == "c" {
			zone = "e"
		}
		return fmt.Sprintf("us-east-1%s", zone), nil
	} else {
		return "", errors.New("Could not generate AZ alias")
	}
}

func substituteAliases(args []string, context string, namespace string) []string {
	for i, arg := range args {
		for _, alias := range NamespaceVariableAliases {
			arg = strings.ReplaceAll(arg, alias, generateNamespaceAlias(context, namespace))
		}
		for _, alias := range ContextVariableAliases {
			arg = strings.ReplaceAll(arg, alias, generateContextAlias(context, namespace))
		}
		for _, alias := range TidbClusterVariableAliases {
			arg = strings.ReplaceAll(arg, alias, generateTidbClusterAlias(context, namespace))
		}
		for _, alias := range AppVariableAliases {
			arg = strings.ReplaceAll(arg, alias, generateAppAlias(context, namespace))
		}
		for _, alias := range AZVariableAliases {
			az, err := generateAZAlias(context, namespace)
			if err == nil {
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

	kubectlCmd := args[0]

	var confirm bool
	if isConfirmableCmd(kubectlCmd) {
		confirm = true
	}

	var resourceModified bool
	var resourceType, resourceName string
	if isResourceModifiableCmd(kubectlCmd) {
		resourceType = args[1]
		if isModifiableResource(resourceType) {
			resourceName = args[2]
			resourceModified = true
		}
	}

	var last int
	for i, arg := range parsedArgs {
		if resourceModified {
			last = i + 1
			if i == 1 {
				output = append(output, fmt.Sprintf("%s/%s", resourceType, resourceName))
				continue
			} else if i == 2 {
				continue
			}
		}

		if arg != "--" {
			output = append(output, arg)
			last = i + 1
		} else {
			last = i
			break
		}
	}

	if allNamespaces {
		output = append(output, "--all-namespaces")
	}

	if assumeClusterAdmin && isEditCmd(kubectlCmd) {
		output = append(output, "--as=compute:cluster-admin")
	}

	for idx := last; idx < len(parsedArgs); idx++ {
			output = append(output, parsedArgs[idx])
	}

	return output, confirm
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
