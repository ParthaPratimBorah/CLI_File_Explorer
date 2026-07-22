package compare

import (
	"fmt"
	"sort"
)

// CompareDirectories compares two directories.
func CompareDirectories(options Options) (Result, error) {
	var result Result

	if options.FirstPath == "" {
		return result, fmt.Errorf(
			"first directory path cannot be empty",
		)
	}

	if options.SecondPath == "" {
		return result, fmt.Errorf(
			"second directory path cannot be empty",
		)
	}

	firstFiles, err := collectFiles(
		options.FirstPath,
		options,
	)

	if err != nil {
		return result, err
	}

	secondFiles, err := collectFiles(
		options.SecondPath,
		options,
	)

	if err != nil {
		return result, err
	}

	// Check files from the first directory.
	for relativePath, firstFile := range firstFiles {
		secondFile, exists := secondFiles[relativePath]

		// File exists in first directory but not second.
		if !exists {
			result.MissingFiles = append(
				result.MissingFiles,
				firstFile,
			)

			continue
		}

		// If sizes are different, the files are modified.
		if firstFile.Size != secondFile.Size {
			result.DifferentSizes = append(
				result.DifferentSizes,
				SizeDifference{
					RelativePath: relativePath,
					FirstSize:    firstFile.Size,
					SecondSize:   secondFile.Size,
				},
			)

			result.ModifiedFiles = append(
				result.ModifiedFiles,
				relativePath,
			)

			continue
		}

		// Sizes are equal, so calculate hashes.
		firstHash, err := calculateSHA256(
			firstFile.AbsolutePath,
		)

		if err != nil {
			return result, err
		}

		secondHash, err := calculateSHA256(
			secondFile.AbsolutePath,
		)

		if err != nil {
			return result, err
		}

		// Same size but different hash means
		// the contents are different.
		if firstHash != secondHash {
			result.DifferentHashes = append(
				result.DifferentHashes,
				HashDifference{
					RelativePath: relativePath,
					FirstHash:    firstHash,
					SecondHash:   secondHash,
				},
			)

			result.ModifiedFiles = append(
				result.ModifiedFiles,
				relativePath,
			)
		}
	}

	// Check files available only in the second directory.
	for relativePath, secondFile := range secondFiles {
		_, exists := firstFiles[relativePath]

		if !exists {
			result.ExtraFiles = append(
				result.ExtraFiles,
				secondFile,
			)
		}
	}

	sortFileDetails(result.MissingFiles)
	sortFileDetails(result.ExtraFiles)

	sort.Strings(result.ModifiedFiles)

	sort.Slice(
		result.DifferentSizes,
		func(i int, j int) bool {
			return result.DifferentSizes[i].RelativePath <
				result.DifferentSizes[j].RelativePath
		},
	)

	sort.Slice(
		result.DifferentHashes,
		func(i int, j int) bool {
			return result.DifferentHashes[i].RelativePath <
				result.DifferentHashes[j].RelativePath
		},
	)

	return result, nil
}