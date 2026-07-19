//go:build windows

package fileinfo

import (
	"os"
	"syscall"
	"time"
)

// getCreatedTime reads the Windows file creation time.
func getCreatedTime(
	fileInfo os.FileInfo,
) (time.Time, bool) {
	fileData, ok := fileInfo.Sys().(
		*syscall.Win32FileAttributeData)

	if !ok {
		return time.Time{}, false
	}

	createdTime := time.Unix(
		0,
		fileData.CreationTime.Nanoseconds(),
	)

	return createdTime, true
}