package stats

import "time"

//stores information about one file
type FileRecord struct {
	Path         string
	RelativePath string
	Name         string
	Extension    string
	Size         int64
	ModifiedTime time.Time
}