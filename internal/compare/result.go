package compare

//stores basic information about one file.
type FileDetails struct {
	RelativePath string
	AbsolutePath string
	Size         int64
}

//stores files having different sizes.
type SizeDifference struct {
	RelativePath string
	FirstSize    int64
	SecondSize   int64
}

//stores files having different SHA256 hashes.
type HashDifference struct {
	RelativePath string
	FirstHash    string
	SecondHash   string
}

//stores the complete directory comparison result
type Result struct {
	MissingFiles    []FileDetails
	ExtraFiles      []FileDetails
	ModifiedFiles   []string
	DifferentSizes  []SizeDifference
	DifferentHashes []HashDifference
}