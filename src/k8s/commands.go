package k8s

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bitfield/script"
	mdexec "github.com/mdcli/cmd"
	"github.com/urfave/cli/v2"
)

const (
	kubectl = "kubectl"
	k9s     = "k9s"
)

const k8sUsage = `Custom wrapper for k8s commands.`

var BaseK8sFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "context",
		Aliases: []string{"c", "ctx"},
		Value:   "",
		Usage:   "`CONTEXT` from kubeconfig to use",
	},
	&cli.StringFlag{
		Name:    "namespace",
		Aliases: []string{"n", "ns"},
		Value:   "",
		Usage:   "`NAMESPACE` to use",
	},
	&cli.BoolFlag{
		Name:    "strict",
		Aliases: []string{"s"},
		Value:   true,
		Usage:   "Require explicit namespace and context",
	},
	&cli.BoolFlag{
		Name:    "interactive",
		Aliases: []string{"i"},
		Value:   true,
		Usage:   "Enable interactive mode to select context and namespace if not provided",
	},
	&cli.BoolFlag{
		Name:    "dryrun",
		Aliases: []string{"d"},
		Value:   false,
		Usage:   "Enable dry-run mode to show the command to run",
	},
	&cli.BoolFlag{
		Name:    "all-namespaces",
		Aliases: []string{"A"},
		Value:   false,
		Usage:   "Run command across all namespaces",
	},
}

var BaseKctlFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "yes",
		Aliases: []string{"y"},
		Value:   false,
		Usage:   "Automatic yes to confirmation prompts for edit commands",
	},
	&cli.BoolFlag{
		Name:    "assume-cluster-admin",
		Aliases: []string{"cluster-admin"},
		Value:   false,
		Usage:   "Assume cluster-admin role for port-forward",
	},
}

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "kubernetes",
		Aliases: []string{"k", "kube", "k8s"},
		Usage:   k8sUsage,
		Subcommands: []*cli.Command{
			kubectlCommand(),
			tidbKubectlCommand(),
			k9sCommand(),
			tidbK9sCommand(),
			tidbMysqlCommand(),
			tidbDmctlCommand(),
	 	},
	}
}

func kubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "kubectl",
		Aliases: []string{"k", "kc", "kctl"},
		Usage:   "Custom kubectl wrapper",
		Flags: BaseK8sFlags,
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			dryRun := cCtx.Bool("dryrun")
			allNamespaces := cCtx.Bool("all-namespaces")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")

			var err error
			context, err = parseContext(context, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			args, confirm := BuildKubectlArgs(context, namespace, allNamespaces, assumeClusterAdmin, cCtx.Args().Slice())
			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", kubectl, strings.Join(args, " ")))
				return nil
			} else if confirm {
				fmt.Println(fmt.Sprintf("%s %s", kubectl, strings.Join(args, " ")))
				res := getConfirmation("Do you want to execute the above command?")
				if !res {
					fmt.Println("Command canceled")
					return errors.New("Command canceled")
				}
			}

			return mdexec.RunCommand(kubectl, args...)
		},
	}
}

func tidbKubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidbkubectl",
		Aliases: []string{"tk", "tkc", "tkctl"},
		Usage:   "Custom kubectl wrapper for TiDB",
		Flags:  append(BaseK8sFlags, BaseKctlFlags...),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			dryRun := cCtx.Bool("dryrun")
			allNamespaces := cCtx.Bool("all-namespaces")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")

			var err error
			context, err = parseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			args, confirm := BuildKubectlArgs(context, namespace, allNamespaces, assumeClusterAdmin, cCtx.Args().Slice())

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", kubectl, strings.Join(args, " ")))
				return nil
			} else if confirm {
				fmt.Println(fmt.Sprintf("%s %s", kubectl, strings.Join(args, " ")))
				res := getConfirmation("Do you want to execute the above command?")
				if !res {
					fmt.Println("Command canceled")
					return errors.New("Command canceled")
				}
			}

			return mdexec.RunCommand(kubectl, args...)
		},
	}
}

