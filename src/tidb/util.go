package tidb

import (
	"strings"
)

func isTestTidbContext(context string) bool {
	return strings.Contains(context, "test")
}
