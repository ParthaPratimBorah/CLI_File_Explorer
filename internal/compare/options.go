package compare

//stores directory comparison settings
type Options struct {
	FirstPath string
	SecondPath string
	Recursive bool
	ShowHidden bool
}