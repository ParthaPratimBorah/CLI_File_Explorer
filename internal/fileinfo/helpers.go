package fileinfo

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// formatSize converts bytes into a readable value.
func formatSize(size int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
		gigabyte = 1024 * megabyte
	)

	switch {
	case size >= gigabyte:
		return fmt.Sprintf(
			"%.2f GB",
			float64(size)/float64(gigabyte),
		)

	case size >= megabyte:
		return fmt.Sprintf(
			"%.2f MB",
			float64(size)/float64(megabyte),
		)

	case size >= kilobyte:
		return fmt.Sprintf(
			"%.2f KB",
			float64(size)/float64(kilobyte),
		)

	default:
		return fmt.Sprintf("%d B", size)
	}
}

// getFileType returns a simple item type.
func getFileType(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()

	switch {
	case mode.IsDir():
		return "Directory"

	case mode.IsRegular():
		return "Regular file"

	case mode&os.ModeSymlink != 0:
		return "Symbolic link"

	case mode&os.ModeNamedPipe != 0:
		return "Named pipe"

	case mode&os.ModeSocket != 0:
		return "Socket"

	default:
		return "Other"
	}
}

// isHidden checks names beginning with a dot.
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// calculateChecksum calculates a SHA-256 checksum.
//
// A checksum is only calculated for regular files.
func calculateChecksum(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", fmt.Errorf(
			"could not open file for checksum: %w",
			err,
		)
	}

	defer file.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, file)

	if err != nil {
		return "", fmt.Errorf(
			"could not calculate checksum: %w",
			err,
		)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// getRelativePath calculates the path from the current directory.
func getRelativePath(absolutePath string) string {
	currentDirectory, err := os.Getwd()

	if err != nil {
		return absolutePath
	}

	relativePath, err := filepath.Rel(
		currentDirectory,
		absolutePath,
	)

	if err != nil {
		return absolutePath
	}

	return relativePath
}
