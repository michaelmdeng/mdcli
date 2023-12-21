package wiki

import (
	"os"
	"path"

	"github.com/mdcli/cmd"
	"github.com/urfave/cli/v2"
)

const wikiUsage = `Provides commands for managing my personal wiki`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "wiki",
		Aliases: []string{"w"},
		Usage:   wikiUsage,
		Subcommands: []*cli.Command{
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
		Action: func(cCtx *cli.Context) error {
			inputPath := cCtx.Args().First()
			outputPath, err := HtmlOutputPath(inputPath)
			if err != nil {
				return err
			}

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cCtx.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cCtx.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			return Convert(inputPath, outputPath, templateAbsPath, cssAbsPath, cCtx.Bool("force"))
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
		Action: func(cCtx *cli.Context) error {
			inputDir := cCtx.Args().First()
			htmlDir := path.Join(inputDir, "../html")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cCtx.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cCtx.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			return Transform(inputDir, htmlDir, templateAbsPath, cssAbsPath, cCtx.Bool("force"))
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
		Action: func(cCtx *cli.Context) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cCtx.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cCtx.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			inputPath := cCtx.Args().First()
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

				err = Convert(inputPath, outputPath, templateAbsPath, cssAbsPath, cCtx.Bool("force"))
				if err != nil {
					return err
				}
			}

			err = cmd.RunCommand(cCtx.String("browser"), outputPath)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
