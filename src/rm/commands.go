package rm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

const remarkableUsage = `Commands for working with reMarkable device.`

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "remarkable",
		Aliases: []string{"rm", "rem"},
		Usage:   remarkableUsage,
		Commands: []*cli.Command{
			documentsCommand(),
			downloadCommand(),
		},
	}
}

func documentsCommand() *cli.Command {
	return &cli.Command{
		Name:    "documents",
		Aliases: []string{"docs", "doc"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "Output format for response",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			documents, err := getDocuments()
			if err != nil {
				return err
			}

			output, err := json.Marshal(documents)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
}

func downloadCommand() *cli.Command {
	return &cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "documentId",
				Aliases: []string{"id", "d"},
				Value:   "",
				Usage:   "Document ID to download",
			},
			&cli.StringFlag{
				Name:    "documentName",
				Aliases: []string{"name", "n"},
				Value:   "",
				Usage:   "Document name to download",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "File to download to",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			documentId := cmd.String("documentId")
			documentName := cmd.String("documentName")
			var err error
			if documentId != "" && documentName != "" {
				return fmt.Errorf("only one of documentId or documentName can be specified")
			} else if documentId == "" && documentName == "" {
				return fmt.Errorf("one of documentId or documentName must be specified")
			} else if documentName != "" {
				documentId, err = getDocumentIdByName(documentName)
				if err != nil {
					return err
				}
			}

			err = downloadDocument(documentId, cmd.String("output"))
			if err != nil {
				return err
			}

			return nil
		},
	}
}
