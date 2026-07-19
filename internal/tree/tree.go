package tree

import (
	"fmt"
	"io"
	"os"
)

// prints the folder structure
func PrintTree(rootPath string, options Options, writer io.Writer) (Result, error) {
	var result Result

	// Check whether the path exists
	fileInfo, err := os.Stat(rootPath)

	if err != nil {
		return result, fmt.Errorf("could not access path %s: %w", rootPath, err)
	}

	// The tree command must start with a directory
	if !fileInfo.IsDir() {
		return result, fmt.Errorf("%s is not a directory", rootPath)
	}

	rootName := getRootName(rootPath)

	// Print the root folder name.
	fmt.Fprintln(writer, rootName)

	err = walkDirectory(rootPath, "", 1, options, writer, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}
