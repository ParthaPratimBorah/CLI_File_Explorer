package search

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Search starts searching using the provided options.
func Search(options Options) ([]Result, error) {
	var results []Result

	if options.Pattern == "" {
		return results, fmt.Errorf(
			"search pattern cannot be empty",
		)
	}

	// Create the matcher only once.
	matcher, err := newMatcher(options)

	if err != nil {
		return results, err
	}

	// Check whether the starting path exists.
	rootInfo, err := os.Stat(options.RootPath)

	if err != nil {
		return results, fmt.Errorf(
			"could not access path %s: %w",
			options.RootPath,
			err,
		)
	}

	// If the provided path is a file,
	// check only that file.
	if !rootInfo.IsDir() {
		if shouldInclude(rootInfo, options) &&
			matcher.matches(rootInfo.Name()) {

			result := createResult(
				options.RootPath,
				rootInfo,
			)

			results = append(results, result)
		}

		return results, nil
	}

	// Search inside the directory.
	err = searchDirectory(
		options.RootPath,
		options,
		matcher,
		&results,
	)

	if err != nil {
		return results, err
	}

	// Sort the results alphabetically by path.
	sort.Slice(results, func(i int, j int) bool {
		return results[i].Path < results[j].Path
	})

	return results, nil
}

// searchDirectory searches inside one directory.
func searchDirectory(
	currentPath string,
	options Options,
	matcher *fileMatcher,
	results *[]Result,
) error {
	entries, err := os.ReadDir(currentPath)

	if err != nil {
		return fmt.Errorf(
			"could not read directory %s: %w",
			currentPath,
			err,
		)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(
			currentPath,
			entry.Name(),
		)

		fileInfo, err := entry.Info()

		if err != nil {
			return fmt.Errorf(
				"could not read information for %s: %w",
				fullPath,
				err,
			)
		}

		// Check whether this item should be included
		// and whether its name matches the pattern.
		if shouldInclude(fileInfo, options) &&
			matcher.matches(entry.Name()) {

			result := createResult(
				fullPath,
				fileInfo,
			)

			*results = append(
				*results,
				result,
			)
		}

		// Search inside child directories.
		if entry.IsDir() && options.Recursive {
			err = searchDirectory(
				fullPath,
				options,
				matcher,
				results,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// shouldInclude checks whether the item type is allowed.
func shouldInclude(
	fileInfo os.FileInfo,
	options Options,
) bool {
	if options.FilesOnly && fileInfo.IsDir() {
		return false
	}

	if options.DirsOnly && !fileInfo.IsDir() {
		return false
	}

	return true
}

// createResult creates a search result.
func createResult(
	path string,
	fileInfo os.FileInfo,
) Result {
	return Result{
		Path:        path,
		Name:        fileInfo.Name(),
		IsDirectory: fileInfo.IsDir(),
		Size:        fileInfo.Size(),
	}
}