package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bitfield/script"
	mdexec "github.com/michaelmdeng/mdcli/internal/cmd"
	"github.com/urfave/cli/v3"
)

const (
	Kubectl = "kubectl"
	K9s     = "k9s"
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
		Name:    "debug",
		Aliases: []string{"d"},
		Value:   false,
		Usage:   "Preview the actual command to be executed",
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
		Aliases: []string{"k8s"},
		Usage:   k8sUsage,
		Commands: []*cli.Command{
			kubectlCommand(),
			k9sCommand(),
		},
	}
}

func kubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "kubectl",
		Aliases: []string{ "kc", "kctl"},
		Usage:   "Custom kubectl wrapper",
		Flags:   BaseK8sFlags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			strict := cmd.Bool("strict")
			context := cmd.String("context")
			namespace := cmd.String("namespace")
			interactive := cmd.Bool("interactive")
			dryRun := cmd.Bool("dryrun")
			allNamespaces := cmd.Bool("all-namespaces")
			assumeClusterAdmin := cmd.Bool("assume-cluster-admin")

			var err error
			context, err = ParseContext(context, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			builder := NewKubeBuilder()
			args, confirm := builder.BuildKubectlArgs(context, namespace, allNamespaces, assumeClusterAdmin, cmd.Args().Slice())
			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", Kubectl, strings.Join(args, " ")))
				return nil
			} else if confirm {
				fmt.Println(fmt.Sprintf("%s %s", Kubectl, strings.Join(args, " ")))
				res := mdexec.GetConfirmation("Do you want to execute the above command?")
				if !res {
					fmt.Println("Command canceled")
					return errors.New("Command canceled")
				}
			}

			return mdexec.RunCommand(Kubectl, args...)
		},
	}
}

func k9sCommand() *cli.Command {
	return &cli.Command{
		Name:  "k9s",
		Usage: "Custom k9s wrapper",
		Flags: BaseK8sFlags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			strict := cmd.Bool("strict")
			context := cmd.String("context")
			namespace := cmd.String("namespace")
			interactive := cmd.Bool("interactive")
			dryRun := cmd.Bool("dryrun")
			allNamespaces := cmd.Bool("all-namespaces")

			var err error
			context, err = ParseContext(context, interactive, "", strict)
			if err != nil {
				return err
			}

			namespace, allNamespaces, err = ParseNamespace(namespace, allNamespaces, interactive, context, "", strict)
			if err != nil {
				return err
			}

			builder := NewKubeBuilder()
			args, err := builder.BuildK9sArgs(context, namespace, allNamespaces, cmd.Args().Slice())
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println(fmt.Sprintf("%s %s", K9s, strings.Join(args, " ")))
				return nil
			}

			_, err = script.Exec(fmt.Sprintf("%s %s", K9s, strings.Join(args, " "))).Stdout()
			return err
		},
	}
}
