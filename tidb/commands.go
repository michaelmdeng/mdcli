package tidb

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bitfield/script"
	mdexec "github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/michaelmdeng/mdcli/internal/config"
	mdk8s "github.com/michaelmdeng/mdcli/k8s"
	"github.com/urfave/cli/v2"
)

var BaseTidbFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "disable-tls",
		Aliases: []string{"tls"},
		Value:   false,
		Usage:   "Whether to disable TLS. Defaults to false.",
	},
}

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidb",
		Aliases: []string{"ti", "db"},
		Usage:   `Commands for managing TiDB on K8s`,
		Subcommands: []*cli.Command{
			tidbSecretCommand(),
			tidbKubectlCommand(),
			tidbK9sCommand(),
			tidbMysqlCommand(),
			tidbDmctlCommand(),
			tidbPdctlCommand(),
			ticdcCommand(),
			BasePdCommand(),
			BaseTikvCommand(),
		},
	}
}

func tidbSecretCommand() *cli.Command {
	return &cli.Command{
		Name:    "secret",
		Aliases: []string{"pass", "password"},
		Usage:   "Fetch tidb root user password",
		Flags:   append(mdk8s.BaseK8sFlags),
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			allNamespaces := cCtx.Bool("all-namespaces")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			rootPass, err := getTidbSecret(context, namespace)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to get tidb secret: %v", err), 1)
			}

			fmt.Print(rootPass)
			return nil
		},
	}
}

func tidbKubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "kubectl",
		Aliases: []string{"kc", "kctl", "tkc", "tkctl"},
		Usage:   "kubectl wrapper for TiDB",
		Flags:   append(mdk8s.BaseK8sFlags, mdk8s.BaseKctlFlags...),
		Action: func(cCtx *cli.Context) error {
			cfg := config.NewConfig()

			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			debug := cCtx.Bool("debug")
			allNamespaces := cCtx.Bool("all-namespaces")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")
			confirmed := cCtx.Bool("yes")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if isTestTidbContext(context) {
				assumeClusterAdmin = assumeClusterAdmin || cfg.EnableClusterAdminForTest
			}
			builder := NewTidbKubeBuilder()
			args, confirm := builder.BuildKubectlArgs(context, namespace, allNamespaces, assumeClusterAdmin, cCtx.Args().Slice())

			needsConfirm := confirm && !confirmed
			if debug || needsConfirm {
				colorDebugPrintfln(context, "%s %s", mdk8s.Kubectl, strings.Join(args, " "))
			}

			if needsConfirm {
				res := mdexec.GetConfirmation("Do you want to execute the above command?")
				if !res {
					return cli.Exit("Command canceled by user", 1)
				}
			}

			// Use the helper function for external command errors
			return mdexec.ExitError(mdexec.RunCommand(mdk8s.Kubectl, args...))
		},
	}
}

func tidbK9sCommand() *cli.Command {
	return &cli.Command{
		Name:    "k9s",
		Aliases: []string{"tk9s"},
		Usage:   "k9s wrapper for TiDB",
		Flags:   mdk8s.BaseK8sFlags,
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			debug := cCtx.Bool("debug")
			allNamespaces := cCtx.Bool("all-namespaces")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			builder := NewTidbKubeBuilder()
			args, err := builder.BuildK9sArgs(context, namespace, allNamespaces, cCtx.Args().Slice())
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to build k9s args: %v", err), 1)
			}

			// Check cellauth login status without printing output
			err = mdexec.RunCommandDiscardOutput("cellauth", "token", "--region", "us-east-1", context)
			if err != nil {
				return cli.Exit(fmt.Sprintf("cellauth check failed: %v", err), 1)
			}

			if debug {
				colorDebugPrintfln(context, "%s %s", mdk8s.K9s, strings.Join(args, " "))
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", mdk8s.K9s, strings.Join(args, " "))).Stdout()
			return mdexec.ExitError(err)
		},
	}
}

