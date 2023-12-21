package k8s

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/bitfield/script"
)

const k8sUsage = `Custom wrapper for k8s commands.

Provides better short flags for kubectl and allows the user to select the
context and namespace interactively via fzf.`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "k8s",
		Aliases: []string{"k"},
		Usage:   k8sUsage,
		Subcommands: []*cli.Command{
			kubectlCommand(),
			k9sCommand(),
		},
	}
}

func kubectlCommand() *cli.Command {
	return &cli.Command{
		Name:    "kubectl",
		Aliases: []string{"k", "kctl"},
		Usage:   "Custom kubectl wrapper",
		Flags: []cli.Flag{
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
		},
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")

			var err error
			if context == "" && interactive {
				context, err = getContextInteractive()
				if strict && err != nil {
					return cli.Exit(err, 1)
				} else if err != nil {
					context = ""
				}
			}
			if context == "" && strict {
				return cli.Exit("context and namespace must be specified in strict-mode", 1)
			}

			if namespace == "" && interactive {
				namespace, err = getNamespaceInteractive(context)
				if strict && err != nil {
					return cli.Exit(err, 1)
				} else if err != nil {
					namespace = ""
				}
			}
			if namespace == "" && strict {
				return cli.Exit("context and namespace must be specified in strict-mode", 1)
			}

			var cmd string
			cmd = "kubectl"
			if context != "" {
				cmd = fmt.Sprintf("%s --context %s", cmd, context)
			}

			var allNs bool
			if namespace == "*" {
				allNs = true
			}
			if namespace != "" && !allNs {
				cmd = fmt.Sprintf("%s --namespace %s", cmd, namespace)
			}

			if allNs {
				cmd = fmt.Sprintf("%s %s %s", cmd, strings.Join(cCtx.Args().Slice(), " "), "--all-namespaces")
			} else {
				cmd = fmt.Sprintf("%s %s", cmd, strings.Join(cCtx.Args().Slice(), " "))
			}
			_, err = script.Exec(cmd).Stdout()
			return err
		},
	}
}

func k9sCommand() *cli.Command {
	return &cli.Command{
		Name:    "k9s",
		Usage:   "Custom k9s wrapper",
		Flags: []cli.Flag{
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
		},
		Action: func(cCtx *cli.Context) error {
			strict := cCtx.Bool("strict")
			context := cCtx.String("context")
			namespace := cCtx.String("namespace")
			interactive := cCtx.Bool("interactive")

			var err error
			if context == "" && interactive {
				context, err = getContextInteractive()
				if strict && err != nil {
					return cli.Exit(err, 1)
				} else if err != nil {
					context = ""
				}
			}
			if context == "" && strict {
				return cli.Exit("context and namespace must be specified in strict-mode", 1)
			}

			if namespace == "" && interactive {
				namespace, err = getNamespaceInteractive(context)
				if strict && err != nil {
					return cli.Exit(err, 1)
				} else if err != nil {
					namespace = ""
				}
			}
			if namespace == "" && strict {
				return cli.Exit("context and namespace must be specified in strict-mode", 1)
			}

			var cmd string
			cmd = "k9s"
			if context != "" {
				cmd = fmt.Sprintf("%s --context %s", cmd, context)
			}

			var allNs bool
			if namespace == "*" {
				allNs = true
			}
			if namespace != "" && !allNs {
				cmd = fmt.Sprintf("%s --namespace %s", cmd, namespace)
			}

			if allNs {
				cmd = fmt.Sprintf("%s -A %s", cmd, strings.Join(cCtx.Args().Slice(), " "))
			}
			_, err = script.Exec(cmd).Stdout()
			return err
		},
	}
}
