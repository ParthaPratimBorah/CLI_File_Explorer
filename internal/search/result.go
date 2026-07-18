package search

// Result stores information about one matched item.
type Result struct {
	Path        string
	Name        string
	IsDirectory bool
	Size        int64
}
