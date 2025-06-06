package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// expandPath expands ~ and returns an absolute path.
// It takes a path string which might start with ~/, ~/ or be relative/absolute.
// It returns the absolute path or an error if the home directory cannot be determined.
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		// Handle cases like "~" or "~/..."
		if path == "~" {
			path = home
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		} else {
			// If it's just "~something", treat it relative to CWD after this point,
			// which filepath.Abs will handle. This case is less common for config paths.
			// Or, alternatively, decide if "~something" without a slash should also be relative to home.
			// Let's assume for now it's not relative to home unless followed by /.
			// If needed, add: path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for '%s': %w", path, err)
	}
	return absPath, nil
}

// findScratchDirectory searches for a directory matching the name within the scratch path.
// It leverages listScratchDirectories and then checks for exact matches (YYYY-MM-DD-name)
// and suffix matches (-name) among the valid scratch directories.
// Returns the full path of the found directory or an empty string if not found.
// Returns an error if listing directories fails or multiple suffix matches are found.
func findScratchDirectory(scratchPath, name string) (string, error) {
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	directories, err := listScratchDirectories(absScratchPath)
	if err != nil {
		return "", err
	}

	var exactMatch string
	var suffixMatches []string
	suffixToCheck := "-" + name
	datePrefix := time.Now().Format("2006-01-02")
	exactName := fmt.Sprintf("%s-%s", datePrefix, name)

	for _, fullPath := range directories {
		entryName := filepath.Base(fullPath)

		if entryName == exactName {
			exactMatch = fullPath
			break
		}

		if strings.HasSuffix(entryName, suffixToCheck) {
			suffixMatches = append(suffixMatches, fullPath)
		}
	}

	if exactMatch != "" {
		return exactMatch, nil
	}

	if len(suffixMatches) == 1 {
		return suffixMatches[0], nil
	}

	if len(suffixMatches) > 1 {
		sort.Strings(suffixMatches)
		// Return the last one (most recent date assuming YYYY-MM-DD prefix)
		return suffixMatches[len(suffixMatches)-1], nil
	}

	return "", nil
}

// createScratchDirectory creates a new dated directory within the scratch path.
// It formats the name as YYYY-MM-DD-name and creates the directory.
// If createReadme is true, it also creates an empty README.md file inside.
// Returns the full path of the created directory or an error.
func createScratchDirectory(scratchPath, name string, createReadme bool) (string, error) {
	today := time.Now().Format("2006-01-02")
	newDirName := fmt.Sprintf("%s-%s", today, name)
	newDirPath := filepath.Join(scratchPath, newDirName)

	if _, err := os.Stat(newDirPath); !os.IsNotExist(err) {
		if err == nil {
			return "", fmt.Errorf("directory '%s' already exists unexpectedly", newDirPath)
		}
		return "", fmt.Errorf("failed to check directory status '%s': %w", newDirPath, err)
	}

	if err := os.Mkdir(newDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory '%s': %w", newDirPath, err)
	}

	if createReadme {
		readmePath := filepath.Join(newDirPath, "README.md")
		file, err := os.Create(readmePath)
		if err != nil {
			return newDirPath, fmt.Errorf("directory created, but failed to create README.md: %w", err)
		}
		file.Close()
	}

	return newDirPath, nil
}

// listScratchDirectories lists all directories directly under the given scratch path
// that match the format YYYY-MM-DD-<name>.
// It returns a slice of full paths to the matching directories.
func listScratchDirectories(scratchPath string) ([]string, error) {
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	if _, err := os.Stat(absScratchPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("scratch directory '%s' does not exist", absScratchPath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check scratch directory '%s': %w", absScratchPath, err)
	}

	entries, err := os.ReadDir(absScratchPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scratch directory '%s': %w", absScratchPath, err)
	}

	// Compile regex to match YYYY-MM-DD-<name> format
	// ^         - start of string
	// \d{4}     - exactly four digits (year)
	// -         - literal hyphen
	// \d{2}     - exactly two digits (month)
	// -         - literal hyphen
	// \d{2}     - exactly two digits (day)
	// -         - literal hyphen
	// .+        - one or more characters (name part)
	// $         - end of string
	scratchDirPattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-.+$`)

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() && scratchDirPattern.MatchString(entry.Name()) {
			fullPath := filepath.Join(absScratchPath, entry.Name())
			directories = append(directories, fullPath)
		}
	}

	sort.Strings(directories)

	return directories, nil
}
