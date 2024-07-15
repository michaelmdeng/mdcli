package k8s

import (
	"fmt"
	"errors"
	"strings"
)

type KubeBuilder struct {
	Substitutions []Substitution
}

type Substitution struct {
	Aliases []string
	Generate func(context, namespace string) (string, error)
}

var (
	ContextSubstitution = Substitution{
		Aliases: []string{
			"ctx", "c",
		},
		Generate: func(context, namespace string) (string, error) {
			return context, nil
		},
	}

	NamespaceSubstitution = Substitution{
		Aliases: []string{
			"ns", "n",
		},
		Generate: func(context, namespace string) (string, error) {
			return namespace, nil
		},
	}
)

func NewKubeBuilder() KubeBuilder {
	baseSubstitutions := []Substitution{
		ContextSubstitution,
		NamespaceSubstitution,
	}
	return KubeBuilder{Substitutions: baseSubstitutions}
}

func NewKubeBuilderWithSubstitutions(substitutions []Substitution) KubeBuilder {
	baseSubstitutions := []Substitution{
		ContextSubstitution,
		NamespaceSubstitution,
	}
	substitutions = append(baseSubstitutions, substitutions...)
	return KubeBuilder{Substitutions: substitutions}
}

func (b *KubeBuilder) substitute(args []string, context, namespace string) []string {
	for i, arg := range args {
		for _, sub := range b.Substitutions {
			for _, alias := range sub.Aliases {
				substition, err := sub.Generate(context, namespace)
				if err != nil {
					continue
				}

				arg = strings.ReplaceAll(arg, fmt.Sprintf("%%%s", alias), substition)
			}
		}

		args[i] = arg
	}
	return args
}

func (b *KubeBuilder) BuildKubectlArgs(context string, namespace string, allNamespaces bool, assumeClusterAdmin bool, args []string) ([]string, bool) {
	parsedArgs := b.substitute(args, context, namespace)

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

func (b *KubeBuilder) BuildK9sArgs(context string, namespace string, allNamespaces bool, args []string) ([]string, error) {
	parsedArgs := b.substitute(args, context, namespace)

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
