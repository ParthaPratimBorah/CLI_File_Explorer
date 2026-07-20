package fileinfo

import (
	"os"
	"syscall"
	"time"
)

//reads the creation time on Windows.
func getCreatedTime(fileInfo os.FileInfo) (time.Time, bool) {
	fileData, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)

	if !ok {
		return time.Time{}, false
	}

	createdTime := time.Unix(0,fileData.CreationTime.Nanoseconds())

	return createdTime, true
}