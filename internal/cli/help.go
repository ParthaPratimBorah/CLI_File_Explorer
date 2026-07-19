package cli

import "fmt"

// showHelp prints the application help message.
func (app *App) showHelp() {
	fmt.Fprintln(app.Writer, "Go File Explorer CLI")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Usage:")
	fmt.Fprintln(app.Writer, "  explorer [global flags] <command> [command flags]")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Available Commands:")
	fmt.Fprintln(app.Writer, "  tree       Display a directory tree")
	fmt.Fprintln(app.Writer, "  search     Search files and directories")
	fmt.Fprintln(app.Writer, "  copy       Copy files or directories")
	fmt.Fprintln(app.Writer, "  move       Move files or directories")
	fmt.Fprintln(app.Writer, "  rename     Rename files or directories")
	fmt.Fprintln(app.Writer, "  delete     Delete a file or directory")
	fmt.Fprintln(app.Writer, "  compare    Compare two files")
	fmt.Fprintln(app.Writer, "  duplicate  Find duplicate files")
	fmt.Fprintln(app.Writer, "  stats      Display filesystem statistics")
	fmt.Fprintln(app.Writer, "  export     Export filesystem information")
	fmt.Fprintln(app.Writer, "  hash       Calculate a file hash")
	fmt.Fprintln(app.Writer, "  encode     Encode file data")
	fmt.Fprintln(app.Writer, "  decode     Decode file data")
	fmt.Fprintln(app.Writer, "  info       Display file information")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Global Flags:")
	fmt.Fprintln(app.Writer, "  --help         Show help")
	fmt.Fprintln(app.Writer, "  --version      Show application version")
	fmt.Fprintln(app.Writer, "  --verbose      Show detailed information")
	fmt.Fprintln(app.Writer, "  --recursive    Enable recursive processing")
	fmt.Fprintln(app.Writer, "  --output FILE  Write output to a file")
	fmt.Fprintln(app.Writer, "  --log FILE     Write logs to a file")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Tree Flags:")
	fmt.Fprintln(app.Writer, "  --depth NUMBER  Limit directory depth")
	fmt.Fprintln(app.Writer, "  --hidden        Show hidden files")
	fmt.Fprintln(app.Writer, "  --size          Show file sizes")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Search Flags:")
	fmt.Fprintln(app.Writer, "  --exact         Search exact names")
	fmt.Fprintln(app.Writer, "  --ignore-case   Ignore uppercase and lowercase")
	fmt.Fprintln(app.Writer, "  --regex         Use regular expressions")
	fmt.Fprintln(app.Writer, "  --extension     Search by file extension")
	fmt.Fprintln(app.Writer, "  --recursive     Search inside subdirectories")
	fmt.Fprintln(app.Writer, "  --files-only    Show only files")
	fmt.Fprintln(app.Writer, "  --dirs-only     Show only directories")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Copy Flags:")
	fmt.Fprintln(app.Writer, "  --recursive     Copy a directory recursively")
	fmt.Fprintln(app.Writer, "  --overwrite     Replace an existing destination")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Move Flags:")
	fmt.Fprintln(app.Writer, "  --overwrite     Replace an existing destination")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Rename Flags:")
	fmt.Fprintln(app.Writer, "  --new-name NAME  Rename one item")
	fmt.Fprintln(app.Writer, "  --batch          Enable batch rename")
	fmt.Fprintln(app.Writer, "  --prefix TEXT    Add text before names")
	fmt.Fprintln(app.Writer, "  --suffix TEXT    Add text after names")
	fmt.Fprintln(app.Writer, "  --replace TEXT   Text to replace")
	fmt.Fprintln(app.Writer, "  --with TEXT      Replacement text")
	fmt.Fprintln(app.Writer, "  --regex PATTERN  Regex rename pattern")
	fmt.Fprintln(app.Writer, "  --overwrite      Replace an existing item")
	fmt.Fprintln(app.Writer, "  --force          Skip batch confirmation")
	fmt.Fprintln(app.Writer)

	fmt.Fprintln(app.Writer, "Delete Flags:")
	fmt.Fprintln(app.Writer, "  --recursive     Delete a directory and its contents")
	fmt.Fprintln(app.Writer, "  --force         Delete without confirmation")
}
