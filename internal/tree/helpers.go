package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// .env files bur hide koribile if --hidden flag is not used
func filterEntries(
	entries []os.DirEntry,
	showHidden bool,
) []os.DirEntry {
	if showHidden {
		return entries
	}

	var visibleEntries []os.DirEntry

	for _, entry := range entries {
		if !isHidden(entry.Name()) {
			visibleEntries = append(
				visibleEntries,
				entry,
			)
		}
	}

	return visibleEntries
}

// checks whether dot (.env) files
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// returns the folder name that should be displayed
func getRootName(rootPath string) string {
	absolutePath, err := filepath.Abs(rootPath)

	if err != nil {
		return filepath.Base(rootPath)
	}

	return filepath.Base(absolutePath)
}

// formatSize converts bytes to sizes like kb mb etc
func formatSize(size int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
		gigabyte = 1024 * megabyte
	)

	switch {
	case size >= gigabyte:
		return fmt.Sprintf(
			"%.2f GB",
			float64(size)/float64(gigabyte),
		)

	case size >= megabyte:
		return fmt.Sprintf(
			"%.2f MB",
			float64(size)/float64(megabyte),
		)

	case size >= kilobyte:
		return fmt.Sprintf(
			"%.2f KB",
			float64(size)/float64(kilobyte),
		)

	default:
		return fmt.Sprintf("%d B", size)
	}
}
