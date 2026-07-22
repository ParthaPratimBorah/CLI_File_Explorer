package stats

//stores one file extension and its count.
type ExtensionCount struct {
	Extension string
	Count int
}

//stores one directory and its total size
type DirectorySize struct {
	Path string
	RelativePath string
	Size int64
	ReadableSize string
}

//stores complete statistics report
type Result struct {
	RootPath string

	TotalFiles int
	TotalDirectories int

	TotalStorage int64
	ReadableTotalStorage string

	AverageSize float64
	ReadableAverageSi