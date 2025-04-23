package tidb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mdexec "github.com/michaelmdeng/mdcli/internal/cmd"
	mdk8s "github.com/michaelmdeng/mdcli/k8s"
	"github.com/urfave/cli/v2"
)

func BaseTikvCommand() *cli.Command {
	return &cli.Command{
		Name:    "tikv",
		Aliases: []string{"kv"},
		Usage:   `Commands for handling TiKVs on K8s`,
		Subcommands: []*cli.Command{
			tikvDeleteCommand(),
			tikvGetCommand(),
			tikvStoreCommand(),
		},
	}
}

type getTikvOutput struct {
	name       string
	storeId    int
	instanceId string
	dataVol    string
	walVol     string
	raftVol    string
}

func tikvGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Fetch tikv info",
		Flags: append(mdk8s.BaseK8sFlags),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			debug := cCtx.Bool("debug") && !mdexec.IsPipe()
			allNamespaces := cCtx.Bool("all-namespaces")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			tikvName := cCtx.Args().Get(0)
			if tikvName == "" {
				return cli.Exit("tikv name is required", 1)
			}
			clusterName := strings.TrimPrefix(namespace, "tidb-")
			tikvName = strings.TrimPrefix(tikvName, clusterName+"-")
			tikvName = strings.TrimPrefix(tikvName, "tikv-")
			tikvNum := tikvName
			tikvName = fmt.Sprintf("%s-tikv-%s", clusterName, tikvName)

			tikvOutput := map[string]any{}
			tikvOutput["name"] = tikvName

			builder := NewTidbKubeBuilder()
			args, _ := builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "tc", clusterName, "-o", "jsonpath='{.status.tikv.stores}'"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err := mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			output = output[1 : len(output)-1]

			var tikvStores map[string]any
			err = json.Unmarshal([]byte(output), &tikvStores)
			if err != nil {
				return err
			}

			var storeId int
			for _, store := range tikvStores {
				store := store.(map[string]any)
				if strings.HasPrefix(store["ip"].(string), tikvName) {
					storeId, err = strconv.Atoi(store["id"].(string))
					if err != nil {
						return err
					}
				}
			}
			tikvOutput["storeId"] = storeId

			dataPvc := fmt.Sprintf("tikv-%s-tikv-%v", clusterName, tikvNum)
			walPvc := fmt.Sprintf("tikv-wal-%s-tikv-%v", clusterName, tikvNum)
			raftPvc := fmt.Sprintf("tikv-raft-%s-tikv-%v", clusterName, tikvNum)

			args, _ = builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "pvc", dataPvc, walPvc, raftPvc, "-o", "jsonpath='{.items[*].spec.volumeName}'"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err = mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			pvs := strings.Split(output[1:len(output)-1], " ")
			dataPv := pvs[0]
			walPv := pvs[1]
			raftPv := pvs[2]

			args, _ = builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "pv", dataPv, walPv, raftPv, "-o", `jsonpath='{range .items[*]}{"{\""}{.metadata.name}{"\":\""}{.spec.csi.volumeHandle}{"\"}\n"}{end}'`})
			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err = mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			// Trim enclosing quotes and whitespace
			output = strings.TrimSpace(output[1 : len(output)-1])
			var pvHandle map[string]string
			for _, line := range strings.Split(output, "\n") {
				err = json.Unmarshal([]byte(line), &pvHandle)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
				for pv, handle := range pvHandle {
					if pv == dataPv {
						tikvOutput["dataVol"] = handle
					} else if pv == walPv {
						tikvOutput["walVol"] = handle
					} else if pv == raftPv {
						tikvOutput["raftVol"] = handle
					}
				}
			}

			args, _ = builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "pod", tikvName, "-o", "jsonpath='{.spec.nodeName}'"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err = mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			nodeName := output[1 : len(output)-1]

			args, _ = builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "node", nodeName, "-o", "jsonpath='{.metadata.labels.node\\.airbnb\\.com/instance-id}'"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err = mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			instanceId := output[1 : len(output)-1]
			tikvOutput["instanceId"] = instanceId

			out, err := json.Marshal(tikvOutput)
			if err != nil {
				return err
			}
			fmt.Println(string(out))

			return nil
		},
	}
}

func tikvStoreCommand() *cli.Command {
	return &cli.Command{
		Name:  "store",
		Usage: "Fetch tikv store info",
		Flags: append(mdk8s.BaseK8sFlags),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			debug := cCtx.Bool("debug") && !mdexec.IsPipe()
			allNamespaces := cCtx.Bool("all-namespaces")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			tikvName := cCtx.Args().Get(0)
			clusterName := strings.TrimPrefix(namespace, "tidb-")
			tikvName = strings.TrimPrefix(tikvName, clusterName+"-")
			tikvName = strings.TrimPrefix(tikvName, "tikv-")
			tikvName = fmt.Sprintf("%s-tikv-%s", clusterName, tikvName)

			builder := NewTidbKubeBuilder()
			args, _ := builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"get", "tc", clusterName, "-o", "jsonpath='{.status.tikv.stores}'"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			output, err := mdexec.CaptureCommand(mdk8s.Kubectl, args...)
			if err != nil {
				return err
			}
			output = output[1 : len(output)-1]

			var tikvStores map[string]any
			err = json.Unmarshal([]byte(output), &tikvStores)
			if err != nil {
				return err
			}

			var storeId int
			for _, store := range tikvStores {
				store := store.(map[string]any)
				if strings.HasPrefix(store["ip"].(string), tikvName) {
					storeId, err = strconv.Atoi(store["id"].(string))
					if err != nil {
						return err
					}

					fmt.Println(storeId)
				}
			}

			return nil
		},
	}
}

func tikvDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete tikv store pod safely",
		Flags: append(mdk8s.BaseK8sFlags),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			debug := cCtx.Bool("debug") && !mdexec.IsPipe()
			allNamespaces := cCtx.Bool("all-namespaces")

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			tikvName := cCtx.Args().Get(0)
			clusterName := strings.TrimPrefix(namespace, "tidb-")
			tikvName = strings.TrimPrefix(tikvName, clusterName+"-")
			tikvName = strings.TrimPrefix(tikvName, "tikv-")
			tikvName = fmt.Sprintf("%s-tikv-%s", clusterName, tikvName)

			builder := NewTidbKubeBuilder()
			args, _ := builder.BuildKubectlArgs(context, namespace, allNamespaces, false, []string{"annotate", "pod", tikvName, "tidb.pingcap.com/evict-leader=delete-pod"})

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			return mdexec.RunCommand(mdk8s.Kubectl, args...)
		},
	}
}
