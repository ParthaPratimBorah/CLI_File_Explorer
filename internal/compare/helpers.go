package compare

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//checks whether a name begins with a dot.
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

//converts bytes into a readable size.
func formatSize(size int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
		gigabyte = 1024 * megabyte
	)

	switch {
	case size >= gigabyte:
		return fmt.Sprintf( "%.2f GB", float64(size)/float64(gigabyte))

	case size >= megabyte:
		return fmt.Sprintf( "%.2f MB", float64(size)/float64(megabyte))

	case size >= kilobyte:
		return fmt.Sprintf( "%.2f KB", float64(size)/float64(kilobyte) )

	default:
		return fmt.Sprintf("%d B", size)
	}
}

//cleans a relative path.
func normalizeRelativePath(path string) string {
	path = filepath.Clean(path)
	return filepath.ToSlash(path)
}

// sorts files by relative path
func sortFileDetails(files []FileDetails) {
	sort.Slice(
		files,
		func(i int, j int) bool {
			return files[i].RelativePath < files[j].RelativePath
		},
	)
}

//allows the CLI package to display readable sizes.
func FormatSize(size int64) string {
	return formatSize(size)
}