package duplicate

import (
	"fmt"
	"strings"
)

// check koribole if . di start hua files ase niki
func isHidden (fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}


//size to readable format anibole

func formatSize(size int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
		gigabyte = 1024 * megabyte
	)

	switch {
		case size >= gigabyte:
			return fmt.Sprintf("%.2f GB", float64(size)/float64(gigabyte))
		case size >= megabyte:
			return fmt.Sprintf("%.2f MB", float64(size)/float64(megabyte))
		case size >= kilobyte:
			return fmt.Sprintf("%.2f KB", float64(size)/float64(kilobyte))
		default:
			return fmt.Sprintf("%d B", size)
	}	
}