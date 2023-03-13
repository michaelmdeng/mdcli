package wiki

import (
	"os"
	"path"
	"strings"

	"github.com/mdcli/cmd"
	"github.com/urfave/cli/v2"
)

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "wiki",
		Aliases: []string{"w"},
		Usage:   "Tasks for personal wiki",
		Subcommands: []*cli.Command{
			convertCommand(),
			transformCommand(),
			openCommand(),
		},
	}
}

func convertCommand() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   "Convert a wiki page to html",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "css",
			Aliases: []string{"c"},
			Value: ".local/share/pandoc/templates/default.css",
			Usage: "CSS template `FILE`",
		    },
		    &cli.StringFlag{
			Name:  "template",
			Aliases: []string{"t"},
			Value: ".local/share/pandoc/templates/default.html5",
			Usage: "HTML template `FILE`",
		    },
		    &cli.BoolFlag{
			Name:  "force",
			Aliases: []string{"f"},
			Usage: "Force conversion of md to html",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			inputPath := cCtx.Args().First()

			fileNameExt := path.Base(inputPath)
			fileExt := path.Ext(inputPath)
			fileName := strings.TrimSuffix(fileNameExt, fileExt)
			htmlPath := path.Join(path.Dir(inputPath), "../html")
			outputPath := path.Join(htmlPath, fileName + ".html")

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

func transformCommand() *cli.Command {
	return &cli.Command{
		Name:    "transform",
		Aliases: []string{"t"},
		Usage:   "Convert an entire wiki folder to html",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "css",
			Aliases: []string{"c"},
			Value: ".local/share/pandoc/templates/default.css",
			Usage: "CSS template `FILE`",
		    },
		    &cli.StringFlag{
			Name:  "template",
			Aliases: []string{"t"},
			Value: ".local/share/pandoc/templates/default.html5",
			Usage: "HTML template `FILE`",
		    },
		    &cli.BoolFlag{
			Name:  "force",
			Aliases: []string{"f"},
			Usage: "Force conversion of md to html",
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

func openCommand() *cli.Command {
	return &cli.Command{
		Name:    "open",
		Aliases: []string{"o"},
		Usage:   "Open a wiki page in the browser",
		Flags: []cli.Flag{
		    &cli.StringFlag{
			Name:  "css",
			Aliases: []string{"c"},
			Value: ".local/share/pandoc/templates/default.css",
			Usage: "CSS template `FILE`",
		    },
		    &cli.StringFlag{
			Name:  "template",
			Aliases: []string{"t"},
			Value: ".local/share/pandoc/templates/default.html5",
			Usage: "HTML template `FILE`",
		    },
		    &cli.StringFlag{
			Name:  "browser",
			Aliases: []string{"b"},
			Value: "firefox",
			Usage: "Browser to open the page in",
		    },
		    &cli.BoolFlag{
			Name:  "force",
			Aliases: []string{"f"},
			Usage: "Force conversion of md to html",
		    },
		},
		Action: func(cCtx *cli.Context) error {
			inputPath := cCtx.Args().First()

			fileNameExt := path.Base(inputPath)
			fileExt := path.Ext(inputPath)
			fileName := strings.TrimSuffix(fileNameExt, fileExt)
			htmlPath := path.Join(path.Dir(inputPath), "../html")
			outputPath := path.Join(htmlPath, fileName + ".html")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cssPath := cCtx.String("css")
			cssAbsPath := path.Join(homeDir, cssPath)
			templatePath := cCtx.String("template")
			templateAbsPath := path.Join(homeDir, templatePath)

			err = Convert(inputPath, outputPath, templateAbsPath, cssAbsPath, cCtx.Bool("force"))
			if err != nil {
				return err
			}
			err = cmd.RunCommand(cCtx.String("browser"), outputPath)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
