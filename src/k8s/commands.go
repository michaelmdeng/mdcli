package k8s

import (
	"fmt"
	"strings"

	"github.com/bitfield/script"
	mdexec "github.com/mdcli/cmd"
	"github.com/urfave/cli/v2"
)

const k8sUsage = `Custom wrapper for k8s commands.

Provides better short flags for kubectl and allows the user to select the
context and namespace interactively via fzf.`

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

			var err error
			context, err = parseContext(context, noop, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, noop, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			cmd := "kubectl"
			args := make([]string, 0)
			if context != "" {
				args = append(args, "--context", context)
			}

			if !allNamespaces && namespace != "" {
				args = append(args, "--namespace", namespace)
			}

			args = append(args, cCtx.Args().Slice()...)

			if allNamespaces {
				args = append(args, "--all-namespaces")
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
				return nil
			}

			return mdexec.RunCommand(cmd, args...)
		},
	}
}

func tidbKubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "tidbkubectl",
		Aliases: []string{"tk", "tkc", "tkctl"},
		Usage:   "Custom kubectl wrapper for TiDB",
		Flags: BaseK8sFlags,
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")
			dryRun := cCtx.Bool("dryrun")
			allNamespaces := cCtx.Bool("all-namespaces")

			var err error
			context, err = parseContext(context, tidbContext, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, tidbNamespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			cmd := "kubectl"
			args := make([]string, 0)
			if context != "" {
				args = append(args, "--context", context)
			}

			if namespace != "" && !allNamespaces {
				args = append(args, "--namespace", namespace)
			}

			args = append(args, cCtx.Args().Slice()...)

			if allNamespaces {
				args = append(args, "--all-namespaces")
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
				return nil
			}

			return mdexec.RunCommand(cmd, args...)
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
			context, err = parseContext(context, noop, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, noop, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			cmd := "k9s"
			args := make([]string, 0)
			if context != "" {
				args = append(args, "--context", context)
			}

			if allNamespaces {
				args = append(args, "--all-namespaces")
			}
			if namespace != "" && !allNamespaces {
				args = append(args, "--namespace", namespace)
			}

			var trailArg string
			if cCtx.Args().Len() == 2 && cCtx.Args().Get(0) == "get" {
				trailArg = cCtx.Args().Get(1)
				args = append(args, "-c", trailArg)
			} else if cCtx.Args().Len() == 1 {
				trailArg = cCtx.Args().Get(0)
				args = append(args, "-c", trailArg)
			} else if cCtx.Args().Len() == 0 {
				// do nothing
			} else {
				return cli.Exit(fmt.Sprintf("too many arguments provided to k9s: %s", strings.Join(cCtx.Args().Slice(), " ")), 1)
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
				return nil
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))).Stdout()
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
			context, err = parseContext(context, tidbContext, interactive, "^m-tidb-", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = parseNamespace(namespace, tidbNamespace, allNamespaces, interactive, context, "^tidb-", strict)
			if err != nil {
				return err
			}

			cmd := "k9s"
			args := make([]string, 0)
			if context != "" {
				args = append(args, "--context", context)
			}

			if allNamespaces {
				args = append(args, "--all-namespaces")
			}
			if namespace != "" && !allNamespaces {
				args = append(args, "--namespace", namespace)
			}

			var trailArg string
			if cCtx.Args().Len() == 2 && cCtx.Args().Get(0) == "get" {
				trailArg = cCtx.Args().Get(1)
				args = append(args, "-c", trailArg)
			} else if cCtx.Args().Len() == 1 {
				trailArg = cCtx.Args().Get(0)
				args = append(args, "-c", trailArg)
			} else if cCtx.Args().Len() == 0 {
				// do nothing
			} else {
				return cli.Exit(fmt.Sprintf("too many arguments provided to k9s: %s", strings.Join(cCtx.Args().Slice(), " ")), 1)
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
				return nil
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))).Stdout()
			return err
		},
	}
}