func tidbMysqlCommand() *cli.Command {
	return &cli.Command{
		Name:    "mysql",
		Aliases: []string{"tmysql"},
		Usage:   "Connect to a TiDB via mysql client",
		Flags: append(append(mdk8s.BaseK8sFlags,
			mdk8s.BaseKctlFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
			&cli.StringFlag{
				Name:  "pod-name",
				Usage: "Pod name to connect to. Defaults to %tc-%az-%pod.",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"P"},
				Value:   -1,
				Usage:   "Port to forward to. Defaults to random 40xx.",
			},
		),
		Action: func(cCtx *cli.Context) error {
			cfg := config.NewConfig()

			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			podName := cCtx.String("pod-name")
			port := cCtx.Int("port")
			debug := cCtx.Bool("debug")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")

			if port == -1 {
				port = rand.Intn(100) + 4000
			}

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, _, err = ParseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			rootPass, err := getTidbSecret(context, namespace)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to get tidb secret: %v", err), 1)
			}

			if isTestTidbContext(context) {
				assumeClusterAdmin = assumeClusterAdmin || cfg.EnableClusterAdminForTest
			}
			if podName == "" {
				podName = "%tc-%az-tidb"
			}
			podName = fmt.Sprintf("%s-%d", podName, pod)
			builder := NewTidbKubeBuilder()
			portForwardCmd, _ := builder.BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"port-forward", podName, fmt.Sprintf("%d:4000", port)})
			cmd := exec.Command("kubectl", portForwardCmd...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			defer func() {
				debugPrintln("Stopping port-forward...")
				if cmd.Process != nil {
					if err := cmd.Process.Kill(); err != nil {
						debugPrintln("Error stopping port-forward:", err)
					}
				}
			}()

			if debug {
				colorDebugPrintfln(context, "%s %s", "kubectl", strings.Join(portForwardCmd, " "))
			}

			portForwardErr := make(chan error)
			go func() {
				debugPrintfln("Starting port-forward from %s:4000 to %d", podName, port)
				if err = cmd.Run(); err != nil {
					debugPrintln(err)
					portForwardErr <- err
				}
			}()

			// Wait a bit for port-forward to establish or fail
			select {
			case err = <-portForwardErr:
				return mdexec.ExitError(fmt.Errorf("port-forward failed: %w", err))
			case <-time.After(2 * time.Second):
				// Port forward likely started, continue
			}

			mysqlArgs := []string{"-h", "127.0.0.1", "-P", fmt.Sprintf("%d", port), "-u", "root", "-p" + rootPass, "--prompt=tidb> "}
			redactedMysqlArgs := []string{"-h", "127.0.0.1", "-P", fmt.Sprintf("%d", port), "-u", "root", "-pPASS", "--prompt=tidb> "}

			if debug {
				colorDebugPrintfln(context, "%s %s", "mysql", strings.Join(redactedMysqlArgs, " "))
			}

			mysqlCmd := exec.Command("mysql", mysqlArgs...)
			mysqlCmd.Stdin = os.Stdin
			mysqlCmd.Stdout = os.Stdout
			mysqlCmd.Stderr = os.Stderr
			if err = mysqlCmd.Run(); err != nil {
				return mdexec.ExitError(err)
			}

			return nil
		},
	}
}

func tidbDmctlCommand() *cli.Command {
	return &cli.Command{
		Name:    "dmctl",
		Aliases: []string{"tdmctl"},
		Usage:   "Connect to dmctl on a TiDB cluster",
		Flags: append(append(append(mdk8s.BaseK8sFlags,
			mdk8s.BaseKctlFlags...), BaseTidbFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
			&cli.BoolFlag{
				Name:    "worker",
				Aliases: []string{"w"},
				Value:   false,
				Usage:   "Whether to connect to dm-worker. Defaults to false (dm-master).",
			},
		),
		Action: func(cCtx *cli.Context) error {
			cfg := config.NewConfig()

			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")
			useWorker := cCtx.Bool("worker")
			disableTls := cCtx.Bool("disable-tls")
			debug := cCtx.Bool("debug")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, _, err = ParseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			clusterName := strings.TrimPrefix(namespace, "tidb-")

			var podName, container, tlsPath string
			if useWorker {
				podName = fmt.Sprintf("%s-dm-worker-%d", clusterName, pod)
				container = "dm-worker"
				tlsPath = "/var/lib/dm-worker-tls"
			} else {
				podName = fmt.Sprintf("%s-dm-master-%d", clusterName, pod)
				container = "dm-master"
				tlsPath = "/var/lib/dm-master-tls"
			}

			var dmMasterEndpoint, dmctlCmd string
			if disableTls {
				dmMasterEndpoint = fmt.Sprintf("http://%s-dm-master:8261", clusterName)
				dmctlCmd = fmt.Sprintf(`./dmctl --master-addr %s `, dmMasterEndpoint)
			} else {
				dmMasterEndpoint = fmt.Sprintf("https://%s-dm-master:8261", clusterName)
				dmctlCmd = fmt.Sprintf(`./dmctl --master-addr %s --ssl-cert %s/tls.crt --ssl-key %s/tls.key --ssl-ca %s/ca.crt`, dmMasterEndpoint, tlsPath, tlsPath, tlsPath)
			}

			if len(cCtx.Args().Slice()) > 0 {
				dmctlCmd += " "
				dmctlCmd += strings.Join(cCtx.Args().Slice(), " ")
			}

			if isTestTidbContext(context) {
				assumeClusterAdmin = assumeClusterAdmin || cfg.EnableClusterAdminForTest
			}
			builder := NewTidbKubeBuilder()
			execArgs, _ := builder.BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"exec", "-it", podName, "-c", container, "--", "bin/sh", "-c", dmctlCmd})

			if debug {
				colorDebugPrintfln(context, "%s %s", "kubectl", strings.Join(execArgs, " "))
			}

			return mdexec.ExitError(mdexec.RunCommand("kubectl", execArgs...))
		},
	}
}

