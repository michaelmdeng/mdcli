package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

const splitUsage = `Split a file into pieces with file-extension-based output`

func splitFileName(inputPath, outputPrefix string, num int) string {
	ext := filepath.Ext(inputPath)
	if ext == "" {
		return fmt.Sprintf("%s%d", outputPrefix, num)
	}
	return fmt.Sprintf("%s-%d%s", outputPrefix, num, ext)
}

func splitByLines(inputPath, outputPrefix string, numLines int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	partNum := 0
	linesWritten := 0
	var outFile *os.File

	for scanner.Scan() {
		if linesWritten == 0 {
			filename := splitFileName(inputPath, outputPrefix, partNum)
			outFile, err = os.Create(filename)
			if err != nil {
				return err
			}
			defer outFile.Close()
		}

		if _, err := outFile.WriteString(scanner.Text() + "\n"); err != nil {
			return err
		}
		linesWritten++

		if linesWritten >= numLines {
			linesWritten = 0
			partNum++
			outFile.Close()
		}
	}
	
	if outFile != nil && linesWritten > 0 {
		outFile.Close()
	}
	return scanner.Err()
}

func parseByteSize(s string) (int, error) {
	re := regexp.MustCompile(`^(\d+)([KkMGBg][i]?)?$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid byte size format")
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	if len(matches[2]) == 0 {
		return num, nil
	}

	switch strings.ToLower(matches[2]) {
	case "k":
		return num * 1000, nil
	case "ki":
		return num * 1024, nil
	case "m":
		return num * 1000 * 1000, nil
	case "mi":
		return num * 1024 * 1024, nil
	case "g":
		return num * 1000 * 1000 * 1000, nil
	case "gi":
		return num * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unknown size suffix: %s", matches[2])
	}
}

func parseLineAwareFlag(value string) (int, bool, error) {
	if strings.HasPrefix(value, "l/") {
		sizeStr := strings.TrimPrefix(value, "l/")
		size, err := parseByteSize(sizeStr)
		return size, true, err
	}
	size, err := parseByteSize(value)
	return size, false, err
}

func splitByBytes(inputPath, outputPrefix string, byteSize int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	partNum := 0
	buf := make([]byte, byteSize)

	for {
		n, err := file.Read(buf)
		if n > 0 {
			filename := splitFileName(inputPath, outputPrefix, partNum)
			outFile, err := os.Create(filename)
			if err != nil {
				return err
			}
			if _, werr := outFile.Write(buf[:n]); werr != nil {
				outFile.Close()
				return werr
			}
			outFile.Close()
			partNum++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func splitByChunks(inputPath, outputPrefix string, numChunks int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := stat.Size()
	chunkSize := fileSize / int64(numChunks)
	for i := 0; i < numChunks; i++ {
		filename := splitFileName(inputPath, outputPrefix, i)
		outFile, err := os.Create(filename)
		if err != nil {
			return err
		}

		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numChunks-1 {
			end = fileSize
		}

		if _, err := file.Seek(start, 0); err != nil {
			outFile.Close()
			return err
		}

		_, err = io.CopyN(outFile, file, end-start)
		outFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func createSplitFile(inputPath, outputPrefix string, partNum int) *os.File {
	filename := splitFileName(inputPath, outputPrefix, partNum)
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}
	return file
}

func splitByBytesLineAware(inputPath, outputPrefix string, byteSize int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	partNum := 0
	currentFile := createSplitFile(inputPath, outputPrefix, partNum)
	if currentFile == nil {
		return fmt.Errorf("failed to create output file")
	}
	defer currentFile.Close()
	currentSize := 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lineSize := len(line)
		if currentSize > 0 && currentSize+lineSize > byteSize {
			currentFile.Close()
			partNum++
			currentFile = createSplitFile(inputPath, outputPrefix, partNum)
			if currentFile == nil {
				return fmt.Errorf("failed to create output file")
			}
			currentSize = 0
		}

		if _, werr := currentFile.Write(line); werr != nil {
			return werr
		}
		currentSize += lineSize

		if err == io.EOF {
			break
		}
	}
	return nil
}

func splitByChunksLineAware(inputPath, outputPrefix string, numChunks int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	partNum := 0
	bytesPerChunk := 0
	file.Seek(0, 0)
	size, err := file.Seek(0, 2)
	if err != nil {
		return err
	}
	bytesPerChunk = int(size) / numChunks

	file.Seek(0, 0)
	currentSize := 0
	currentFile := createSplitFile(inputPath, outputPrefix, partNum)
	if currentFile == nil {
		return fmt.Errorf("failed to create output file")
	}
	defer currentFile.Close()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lineSize := len(line)
		if currentSize > 0 && currentSize+lineSize > bytesPerChunk && partNum < numChunks-1 {
			currentFile.Close()
			partNum++
			currentFile = createSplitFile(inputPath, outputPrefix, partNum)
			if currentFile == nil {
				return fmt.Errorf("failed to create output file")
			}
			currentSize = 0
		}

		if _, werr := currentFile.Write(line); werr != nil {
			return werr
		}
		currentSize += lineSize

		if err == io.EOF {
			break
		}
	}
	return nil
}

func BaseCommand() *cli.Command {
	return &cli.Command{
		Name:    "split",
		Aliases: []string{"sp"},
		Usage:   splitUsage,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "lines",
				Aliases: []string{"l"},
				Usage:   "Split into chunks of `l` lines",
			},
			&cli.StringFlag{
				Name:    "chunks",
				Aliases: []string{"n"},
				Usage:   "Split into `n` chunks. Use `l/n` for line-aware splitting",
			},
			&cli.StringFlag{
				Name:    "bytes",
				Aliases: []string{"b"},
				Usage:   "Split into chunks of `b` bytes. Supports unit suffixes, ex. Mi. Use `l/b` for line-aware splitting.)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() < 1 {
				return cli.Exit("Input file required", 1)
			}
			inputPath := cCtx.Args().Get(0)

			// Determine output prefix
			var outputPrefix string
			if cCtx.NArg() >= 2 {
				outputPrefix = cCtx.Args().Get(1)
			} else {
				// Set output prefix to input path without extension
				base := filepath.Base(inputPath)
				ext := filepath.Ext(base)
				nameWithoutExt := strings.TrimSuffix(base, ext)
				outputPrefix = filepath.Join(filepath.Dir(inputPath), nameWithoutExt)
			}

			switch {
			case cCtx.IsSet("lines"):
				return splitByLines(inputPath, outputPrefix, cCtx.Int("lines"))
			case cCtx.IsSet("chunks"):
				val := cCtx.String("chunks")
				if strings.HasPrefix(val, "l/") {
					chunks, err := strconv.Atoi(val[2:])
					if err != nil {
						return fmt.Errorf("invalid chunk size: %s", val[2:])
					}
					return splitByChunksLineAware(inputPath, outputPrefix, chunks)
				}
				chunks, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("invalid chunk size: %s", val)
				}
				return splitByChunks(inputPath, outputPrefix, chunks)
			case cCtx.IsSet("bytes"):
				val := cCtx.String("bytes")
				byteSize, byLine, err := parseLineAwareFlag(val)
				if err != nil {
					return err
				}
				if byLine {
					return splitByBytesLineAware(inputPath, outputPrefix, byteSize)
				}
				return splitByBytes(inputPath, outputPrefix, byteSize)
			default:
				return cli.Exit("One of -l, -c, or -b required", 1)
			}
		},
	}
}
