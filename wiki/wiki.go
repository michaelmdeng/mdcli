package wiki

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/michaelmdeng/mdcli/internal/cmd"
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
		outputPath := path.Join(outputDir, strings.TrimSuffix(p, ".md")+".html")
		return Convert(inputPath, outputPath, template, css, force)
	})
}

func Convert(input string, output string, template string, css string, force bool) error {
	if force {
		return pandocConvert(input, output, template, css, "")
	}

	inputStat, err := os.Stat(input)
	if errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var shouldConvert bool
	if outputStat, err := os.Stat(output); os.IsNotExist(err) {
		shouldConvert = true
	} else {
		if inputStat.ModTime().After(outputStat.ModTime()) {
			shouldConvert = true
		} else {
			shouldConvert = false
		}
	}

	if shouldConvert {
		return pandocConvert(input, output, template, css, "")
	}

	return nil
}

func convertTemp(input string, template string, css string) (string, error) {
	output, err := tempHtmlOutputPath(input)
	if err != nil {
		return "", err
	}

	inputName := strings.TrimSuffix(path.Base(input), path.Ext(input))
	err = pandocConvert(input, output, template, css, inputName)
	if err != nil {
		return "", err
	}
	return output, nil
}

func pandocConvert(input string, output string, template string, css string, title string) error {
	fileNameExt := path.Base(output)
	fileExt := path.Ext(output)
	if title == "" {
		title = strings.TrimSuffix(fileNameExt, fileExt)
	}
	err := cmd.RunCommand(
		"pandoc", input,
		"-r", "markdown",
		"-w", "html",
		"-o", output,
		"--metadata", fmt.Sprintf("title=\"%v\"", title),
		"--template", template,
		"--css", css,
		"--embed-resources",
		"--standalone",
	)
	if err != nil {
		return err
	}

	return nil
}

func basePath(path string) (string, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	walkPath := fullPath
	var fullBasePath string
	for {
		dir := filepath.Dir(walkPath)
		if wikiDir := dir; filepath.Base(dir) == "wiki" {
			fullBasePath = filepath.Dir(wikiDir)
			break
		}

		if dir == "." || dir == "/" {
			return "", errors.New("unable to find wiki base path")
		}

		walkPath = dir
	}

	return fullBasePath, nil
}

func wikiPath(path string) (string, error) {
	basePath, err := basePath(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(basePath, "wiki"), nil
}

func htmlPath(path string) (string, error) {
	basePath, err := basePath(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(basePath, "html"), nil
}

func HtmlOutputPath(input string) (string, error) {
	inputPath, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}

	wikiPath, err := wikiPath(inputPath)
	if err != nil {
		return "", err
	}

	relInputPath, err := filepath.Rel(wikiPath, inputPath)
	if err != nil {
		return "", err
	}

	fileExt := path.Ext(inputPath)

	htmlPath, err := htmlPath(inputPath)
	if err != nil {
		return "", err
	}

	relHtmlPath := strings.TrimSuffix(relInputPath, fileExt) + ".html"
	outputPath := path.Join(htmlPath, relHtmlPath)

	return outputPath, nil
}

func tempHtmlOutputPath(input string) (string, error) {
	inputPath, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tmpDir := path.Join(homeDir, "tmp/mdcli")
	err = os.MkdirAll(tmpDir, 0744)
	if err != nil {
		return "", err
	}

	inputName := strings.TrimSuffix(path.Base(inputPath), path.Ext(inputPath))
	tmpFile, err := os.CreateTemp(tmpDir, fmt.Sprintf("%s-*.html", inputName))
	if err != nil {
		return "", err
	}

	err = tmpFile.Close()
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
