package tidb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	mdk8s "github.com/michaelmdeng/mdcli/k8s"
)

var contextAliases = map[string][]string{
	"m-tidb-prod-a-ea1-us": {
		"proda",
		"prod",
	},
	"m-tidb-prod-b-ea1-us": {
		"prodb",
	},
	"m-tidb-prod-c-ea1-us": {
		"prodc",
		"prode",
	},
	"m-tidb-stg-a-ea1-us": {
		"stga",
		"stg",
	},
	"m-tidb-stg-b-ea1-us": {
		"stgb",
	},
	"m-tidb-stg-c-ea1-us": {
		"stgc",
		"stge",
	},
	"m-tidb-test-a-ea1-us": {
		"testa",
		"test",
	},
	"m-tidb-test-b-ea1-us": {
		"testb",
	},
	"m-tidb-test-c-ea1-us": {
		"testc",
		"teste",
	},
}

var ContextsByAlias = make(map[string]string)

func init() {
	for context, contextAliases := range contextAliases {
		for _, alias := range contextAliases {
			ContextsByAlias[alias] = context
		}
	}

	for namespace, namespaceAliases := range prodNamespaceAliases {
		for _, alias := range namespaceAliases {
			ProdNamespacesByAlias[alias] = namespace
		}
	}
	for namespace, namespaceAliases := range stgNamespaceAliases {
		for _, alias := range namespaceAliases {
			StgNamespacesByAlias[alias] = namespace
		}
	}
	for namespace, namespaceAliases := range testNamespaceAliases {
		for _, alias := range namespaceAliases {
			TestNamespacesByAlias[alias] = namespace
		}
	}
}

var prodNamespaceAliases = map[string][]string{
	"tidb-mussel-prod-ml-dr": {
		"mldr",
	},
	"tidb-mussel-prod-ml": {
		"ml",
	},
	"tidb-mussel-prod-ml-dr-1": {
		"mldr1",
	},
	"tidb-mussel-prod-ml-dr-2": {
		"mldr2",
	},
	"tidb-mussel-prod-ml-dr-3": {
		"mldr3",
	},
	"tidb-mussel-prod-dr": {
		"mergedr",
	},
	"tidb-mussel-prod": {
		"merge",
	},
	"tidb-mussel-prod-replace": {
		"replace",
	},
	"tidb-mussel-prod-replace-dr": {
		"replacedr",
	},
	"tidb-restore-operator-prod": {
		"restore",
		"restoreoperator",
	},
	"tidb-migration-operator-prod": {
		"migration",
		"migrationoperator",
	},
}

var stgNamespaceAliases = map[string][]string{
	"tidb-mussel-stg-replace": {
		"stgreplace",
		"replace",
	},
	"tidb-mussel-stg-replace-dr": {
		"replacedr",
	},
	"tidb-mussel-stg": {
		"stg",
		"merge",
	},
	"tidb-mussel-stg-dr": {
		"mergedr",
	},
	"tidb-restore-operator-stg": {
		"restore",
		"restoreoperator",
	},
	"tidb-migration-operator-stg": {
		"migration",
		"migrationoperator",
	},
	"tidb-uds-full-shadow-stg": {
		"udsfullshadow",
		"fullshadow",
	},
	"tidb-uds-full-shadow-stg-br-1": {
		"udsfullshadowbr1",
		"fullshadowbr1",
	},
	"tidb-uds-full-shadow-stg-br-2": {
		"udsfullshadowbr2",
		"fullshadowbr2",
	},
	"tidb-uds-shadow-stg": {
		"udsshadow",
		"shadow",
	},
	"tidb-ingestion-staging": {
		"ingestion",
	},
	"tidb-ingestion-staging-dr": {
		"ingestiondr",
	},
}

