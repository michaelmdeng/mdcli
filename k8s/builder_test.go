package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstitute(t *testing.T) {
	defaultBuilder := NewKubeBuilder()

	customSub := Substitution{
		Aliases: []string{"n"},
		Generate: func(context, namespace string) (string, error) {
			return "custom-namespace", nil
		},
	}
	// NewKubeBuilderWithSubstitutions appends the custom substitutions, so the
	// default substitution for "n" comes first and wins.
	firstWinsBuilder := NewKubeBuilderWithSubstitutions([]Substitution{customSub})

	testCases := []struct {
		name      string
		builder   KubeBuilder
		args      []string
		context   string
		namespace string
		expected  []string
	}{
		{
			name:      "Basic alias substitution",
			builder:   defaultBuilder,
			args:      []string{"get", "pods", "-l", "app=%ns-app"},
			context:   "my-context",
			namespace: "my-namespace",
			expected:  []string{"get", "pods", "-l", "app=my-namespace-app"},
		},
		{
			name:      "Multiple substitutions",
			builder:   defaultBuilder,
			args:      []string{"get", "pods", "-l", "app=%ns-app", "--context=%ctx"},
			context:   "my-context",
			namespace: "my-namespace",
			expected:  []string{"get", "pods", "-l", "app=my-namespace-app", "--context=my-context"},
		},
		{
			name:      "No matching substitutions",
			builder:   defaultBuilder,
			args:      []string{"get", "pods"},
			context:   "my-context",
			namespace: "my-namespace",
			expected:  []string{"get", "pods"},
		},
		{
			name:      "Multiple substitutions in same argument",
			builder:   defaultBuilder,
			args:      []string{"get", "pods", "--name=%ctx-%ns"},
			context:   "my-context",
			namespace: "my-namespace",
			expected:  []string{"get", "pods", "--name=my-context-my-namespace"},
		},
		{
			name:      "First substitution wins",
			builder:   firstWinsBuilder,
			args:      []string{"get", "pods", "-l", "app=%n-app"},
			context:   "my-context",
			namespace: "my-namespace",
			expected:  []string{"get", "pods", "-l", "app=my-namespace-app"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.builder.Substitute(tc.args, tc.context, tc.namespace)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBuildKubectlArgs(t *testing.T) {
	builder := NewKubeBuilder()

	testCases := []struct {
		name               string
		context            string
		namespace          string
		allNamespaces      bool
		assumeClusterAdmin bool
		args               []string
		expectedArgs       []string
		expectedConfirm    bool
	}{
		{
			name:            "Basic command",
			args:            []string{"get", "pods"},
			expectedArgs:    []string{"get", "pods"},
			expectedConfirm: false,
		},
		{
			name:         "With context and namespace",
			context:      "my-context",
			namespace:    "my-namespace",
			args:         []string{"get", "pods"},
			expectedArgs: []string{"--context", "my-context", "--namespace", "my-namespace", "get", "pods"},
		},
		{
			name:          "All namespaces",
			allNamespaces: true,
			args:          []string{"get", "pods"},
			expectedArgs:  []string{"get", "pods", "--all-namespaces"},
		},
		{
			name:               "Assume cluster admin",
			assumeClusterAdmin: true,
			args:               []string{"edit", "deployment", "my-deployment"},
			expectedArgs:       []string{"edit", "deployment", "my-deployment", "--as=compute:cluster-admin"},
			expectedConfirm:    false,
		},
		{
			name:            "Confirmation for mutating commands",
			args:            []string{"delete", "pod", "my-pod"},
			expectedArgs:    []string{"delete", "pod", "my-pod"},
			expectedConfirm: true,
		},
		{
			name:            "Rewrites resource",
			args:            []string{"logs", "job", "my-job"},
			expectedArgs:    []string{"logs", "job/my-job"},
			expectedConfirm: false,
		},
		{
			name:            "Doesn't modify qualified / resource",
			args:            []string{"logs", "job/my-job"},
			expectedArgs:    []string{"logs", "job/my-job"},
			expectedConfirm: false,
		},
		{
			name:            "Empty exec",
			args:            []string{"exec", "my-pod"},
			expectedArgs:    []string{"exec", "my-pod", "-it", "--", "bash"},
			expectedConfirm: false,
		},
		{
			name:            "Empty exec for rewritable resource",
			args:            []string{"exec", "job", "my-job"},
			expectedArgs:    []string{"exec", "job/my-job", "-it", "--", "bash"},
			expectedConfirm: false,
		},
		{
			name:         "Empty arguments",
			context:      "my-context",
			namespace:    "my-namespace",
			args:         []string{},
			expectedArgs: []string{"--context", "my-context", "--namespace", "my-namespace"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, confirm := builder.BuildKubectlArgs(tc.context, tc.namespace, tc.allNamespaces, tc.assumeClusterAdmin, tc.args)
			assert.Equal(t, tc.expectedArgs, result)
			assert.Equal(t, tc.expectedConfirm, confirm)
		})
	}
}

func TestBuildK9sArgs(t *testing.T) {
	builder := NewKubeBuilder()

	testCases := []struct {
		name          string
		context       string
		namespace     string
		allNamespaces bool
		args          []string
		expected      []string
		expectError   bool
	}{
		{
			name:          "Basic command",
			context:       "",
			namespace:     "",
			allNamespaces: false,
			args:          []string{"get", "pods"},
			expected:      []string{"-c", "pods"},
			expectError:   false,
		},
		{
			name:          "With context and namespace",
			context:       "my-context",
			namespace:     "my-namespace",
			allNamespaces: false,
			args:          []string{"get", "pods"},
			expected:      []string{"--context", "my-context", "-n", "my-namespace", "-c", "pods"},
			expectError:   false,
		},
		{
			name:          "All namespaces",
			context:       "",
			namespace:     "",
			allNamespaces: true,
			args:          []string{"get", "pods"},
			expected:      []string{"-c", "pods", "--all-namespaces"},
			expectError:   false,
		},
		{
			name:          "Empty args",
			context:       "my-context",
			namespace:     "my-namespace",
			allNamespaces: false,
			args:          []string{},
			expected:      []string{"--context", "my-context", "-n", "my-namespace", "-c", "pods"},
			expectError:   false,
		},
		{
			name:          "Too many args",
			context:       "",
			namespace:     "",
			allNamespaces: false,
			args:          []string{"get", "pods", "extra"},
			expected:      nil,
			expectError:   true,
		},
		{
			name:          "Get with resource",
			context:       "",
			namespace:     "",
			allNamespaces: false,
			args:          []string{"get", "deployments"},
			expected:      []string{"-c", "deployments"},
			expectError:   false,
		},
		{
			name:          "Resource only",
			context:       "",
			namespace:     "",
			allNamespaces: false,
			args:          []string{"deployments"},
			expected:      []string{"-c", "deployments"},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := builder.BuildK9sArgs(tc.context, tc.namespace, tc.allNamespaces, tc.args)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
