package tidb

import (
	"strings"

	mdexec "github.com/michaelmdeng/mdcli/cmd"
)

func getTidbSecret(context, namespace string) (string, error) {
	args := make([]string, 0)
	args = append(args, "kubectl", "--context", context, "--namespace", namespace, "get", "secret", "tidb-secret", "-o", "json", "|", "jq", "-r", "'.data.root'", "|", "base64", "-d", "|", "tr", "-d", "'\\n'")
	rootPass, err := mdexec.CaptureCommand("bash", "-c", strings.Join(args, " "))
	if err != nil {
		return "", err
	}

	return rootPass, nil
}
