package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// reads and displays one directory
func walkDirectory(
	currentPath string,
	prefix string,
	depth int,
	options Options,
	writer io.Writer,
	result *Result,
) error {
	// Stop when the maximum depth is reached
	if options.MaxDepth > 0 && depth > options.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(currentPath)

	if err != nil {
		return fmt.Errorf( "could not read directory %s: %w", currentPath, err )
	}

	// Remove hidden files when --hidden is not used.
	entries = filterEntries(entries, options.ShowHidden)

	// Sort alphabetically
	sort.Slice(entries, func(i int, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for index, entry := range entries {
		isLastEntry := index == len(entries)-1

		connector := "├── "

		if isLastEntry {
			connector = "└── "
		}

		entryName := entry.Name()

		fmt.Fprint(writer, prefix, connector, entryName)

		fullPath := filepath.Join(currentPath, entryName)

		if entry.IsDir() {
			result.DirectoryCount++
		} else {
			result.FileCount++

			if options.ShowSize {
				fileInfo, infoError := entry.Info()

				if infoError != nil {
					return fmt.Errorf( "could not read information for %s: %w", fullPath, infoError )
				}

				fmt.Fprintf(writer, " (%s)", formatSize(fileInfo.Size()))
			}
		}

		fmt.Fprintln(writer)

		// enter the directory recursively
		if entry.IsDir() {
			newPrefix := prefix

			if isLastEntry {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			err = walkDirectory(fullPath, newPrefix, depth+1, options, writer, result)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
