package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathExists checks whether a file or directory exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// PrepareDestination decides the final destination path.
func PrepareDestination(
	sourcePath string,
	destinationPath string,
	sourceIsDirectory bool,
) string {
	destinationInfo, err := os.Stat(destinationPath)

	if err != nil {
		return destinationPath
	}

	// For a file, copy or move it inside an existing directory.
	if destinationInfo.IsDir() && !sourceIsDirectory {
		return filepath.Join(
			destinationPath,
			filepath.Base(sourcePath),
		)
	}

	return destinationPath
}

// RemoveExistingDestination removes an existing destination.
func RemoveExistingDestination(path string) error {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf(
			"could not check destination %s: %w",
			path,
			err,
		)
	}

	if fileInfo.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err != nil {
		return fmt.Errorf(
			"could not remove existing destination %s: %w",
			path,
			err,
		)
	}

	return nil
}

// validateDifferentPaths checks that source and destination are not the same path.
func validateDifferentPaths(
	sourcePath string,
	destinationPath string,
) error {
	absoluteSource, err := filepath.Abs(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not read source path: %w",
			err,
		)
	}

	absoluteDestination, err := filepath.Abs(destinationPath)

	if err != nil {
		return fmt.Errorf(
			"could not read destination path: %w",
			err,
		)
	}

	absoluteSource = filepath.Clean(absoluteSource)
	absoluteDestination = filepath.Clean(absoluteDestination)

	if absoluteSource == absoluteDestination {
		return fmt.Errorf(
			"source and destination cannot be the same",
		)
	}

	return nil
}

// validateDirectoryDestination prevents copying or moving a directory inside itself.
func validateDirectoryDestination(
	sourcePath string,
	destinationPath string,
) error {
	absoluteSource, err := filepath.Abs(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not read source path: %w",
			err,
		)
	}

	absoluteDestination, err := filepath.Abs(destinationPath)

	if err != nil {
		return fmt.Errorf(
			"could not read destination path: %w",
			err,
		)
	}

	absoluteSource = filepath.Clean(absoluteSource)
	absoluteDestination = filepath.Clean(absoluteDestination)

	sourceWithSeparator := absoluteSource +
		string(os.PathSeparator)

	if strings.HasPrefix(
		absoluteDestination,
		sourceWithSeparator,
	) {
		return fmt.Errorf(
			"cannot place a directory inside itself",
		)
	}

	return nil
}