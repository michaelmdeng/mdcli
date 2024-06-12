package k8s

import (
	"strings"
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
	},
	"tidb-mussel-stg": {
		"stg",
		"merge",
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
	"tidb-uds-shadow-stg": {
		"udsshadow",
		"shadow",
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
	},
	"tidb-rp-mussel-func-test-br-1": {
		"rprestore1",
		"rpbr1",
	},
	"tidb-rp-mussel-func-test-br-2": {
		"rprestore2",
		"rpbr2",
	},
	"tidb-rp-mussel-func-test-br-3": {
		"rprestore3",
		"rpbr3",
	},
	"tidb-rp-func-test-2": {
		"rp2",
		"rpfunc2",
	},
	"tidb-rp-func-test-2-restore-1": {
		"rp2restore1",
		"rp2br1",
	},
	"tidb-rp-func-test-2-restore-2": {
		"rp2restore2",
		"rp2br2",
	},
	"tidb-rp-func-test-2-restore-3": {
		"rp2restore3",
		"rp2br3",
	},
	"tidb-rp-uds-func-test": {
		"rpudsfunc",
		"udsfunc",
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
}

var ProdNamespacesByAlias = make(map[string]string)
var StgNamespacesByAlias = make(map[string]string)
var TestNamespacesByAlias = make(map[string]string)

func inferNamespace(context string, namespaceAlias string) (string, bool) {
	namespaceAlias = strings.ReplaceAll(strings.ToLower(namespaceAlias), "-", "")

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
		return namespaceAlias, false
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

	return namespaceAlias, false
}