var testNamespaceAliases = map[string][]string{
	"tidb-mussel-stag-replace": {
		"stagreplace",
		"loadtest",
	},
	"tidb-loadtest-br-1": {
		"loadtestbr1",
		"loadtest1",
	},
	"tidb-loadtest-br-2": {
		"loadtestbr2",
		"loadtest2",
	},
	"tidb-func-test": {
		"func",
	},
	"tidb-func-test-1": {
		"func1",
	},
	"tidb-func-test-2": {
		"func2",
	},
	"tidb-func-test-3": {
		"func3",
	},
	"tidb-rp-mussel-func-test": {
		"rp",
		"rpfunc",
		"rpmussfunc",
		"rpmusselfunc",
	},
	"tidb-rp-mussel-func-test-br-1": {
		"rpfuncbr1",
		"rpbr1",
	},
	"tidb-rp-mussel-func-test-br-2": {
		"rprestore2",
		"rpbr2",
	},
	"tidb-rp-func-test-2": {
		"rp2",
		"rpfunc2",
	},
	"tidb-rp-func-test-2-br-1": {
		"rpfunc2br1",
		"rp2br1",
	},
	"tidb-rp-func-test-2-br-2": {
		"rpfunc2br2",
		"rp2br2",
	},
	"tidb-rp-uds-func-test": {
		"rpudsfunc",
		"udsfunc",
	},
	"tidb-rp-uds-func-test-br-1": {
		"rpudsfuncbr1",
		"udsfuncbr1",
	},
	"tidb-rp-uds-func-test-br-2": {
		"rpudsfuncbr2",
		"udsfuncbr2",
	},
	"tidb-rp-uds-load-test": {
		"rpudsload",
		"udsload",
	},
	"tidb-restore-operator-test": {
		"restore",
		"restoreoperator",
	},
	"tidb-migration-operator-test": {
		"migration",
		"migrationoperator",
		"migrationtest",
		"migrationoperatortest",
	},
	"tidb-migration-operator-dev": {
		"migrationdev",
		"migrationoperatordev",
	},
	"tidb-test-single-cell": {
		"singlecell",
	},
	"tidb-toolbox-test": {
		"toolbox",
	},
	"tidb-dev-mdeng-test": {
		"mdeng",
	},
	"tidb-release-production": {
		"release",
	},
	"tidb-mussel-test-hightouch-1": {
		"musselhightouch",
	},
	"tidb-mussel-test-ht1-v75": {
		"musselhightouch75",
	},
}

var ProdNamespacesByAlias = make(map[string]string)
var StgNamespacesByAlias = make(map[string]string)
var TestNamespacesByAlias = make(map[string]string)

func inferContext(kubecontext string) (string, bool) {
	contextAlias := strings.ReplaceAll(strings.ToLower(kubecontext), "-", "")

	if context, ok := ContextsByAlias[contextAlias]; ok {
		return context, true
	}

	return kubecontext, false
}

func inferNamespace(context string, namespace string) (string, bool) {
	namespaceAlias := strings.ReplaceAll(strings.ToLower(namespace), "-", "")

	var env string
	if strings.Contains(context, "prod") {
		env = "prod"
	} else if strings.Contains(context, "stg") {
		env = "stg"
	} else if strings.Contains(context, "test") {
		env = "test"
	} else if strings.Contains(context, "dev") {
		env = "test"
	} else {
		return namespace, false
	}

	if env == "prod" {
		if namespace, ok := ProdNamespacesByAlias[namespaceAlias]; ok {
			return namespace, true
		}
	} else if env == "stg" {
		if namespace, ok := StgNamespacesByAlias[namespaceAlias]; ok {
			return namespace, true
		}
	} else if env == "test" {
		if namespace, ok := TestNamespacesByAlias[namespaceAlias]; ok {
			return namespace, true
		}
	}

	return namespace, false
}

func ParseContext(context string, interactive bool, pattern string, strict bool) (string, error) {
	if context != "" {
		context, _ = inferContext(context)
		return context, nil
	}

	if interactive && context == "" {
		var err error
		context, err = mdk8s.GetContextInteractive(pattern)
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
		namespace, _ = inferNamespace(context, namespace)
		return namespace, false, nil
	}

	if interactive && !allNamespaces && namespace == "" {
		var err error
		namespace, err = mdk8s.GetNamespaceInteractive(context, pattern)
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

var (
	TidbClusterSubstitution = mdk8s.Substitution{
		Aliases: []string{
			"tc",
			"t",
		},
		Generate: func(context, namespace string) (string, error) {
			return strings.TrimPrefix(namespace, "tidb-"), nil
		},
	}

	AzSubstitution = mdk8s.Substitution{
		Aliases: []string{
			"az",
			"z",
		},
		Generate: func(context, namespace string) (string, error) {
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
		},
	}

	AppSubstitution = mdk8s.Substitution{
		Aliases: []string{
			"app",
			"ap",
		},
		Generate: func(context, namespace string) (string, error) {
			return strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(namespace, "tidb-"), "-test"), "-stg"), "-prod"), nil
		},
	}
)

func NewTidbKubeBuilder() mdk8s.KubeBuilder {
	substitutions := []mdk8s.Substitution{
		TidbClusterSubstitution,
		AzSubstitution,
		AppSubstitution,
	}
	return mdk8s.NewKubeBuilderWithSubstitutions(substitutions)
}
