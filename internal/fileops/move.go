package fileops

import (
	"fmt"
	"os"
)

// Move moves a file or directory.

func Move(
	sourcePath string,
	destinationPath string,
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
		err = validateDirectoryDestination(
			sourcePath,
			destinationPath,
		)

		if err != nil {
			return err
		}
	}

	if PathExists(destinationPath) {
		if !overwrite {
			return fmt.Errorf(
				"destination already exists: %s",
				destinationPath,
			)
		}

		err = RemoveExistingDestination(destinationPath)

		if err != nil {
			return err
		}
	}

	// First try the normal operating system move.
	err = os.Rename(sourcePath, destinationPath)

	if err == nil {
		return nil
	}

	// os.Rename may fail when moving between two drives.
	// In that case, copy the source and then delete the original.
	copyError := Copy(
		sourcePath,
		destinationPath,
		true,
		true,
	)

	if copyError != nil {
		return fmt.Errorf(
			"could not move item: %w",
			copyError,
		)
	}

	deleteError := Delete(sourcePath, true)

	if deleteError != nil {
		return fmt.Errorf(
			"item was copied but the original could not be deleted: %w",
			deleteError,
		)
	}

	return nil
}
