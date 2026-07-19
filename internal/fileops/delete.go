package fileops

import (
	"fmt"
	"os"
)

// Delete deletes a file or directory.
//
// A non-empty directory requires recursive to be true.
func Delete(path string, recursive bool) error {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return fmt.Errorf(
			"could not access %s: %w",
			path,
			err,
		)
	}

	if !fileInfo.IsDir() {
		err = os.Remove(path)

		if err != nil {
			return fmt.Errorf(
				"could not delete file %s: %w",
				path,
				err,
			)
		}

		return nil
	}

	if recursive {
		err = os.RemoveAll(path)

		if err != nil {
			return fmt.Errorf(
				"could not delete directory %s: %w",
				path,
				err,
			)
		}

		return nil
	}

	// os.Remove can delete an empty directory.
	// It returns an error when the directory is not empty.
	err = os.Remove(path)

	if err != nil {
		return fmt.Errorf(
			"could not delete directory %s; use --recursive if it is not empty: %w",
			path,
			err,
		)
	}

	return nil
}
