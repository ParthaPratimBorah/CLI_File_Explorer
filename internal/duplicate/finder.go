package duplicate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Find searches for duplicate files.
func Find(options Options) ([]Group, error) {
	var duplicateGroups []Group

	// Use SHA256 when no algorithm
	if options.Algorithm == "" {
		options.Algorithm = "sha256"
	}

	options.Algorithm = strings.ToLower(
		options.Algorithm,
	)

	if options.Algorithm != "sha256" &&
		options.Algorithm != "crc32" {

		return duplicateGroups, fmt.Errorf(
			"algorithm must be sha256 or crc32",
		)
	}

	filePaths, err := collectFiles(options)

	if err != nil {
		return duplicateGroups, err
	}

	// First group the files using their sizes.
	filesBySize, err := groupFilesBySize(filePaths)

	if err != nil {
		return duplicateGroups, err
	}

	// Only files with the same size need hashing.
	for fileSize, sameSizeFiles := range filesBySize {
		if len(sameSizeFiles) < 2 {
			continue
		}

		filesByHash := make(map[string][]string)

		for _, filePath := range sameSizeFiles {
			hashValue, err := calculateHash(
				filePath,
				options.Algorithm,
			)

			if err != nil {
				return duplicateGroups, err
			}

			filesByHash[hashValue] = append(
				filesByHash[hashValue],
				filePath,
			)
		}

		for hashValue, locations := range filesByHash {
			// A single file is not a duplicate.
			if len(locations) < 2 {
				continue
			}

			sort.Strings(locations)

			group := Group{
				Hash:         hashValue,
				Size:         fileSize,
				ReadableSize: formatSize(fileSize),
				Locations:    locations,
			}

			duplicateGroups = append(
				duplicateGroups,
				group,
			)
		}
	}

	// Sort larger duplicate groups first.
	sort.Slice(
		duplicateGroups,
		func(i int, j int) bool {
			if duplicateGroups[i].Size ==
				duplicateGroups[j].Size {

				return duplicateGroups[i].Hash <
					duplicateGroups[j].Hash
			}

			return duplicateGroups[i].Size >
				duplicateGroups[j].Size
		},
	)

	return duplicateGroups, nil
}

// collectFiles collects regular files from the root path.
func collectFiles(options Options) ([]string, error) {
	var filePaths []string

	rootInfo, err := os.Stat(options.RootPath)

	if err != nil {
		return filePaths, fmt.Errorf(
			"could not access path %s: %w",
			options.RootPath,
			err,
		)
	}

	// If the root is one file, add only that file.
	if !rootInfo.IsDir() {
		if rootInfo.Mode().IsRegular() {
			filePaths = append(
				filePaths,
				options.RootPath,
			)
		}

		return filePaths, nil
	}

	err = collectFromDirectory(
		options.RootPath,
		options,
		&filePaths,
	)

	if err != nil {
		return filePaths, err
	}

	return filePaths, nil
}

// collectFromDirectory reads files from one directory.
func collectFromDirectory(
	directoryPath string,
	options Options,
	filePaths *[]string,
) error {
	entries, err := os.ReadDir(directoryPath)

	if err != nil {
		return fmt.Errorf(
			"could not read directory %s: %w",
			directoryPath,
			err,
		)
	}

	for _, entry := range entries {
		if !options.ShowHidden &&
			isHidden(entry.Name()) {

			continue
		}

		fullPath := filepath.Join(
			directoryPath,
			entry.Name(),
		)

		fileInfo, err := entry.Info()

		if err != nil {
			return fmt.Errorf(
				"could not read file information for %s: %w",
				fullPath,
				err,
			)
		}

		if entry.IsDir() {
			if options.Recursive {
				err = collectFromDirectory(
					fullPath,
					options,
					filePaths,
				)

				if err != nil {
					return err
				}
			}

			continue
		}

		// Ignore symbolic links, sockets and other special items.
		if fileInfo.Mode().IsRegular() {
			*filePaths = append(
				*filePaths,
				fullPath,
			)
		}
	}

	return nil
}

// groupFilesBySize groups files having the same size.
func groupFilesBySize(
	filePaths []string,
) (map[int64][]string, error) {
	filesBySize := make(map[int64][]string)

	for _, filePath := range filePaths {
		fileInfo, err := os.Stat(filePath)

		if err != nil {
			return filesBySize, fmt.Errorf(
				"could not read file information for %s: %w",
				filePath,
				err,
			)
		}

		fileSize := fileInfo.Size()

		filesBySize[fileSize] = append(
			filesBySize[fileSize],
			filePath,
		)
	}

	return filesBySize, nil
}