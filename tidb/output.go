package tidb

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func colorDebugPrintfln(context string, format string, args ...interface{}) {
	format = fmt.Sprintf("%s\n", format)
	if isProdTidbContext(context) {
		color.New(color.FgRed).Fprintf(os.Stderr, format, args...)
	} else if isStgTidbContext(context) {
		color.New(color.FgYellow).Fprintf(os.Stderr, format, args...)
	} else if isTestTidbContext(context) {
		color.New(color.FgMagenta).Fprintf(os.Stderr, format, args...)
	} else {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func debugPrintfln(format string, args ...interface{}) {
	format = fmt.Sprintf("%s\n", format)
	fmt.Fprintf(os.Stderr, format, args...)
}

func debugPrintln(args ...interface{}) {
	fmt.Fprintln(os.Stderr,  args...)
}
