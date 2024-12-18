package wiki

import (
	"context"
	"os"
	"path"

	mdexec "github.com/mdcli/cmd"
	"github.com/urfave/cli/v3"
)

const wikiUsage = `Provides commands for managing my personal wiki`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "wiki",
		Aliases: []string{"w"},
		Usage:   wikiUsage,
		Commands: []*cli.Command{
			convertCommand(),
			transformCommand(),
			openCommand(),
		},
	}
}

const convertUsage = `Converts a wiki page in markdown to html`

func convertCommand() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   convertUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "css",
				Aliases: []string{"c"},
				Value:   ".local/share/pandoc/templates/default.css",
				Usage:   "CSS template `FILE` to convert with",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Value:   ".local/share/pandoc/templates/default.html5",
				Usage:   "HTML template `FILE` to convert with",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force conversion of md to html",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			inputPath := cmd.Args().First()
			outputPath, err := HtmlOutputPath(inputPath)
			if err != nil {
				return err
			}

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cmd.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cmd.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			return Convert(inputPath, outputPath, templateAbsPath, cssAbsPath, cmd.Bool("force"))
		},
	}
}

const transformUsage = `Converts a wiki folder in markdown to html`

func transformCommand() *cli.Command {
	return &cli.Command{
		Name:    "transform",
		Aliases: []string{"t"},
		Usage:   transformUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "css",
				Aliases: []string{"c"},
				Value:   ".local/share/pandoc/templates/default.css",
				Usage:   "CSS template `FILE` to transform with",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Value:   ".local/share/pandoc/templates/default.html5",
				Usage:   "HTML template `FILE` to transform with",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force conversion of md to html",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			inputDir := cmd.Args().First()
			htmlDir := path.Join(inputDir, "../html")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cmd.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cmd.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			return Transform(inputDir, htmlDir, templateAbsPath, cssAbsPath, cmd.Bool("force"))
		},
	}
}

const openUsage = `Opens a wiki page in the browser

Converts the page to html if necessary, or if force is set`

func openCommand() *cli.Command {
	return &cli.Command{
		Name:    "open",
		Aliases: []string{"o"},
		Usage:   openUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "css",
				Aliases: []string{"c"},
				Value:   ".local/share/pandoc/templates/default.css",
				Usage:   "CSS template `FILE` to convert with",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Value:   ".local/share/pandoc/templates/default.html5",
				Usage:   "HTML template `FILE` to convert with",
			},
			&cli.StringFlag{
				Name:    "browser",
				Aliases: []string{"b"},
				Value:   "firefox",
				Usage:   "Browser to open the page in",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force conversion of md to html",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cmd.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cmd.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			inputPath := cmd.Args().First()
			_, err = basePath(inputPath)
			var outputPath string
			if err != nil {
				outputPath, err = convertTemp(inputPath, templateAbsPath, cssAbsPath)
				if err != nil {
					return err
				}
			} else {
				outputPath, err = HtmlOutputPath(inputPath)
				if err != nil {
					return err
				}

				err = Convert(inputPath, outputPath, templateAbsPath, cssAbsPath, cmd.Bool("force"))
				if err != nil {
					return err
				}
			}

			err = mdexec.RunCommand(cmd.String("browser"), outputPath)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
