package tidb

import (
	"strings"
)

func isTestTidbContext(context string) bool {
	return strings.Contains(context, "test")
}

func isStgTidbContext(context string) bool {
	return strings.Contains(context, "stg")
}

func isProdTidbContext(context string) bool {
	return strings.Contains(context, "prod")
}
