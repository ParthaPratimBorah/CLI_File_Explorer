package fileinfo

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetInfo reads information about a file or directory.
func GetInfo(path string) (Result, error) {
	var result Result

	fileInfo, err := os.Stat(path)

	if err != nil {
		return result, fmt.Errorf(
			"could not access %s: %w",
			path,
			err,
		)
	}

	absolutePath, err := filepath.Abs(path)

	if err != nil {
		return result, fmt.Errorf(
			"could not create absolute path: %w",
			err,
		)
	}

	createdTime, createdKnown := getCreatedTime(fileInfo)

	result = Result{
		FileName: fileInfo.Name(),
		Extension: filepath.Ext(fileInfo.Name()),
		AbsolutePath: absolutePath,
		RelativePath: getRelativePath(absolutePath),
		Directory: filepath.Dir(absolutePath),
		Size: fileInfo.Size(),
		ReadableSize: formatSize(fileInfo.Size()),
		Permissions: fileInfo.Mode().Perm().String(),
		CreatedTime: createdTime,
		CreatedKnown: createdKnown,
		ModifiedTime: fileInfo.ModTime(),
		FileType: getFileType(fileInfo),
		Hidden: isHidden(fileInfo.Name()),
		Checksum: "Not available",
	}

	// Only calculate checksums for regular files.
	if fileInfo.Mode().IsRegular() {
		checksum, checksumError := calculateChecksum(path)

		if checksumError != nil {
			return result, checksumError
		}

		result.Checksum = checksum
	}

	return result, nil
}