func k9sCommand() *cli.Command {
	return &cli.Command{
		Name:    "k9s",
		Usage:   "Custom k9s wrapper",
		Flags: BaseK8sFlags,
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			dryRun := cCtx.Bool("dryrun")
			allNamespaces := cCtx.Bool("all-namespaces")

			var err error
			context, err = parseContext(context, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			args, err := BuildK9sArgs(context, namespace, allNamespaces, cCtx.Args().Slice())
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", k9s, strings.Join(args, " ")))
				return nil
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", k9s, strings.Join(args, " "))).Stdout()
			return err
		},
	}
}

func tidbK9sCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidbk9s",
		Aliases: []string{"tk9s"},
		Usage:   "Custom k9s wrapper for TiDB",
		Flags: BaseK8sFlags,
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			dryRun := cCtx.Bool("dryrun")
			allNamespaces := cCtx.Bool("all-namespaces")

			var err error
			context, err = parseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			args, err := BuildK9sArgs(context, namespace, allNamespaces, cCtx.Args().Slice())
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", k9s, strings.Join(args, " ")))
				return nil
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", k9s, strings.Join(args, " "))).Stdout()
			return err
		},
	}
}

func tidbMysqlCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidbmysql",
		Aliases: []string{"mysql", "tmysql"},
		Usage:   "Connect to a TiDB via mysql client",
		Flags: append(append(BaseK8sFlags,
			BaseKctlFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
			&cli.StringFlag{
				Name:    "pod-name",
				Usage:   "Pod name to connect to. Defaults to %tc-%az-%pod.",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"P"},
				Value:   -1,
				Usage:   "Port to forward to. Defaults to random 40xx.",
			},
		),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			podName := cCtx.String("pod-name")
			port := cCtx.Int("port")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")

			if port == -1 {
				port = rand.Intn(100) + 4000
			}

			var err error
			context, err = parseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, _, err = parseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			args := make([]string, 0)
			args = append(args, "kubectl", "--context", context, "--namespace", namespace, "get", "secret", "tidb-secret", "-o", "json", "|", "jq", "-r", "'.data.root'", "|", "base64", "-d", "|", "tr", "-d", "'\\n'")
			rootPass, err := mdexec.CaptureCommand("bash", "-c", strings.Join(args, " "))
			if err != nil {
				return err
			}

			if podName == "" {
				podName = "%c-%z-tidb"
			}
			podName = fmt.Sprintf("%s-%d", podName,  pod)
			portForwardCmd, _ := BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"port-forward", podName, fmt.Sprintf("%d:4000", port)})
			cmd := exec.Command("kubectl", portForwardCmd...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			defer func() {
				fmt.Println("Stopping port-forward")
				if cmd.Process != nil {
					if err := cmd.Process.Kill(); err != nil {
						fmt.Println("Error stopping port-forward:", err)
					}
				}
			}()

			portForwardErr := make(chan error)
			go func() {
				fmt.Println(fmt.Sprintf("Starting port-forward from %s:4000 to %d", podName, port))
				if err = cmd.Run(); err != nil {
					fmt.Println(err)
					portForwardErr <- err
				}
			}()

			time.Sleep(2 * time.Second)
			select {
			case err = <-portForwardErr:
				return err
			default:
			}

			mysqlCmd := exec.Command("mysql", "-h", "127.0.0.1", "-P", fmt.Sprintf("%d", port), "-u", "root", "-p" + rootPass, "--prompt=tidb> ")
			mysqlCmd.Stdin = os.Stdin
			mysqlCmd.Stdout = os.Stdout
			mysqlCmd.Stderr = os.Stderr
			if err = mysqlCmd.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

func tidbDmctlCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidbdmctl",
		Aliases: []string{"dmctl", "tdmctl"},
		Usage:   "Connect to dmctl on a TiDB cluster",
		Flags: append(append(BaseK8sFlags,
			BaseKctlFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
		),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")

			var err error
			context, err = parseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, _, err = parseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			clusterName := strings.TrimPrefix(namespace, "tidb-")
			podName := fmt.Sprintf("%s-dm-master-%d", clusterName, pod)

			execArgs, _ := BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"exec", "-it", podName, "-c", "dm-master", "--", "bin/sh", "-c", `./dmctl --master-addr https://127.0.0.1:8261 --ssl-cert /var/lib/dm-master-tls/tls.crt --ssl-key /var/lib/dm-master-tls/tls.key --ssl-ca /var/lib/dm-master-tls/ca.crt`})
			fmt.Println(execArgs)
			return mdexec.RunCommand("kubectl", execArgs...)
		},
	}
}
