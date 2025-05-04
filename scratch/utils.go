package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// findScratchDirectory searches for a directory matching the name within the scratch path.
// It prioritizes exact matches (YYYY-MM-DD-name) and then suffix matches (-name).
// Returns the full path of the found directory or an empty string if not found.
// Returns an error if multiple suffix matches are found.
func findScratchDirectory(scratchPath, name string) (string, error) {
	absScratchPath, err := filepath.Abs(scratchPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for scratch directory '%s': %w", scratchPath, err)
	}

	entries, err := os.ReadDir(absScratchPath)
	if err != nil {
		return "", fmt.Errorf("failed to read scratch directory '%s': %w", absScratchPath, err)
	}

	var exactMatch string
	var suffixMatches []string
	suffixToCheck := "-" + name
	datePrefix := time.Now().Format("2006-01-02")
	exactName := fmt.Sprintf("%s-%s", datePrefix, name)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		fullPath := filepath.Join(absScratchPath, entryName)

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
		// Sort by name to get the most recent if dates are prefixes
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
// Returns the full path of the created directory or an error.
func createScratchDirectory(scratchPath, name string) (string, error) {
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

	return newDirPath, nil
}
