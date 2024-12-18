package tidb

import (
	"strings"

	"github.com/fatih/color"
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

func contextColor(context string) color.Attribute {
	if isTestTidbContext(context) {
		return color.FgGreen
	}
	return 0
}
