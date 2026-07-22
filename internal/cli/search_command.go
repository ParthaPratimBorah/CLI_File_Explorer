package cli

import (
	"flag"
	"fmt"
	"io"

	searchpackage "file-explorer/internal/search"
)

//handles the search command
func (app *App) runSearchCommand(args []string) int {
	var exact bool
	var ignoreCase bool
	var useRegex bool
	var extension bool
	var recursive bool
	var filesOnly bool
	var directoriesOnly bool

	flagSet := flag.NewFlagSet("search", flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar( &exact, "exact", false, "search exact names" )
	flagSet.BoolVar( &ignoreCase, "ignore-case", false, "ignore uppercase and lowercase" )
	flagSet.BoolVar( &useRegex, "regex", false, "use regular expression" )
	flagSet.BoolVar( &extension, "extension", false, "search by extension" )
	flagSet.BoolVar( &recursive, "recursive", true, "search subdirectories" )
	flagSet.BoolVar( &filesOnly, "files-only", false, "show only files" )
	flagSet.BoolVar( &directoriesOnly, "dirs-only", false, "show only directories" )

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln( app.Writer, "Error: invalid search command flags:", err)

		return 1
	}

	remainingArguments := flagSet.Args()

	// At least one argument is required for the search pattern.
	if len(remainingArguments) == 0 {
		fmt.Fprintln( app.Writer, "Error: search pattern is required")

		fmt.Fprintln( app.Writer, "Usage: explorer search [flags] <pattern> [path]" )

		return 1
	}

	if filesOnly && directoriesOnly {
		fmt.Fprintln( app.Writer, "Error: --files-only and --dirs-only cannot be used together" )

		return 1
	}

	searchPattern := remainingArguments[0]

	// Use the current directory by default.
	rootPath := "."

	if len(remainingArguments) > 1 {
		rootPath = remainingArguments[1]
	}

	options := searchpackage.Options{
		Pattern: searchPattern,
		RootPath: rootPath,
		Exact: exact,
		IgnoreCase: ignoreCase,
		UseRegex: useRegex,
		Extension: extension,
		Recursive: recursive,
		FilesOnly: filesOnly,
		DirsOnly: directoriesOnly,
	}

	if app.Verbose {
		fmt.Fprintln(app.Writer, "Search pattern:", searchPattern)
		fmt.Fprintln(app.Writer, "Search location:", rootPath)

		app.Logger.Println("Search pattern:", searchPattern)
		app.Logger.Println("Search location:", rootPath)
	}

	results, err := searchpackage.Search(options)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Search error:", err)

		return 1
	}

	if len(results) == 0 {
		fmt.Fprintln(app.Writer, "No matches found.")
		return 0
	}

	fmt.Fprintln(app.Writer, "Search results:")

	for _, result := range results {
		if result.IsDirectory {
			fmt.Fprintf( app.Writer, "[DIR]  %s\n", result.Path )
		} else {
			fmt.Fprintf( app.Writer, "[FILE] %s\n", result.Path )
		}
	}

	fmt.Fprintln(app.Writer)

	fmt.Fprintf(
		app.Writer,
		"Total matches: %d\n",
		len(results),
	)

	return 0
}
