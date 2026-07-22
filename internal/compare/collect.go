package compare

import (
	"fmt"
	"os"
	"path/filepath"
)

//collects regular files from a directory using map
func collectFiles(
	rootPath string,
	options Options,
) (map[string]FileDetails, error) {
	files := make(map[string]FileDetails)

	rootInfo, err := os.Stat(rootPath)

	if err != nil {
		return files, fmt.Errorf(
			"could not access directory %s: %w",
			rootPath,
			err,
		)
	}

	if !rootInfo.IsDir() {
		return files, fmt.Errorf(
			"%s is not a directory",
			rootPath,
		)
	}

	err = collectFromDirectory(
		rootPath,
		rootPath,
		options,
		files,
	)

	if err != nil {
		return files, err
	}

	return files, nil
}

// collectFromDirectory reads one directory.
func collectFromDirectory(
	rootPath string,
	currentPath string,
	options Options,
	files map[string]FileDetails,
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
		if !options.ShowHidden &&
			isHidden(entry.Name()) {

			continue
		}

		fullPath := filepath.Join(
			currentPath,
			entry.Name(),
		)

		if entry.IsDir() {
			if options.Recursive {
				err = collectFromDirectory(
					rootPath,
					fullPath,
					options,
					files,
				)

				if err != nil {
					return err
				}
			}

			continue
		}

		fileInfo, err := entry.Info()

		if err != nil {
			return fmt.Errorf(
				"could not read file information for %s: %w",
				fullPath,
				err,
			)
		}

		// Ignore symbolic links and special files.
		if !fileInfo.Mode().IsRegular() {
			continue
		}

		relativePath, err := filepath.Rel(
			rootPath,
			fullPath,
		)

		if err != nil {
			return fmt.Errorf(
				"could not create relative path for %s: %w",
				fullPath,
				err,
			)
		}

		relativePath = normalizeRelativePath(
			relativePath,
		)

		absolutePath, err := filepath.Abs(
			fullPath,
		)

		if err != nil {
			return fmt.Errorf(
				"could not create absolute path for %s: %w",
				fullPath,
				err,
			)
		}

		files[relativePath] = FileDetails{
			RelativePath: relativePath,
			AbsolutePath: absolutePath,
			Size:         fileInfo.Size(),
		}
	}

	return nil
}