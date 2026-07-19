package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copy copies a file or directory.
func Copy(
	sourcePath string,
	destinationPath string,
	recursive bool,
	overwrite bool,
) error {
	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not access source %s: %w",
			sourcePath,
			err,
		)
	}

	destinationPath = PrepareDestination(
		sourcePath,
		destinationPath,
		sourceInfo.IsDir(),
	)

	err = validateDifferentPaths(
		sourcePath,
		destinationPath,
	)

	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		if !recursive {
			return fmt.Errorf(
				"source is a directory; use --recursive",
			)
		}

		err = validateDirectoryDestination(
			sourcePath,
			destinationPath,
		)

		if err != nil {
			return err
		}
	}

	if PathExists(destinationPath) && !overwrite {
		return fmt.Errorf(
			"destination already exists: %s",
			destinationPath,
		)
	}

	if sourceInfo.IsDir() {
		return copyDirectory(
			sourcePath,
			destinationPath,
			overwrite,
		)
	}

	return copyFile(
		sourcePath,
		destinationPath,
		overwrite,
	)
}

// copyFile copies one file.
func copyFile(
	sourcePath string,
	destinationPath string,
	overwrite bool,
) error {
	if PathExists(destinationPath) {
		if !overwrite {
			return fmt.Errorf(
				"destination already exists: %s",
				destinationPath,
			)
		}

		err := RemoveExistingDestination(destinationPath)

		if err != nil {
			return err
		}
	}

	sourceFile, err := os.Open(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not open source file %s: %w",
			sourcePath,
			err,
		)
	}

	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()

	if err != nil {
		return fmt.Errorf(
			"could not read source file information: %w",
			err,
		)
	}

	// Create the destination's parent folder if needed.
	parentDirectory := filepath.Dir(destinationPath)

	err = os.MkdirAll(parentDirectory, 0755)

	if err != nil {
		return fmt.Errorf(
			"could not create destination directory %s: %w",
			parentDirectory,
			err,
		)
	}

	destinationFile, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		sourceInfo.Mode(),
	)

	if err != nil {
		return fmt.Errorf(
			"could not create destination file %s: %w",
			destinationPath,
			err,
		)
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)

	if err != nil {
		return fmt.Errorf(
			"could not copy file data: %w",
			err,
		)
	}

	return nil
}

// copyDirectory copies a directory and all its contents
func copyDirectory(
	sourcePath string,
	destinationPath string,
	overwrite bool,
) error {
	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not read source directory: %w",
			err,
		)
	}

	// Create the destination directory
	err = os.MkdirAll(
		destinationPath,
		sourceInfo.Mode(),
	)

	if err != nil {
		return fmt.Errorf(
			"could not create destination directory %s: %w",
			destinationPath,
			err,
		)
	}

	entries, err := os.ReadDir(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not read directory %s: %w",
			sourcePath,
			err,
		)
	}

	for _, entry := range entries {
		childSource := filepath.Join(sourcePath, entry.Name())

		childDestination := filepath.Join(destinationPath, entry.Name())

		if entry.IsDir() {
			err = copyDirectory(childSource, childDestination, overwrite)
		} else {
			err = copyFile(
				childSource,
				childDestination,
				overwrite,
			)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
