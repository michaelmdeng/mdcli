package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp" // Add this import
	"sort"
	"strings"
	"time"
)

// findScratchDirectory searches for a directory matching the name within the scratch path.
// It leverages listScratchDirectories and then checks for exact matches (YYYY-MM-DD-name)
// and suffix matches (-name) among the valid scratch directories.
// Returns the full path of the found directory or an empty string if not found.
// Returns an error if listing directories fails or multiple suffix matches are found.
func findScratchDirectory(scratchPath, name string) (string, error) {
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		// Although listScratchDirectories also does this, checking early avoids unnecessary work
		return "", fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	// Use listScratchDirectories to get only valid, full paths
	directories, err := listScratchDirectories(absScratchPath)
	if err != nil {
		// Propagate error from listing/checking the directory
		return "", err
	}

	var exactMatch string
	var suffixMatches []string
	suffixToCheck := "-" + name
	datePrefix := time.Now().Format("2006-01-02")
	exactName := fmt.Sprintf("%s-%s", datePrefix, name)

	// Iterate through the full paths returned by listScratchDirectories
	for _, fullPath := range directories {
		entryName := filepath.Base(fullPath) // Get the directory name from the full path

		// Check for exact match (including today's date)
		if entryName == exactName {
			exactMatch = fullPath
			break // Exact match found, no need to check further
		}

		// Check for suffix match
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
		// Sort by full path name to get the most recent if dates are prefixes
		sort.Strings(suffixMatches)
		// Return the last one (most recent date assuming YYYY-MM-DD prefix)
		return suffixMatches[len(suffixMatches)-1], nil
		// Alternatively, error out:
		// return "", fmt.Errorf("multiple directories found matching suffix '%s': %v", suffixToCheck, suffixMatches)
	}

	return "", nil // No match found
}

// createScratchDirectory creates a new dated directory within the scratch path.
// It formats the name as YYYY-MM-DD-name and creates the directory.
// If createReadme is true, it also creates an empty README.md file inside.
// Returns the full path of the created directory or an error.
func createScratchDirectory(scratchPath, name string, createReadme bool) (string, error) {
	// Format the new directory name
	today := time.Now().Format("2006-01-02")
	newDirName := fmt.Sprintf("%s-%s", today, name)
	// Note: scratchPath is expected to be absolute already by the callers
	newDirPath := filepath.Join(scratchPath, newDirName)

	// Double-check if it already exists (mitigate race conditions)
	// Although findScratchDirectory in the callers should prevent this call if it exists.
	if _, err := os.Stat(newDirPath); !os.IsNotExist(err) {
		if err == nil {
			// Directory surprisingly exists
			return "", fmt.Errorf("directory '%s' already exists unexpectedly", newDirPath)
		}
		// Other stat error
		return "", fmt.Errorf("failed to check directory status '%s': %w", newDirPath, err)
	}

	// Create the new directory
	if err := os.Mkdir(newDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory '%s': %w", newDirPath, err)
	}

	// Create README.md if requested
	if createReadme {
		readmePath := filepath.Join(newDirPath, "README.md")
		file, err := os.Create(readmePath)
		if err != nil {
			// Return error, but the directory was already created. Maybe log this?
			return newDirPath, fmt.Errorf("directory created, but failed to create README.md: %w", err)
		}
		file.Close() // Close the empty file
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

	// Check if the directory exists first
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
		// Check if it's a directory AND matches the pattern
		if entry.IsDir() && scratchDirPattern.MatchString(entry.Name()) {
			fullPath := filepath.Join(absScratchPath, entry.Name())
			directories = append(directories, fullPath)
		}
	}

	// Optionally sort the directories by name for consistent order
	sort.Strings(directories)

	return directories, nil
}
