package cli

import (
	"flag"
	"fmt"
	"io"

	"file-explorer/internal/duplicate"
)

//handles the duplicate command.
func (app *App) runDuplicateCommand(args []string) int {
	var algorithm string
	var recursive bool
	var showHidden bool

	flagSet := flag.NewFlagSet( "duplicate", flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	flagSet.StringVar( &algorithm, "algorithm", "sha256", "hash algorithm: sha256 or crc32")
	flagSet.BoolVar( &recursive, "recursive", true, "search inside subdirectories")
	flagSet.BoolVar( &showHidden, "hidden", false, "include hidden files")

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln( app.Writer, "Error: invalid duplicate flags:", err)

		return 1
	}

	remainingArguments := flagSet.Args()

	// Search the current directory by default.
	rootPath := "."

	if len(remainingArguments) > 0 {
		rootPath = remainingArguments[0]
	}

	if len(remainingArguments) > 1 {
		fmt.Fprintln( app.Writer, "Usage: explorer duplicate [flags] [path]")

		return 1
	}

	options := duplicate.Options{
		RootPath:   rootPath,
		Algorithm:  algorithm,
		Recursive:  recursive,
		ShowHidden: showHidden,
	}

	if app.Verbose {
		fmt.Fprintln( app.Writer, "Duplicate search path:", rootPath)

		fmt.Fprintln( app.Writer, "Hash algorithm:", algorithm)
	}

	groups, err := duplicate.Find(options)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Duplicate finder error:", err)

		return 1
	}

	if len(groups) == 0 {
		fmt.Fprintln( app.Writer, "No duplicate files found.")

		return 0
	}

	fmt.Fprintln( app.Writer, "Duplicate Files")

	fmt.Fprintln(app.Writer)

	for index, group := range groups {
		fmt.Fprintf( app.Writer, "Duplicate Group %d\n", index+1)

		fmt.Fprintf( app.Writer, "Algorithm: %s\n", algorithm)

		fmt.Fprintf( app.Writer, "Hash: %s\n", group.Hash )

		fmt.Fprintf( app.Writer, "Size: %d bytes (%s)\n", group.Size, group.ReadableSize)

		fmt.Fprintln( app.Writer, "Locations:" )

		for _, location := range group.Locations {
			fmt.Fprintf( app.Writer, "  - %s\n", location)
		}

		fmt.Fprintln(app.Writer)
	}

	fmt.Fprintf( app.Writer, "Total duplicate groups: %d\n", len(groups))

	return 0
}