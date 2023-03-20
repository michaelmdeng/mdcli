package wiki

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/mdcli/cmd"
)

var (
	linkWhitespacePattern = regexp.MustCompile(`\s+`)
)

func GenerateLink(text string) string {
	return fmt.Sprintf("[%v](%v)", text, strings.ToLower(linkWhitespacePattern.ReplaceAllString(text, "-")))
}

func Transform(inputDir string, outputDir string, template string, css string, force bool) error {
	fileSystem := os.DirFS(inputDir)
	return fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || path.Ext(p) != ".md" {
			return nil
		}

		inputPath := path.Join(inputDir, p)
		outputPath := path.Join(outputDir, strings.TrimSuffix(p, ".md") + ".html")
		return Convert(inputPath, outputPath, template, css, force)
	})
}

func Convert(input string, output string, template string, css string, force bool) error {
	if force {
		return pandocConvert(input, output, template, css)
	}

	inputStat, err := os.Stat(input)
	if errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var shouldConvert bool
	if outputStat, err := os.Stat(input); os.IsNotExist(err) {
		shouldConvert = true
	} else {
		if inputStat.ModTime().After(outputStat.ModTime()) {
			shouldConvert = true
		} else {
			shouldConvert = false
		}
	}

	if shouldConvert {
		return pandocConvert(input, output, template, css)
	}

	return nil
}

func pandocConvert(input string, output string, template string, css string) error {
	fileNameExt := path.Base(output)
	fileExt := path.Ext(output)
	fileName := strings.TrimSuffix(fileNameExt, fileExt)
	err := cmd.RunCommand(
		"pandoc", input,
		"-r", "markdown",
		"-w", "html",
		"-o", output,
		"--metadata", fmt.Sprintf("title=\"%v\"", fileName),
		"--template", template,
		"--css", css,
	)
	if err != nil {
		return err
	}

	return nil
}
