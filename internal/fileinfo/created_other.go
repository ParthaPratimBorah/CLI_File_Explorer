//go:build !windows

package fileinfo

import (
	"os"
	"time"
)

// Creation time is not provided consistently across all operating systems by the Go standard library.
func getCreatedTime(
	fileInfo os.FileInfo,
) (time.Time, bool) {
	return time.Time{}, false
}