func tidbPdctlCommand() *cli.Command {
	return &cli.Command{
		Name:    "pdctl",
		Aliases: []string{"tpdctl"},
		Usage:   "Connect to pdctl on a TiDB cluster",
		Flags: append(append(append(mdk8s.BaseK8sFlags,
			mdk8s.BaseKctlFlags...), BaseTidbFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
		),
		Action: func(cCtx *cli.Context) error {
			cfg := config.NewConfig()

			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")
			disableTls := cCtx.Bool("disable-tls")
			debug := cCtx.Bool("debug")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, _, err = ParseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			clusterName := strings.TrimPrefix(namespace, "tidb-")
			var pdctlCmd, podName, container, tlsPath string
			podName = fmt.Sprintf("%s-pd-%d", clusterName, pod)
			container = "pd"
			if disableTls {
				pdctlCmd = "./pd-ctl -u http://127.0.0.1:2379 -i"
			} else {
				tlsPath = "/var/lib/cluster-client-tls"

				pdctlCmd = fmt.Sprintf(`./pd-ctl -u https://127.0.0.1:2379 --cert %s/tls.crt --key %s/tls.key --cacert %s/ca.crt -i`, tlsPath, tlsPath, tlsPath)
			}

			if isTestTidbContext(context) {
				assumeClusterAdmin = assumeClusterAdmin || cfg.EnableClusterAdminForTest
			}
			builder := NewTidbKubeBuilder()
			execArgs, _ := builder.BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"exec", "-it", podName, "-c", container, "--", "bin/sh", "-c", pdctlCmd})

			if debug {
				colorDebugPrintfln(context, "%s %s", "kubectl", strings.Join(execArgs, " "))
			}

			return mdexec.ExitError(mdexec.RunCommand("kubectl", execArgs...))
		},
	}
}

func ticdcCommand() *cli.Command {
	return &cli.Command{
		Name:    "ticdc",
		Aliases: []string{"cdc"},
		Usage:   "Execute cdc command on a TiDB cluster",
		Flags: append(append(append(mdk8s.BaseK8sFlags,
			mdk8s.BaseKctlFlags...), BaseTidbFlags...),
			&cli.IntFlag{
				Name:    "pod",
				Aliases: []string{"p"},
				Value:   0,
				Usage:   "Pod number to connect to. Defaults to 0.",
			},
		),
		Action: func(cCtx *cli.Context) error {
			cfg := config.NewConfig()

			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			pod := cCtx.Int("pod")
			assumeClusterAdmin := cCtx.Bool("assume-cluster-admin")
			disableTls := cCtx.Bool("disable-tls")
			debug := cCtx.Bool("debug")

			context = inferContextFromNamespace(context, namespace)

			var err error
			context, err = ParseContext(context, interactive, "^m-tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			namespace, _, err = ParseNamespace(namespace, false, interactive, context, "^tidb-", strict)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			clusterName := strings.TrimPrefix(namespace, "tidb-")
			podName := fmt.Sprintf("%s-ticdc-%d", clusterName, pod)
			var pdEndpoint, cdcCmd string
			if disableTls {
				pdEndpoint = fmt.Sprintf("http://%s-pd:2379", clusterName)
				cdcCmd = fmt.Sprintf("./cdc cli --pd %s %s", pdEndpoint, strings.Join(cCtx.Args().Slice(), " "))
			} else {
				tlsPath := "/var/lib/cluster-client-tls"
				pdEndpoint = fmt.Sprintf("https://%s-pd:2379", clusterName)
				cdcCmd = fmt.Sprintf("./cdc cli --pd %s --cert %s/tls.crt --key %s/tls.key --ca %s/ca.crt %s", pdEndpoint, tlsPath, tlsPath, tlsPath, strings.Join(cCtx.Args().Slice(), " "))
			}

			if isTestTidbContext(context) {
				assumeClusterAdmin = assumeClusterAdmin || cfg.EnableClusterAdminForTest
			}
			builder := NewTidbKubeBuilder()
			args, _ := builder.BuildKubectlArgs(context, namespace, false, assumeClusterAdmin, []string{"exec", "-it", podName, "-c", "ticdc", "--", "bin/sh", "-c", cdcCmd})

			if debug {
				colorDebugPrintfln(context, "%s %s", "kubectl", strings.Join(args, " "))
			}

			return mdexec.ExitError(mdexec.RunCommand("kubectl", args...))
		},
	}
}
