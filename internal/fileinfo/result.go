package fileinfo

import "time"

// stores information about one file or directory
type Result struct {
	FileName     string
	Extension    string
	AbsolutePath string
	RelativePath string
	Directory    string
	Size         int64
	ReadableSize string
	Permissions  string
	CreatedTime  time.Time
	CreatedKnown bool
	ModifiedTime time.Time
	FileType     string
	Hidden       bool
	Checksum     string
}
