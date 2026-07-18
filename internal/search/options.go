package search

//SEARCH
type Options struct {
	Pattern string
	RootPath string
	Exact bool
	IgnoreCase bool
	UseRegex bool
	Extension bool
	Recursive bool
	FilesOnly bool
	DirsOnly bool
}
