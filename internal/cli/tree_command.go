package cli

import (
	"flag"
	"fmt"
	"io"

	treepackage "file-explorer/internal/tree"
)

// runTreeCommand handles the tree command.
func (app *App) runTreeCommand(args []string) int {
	var maxDepth int
	var showHidden bool
	var showSize bool

	flagSet := flag.NewFlagSet("tree", flag.ContinueOnError)

	// Hide the default flag package output.
	flagSet.SetOutput(io.Discard)

	flagSet.IntVar(
		&maxDepth,
		"depth",
		0,
		"maximum tree depth",
	)

	flagSet.BoolVar(
		&showHidden,
		"hidden",
		false,
		"show hidden files",
	)

	flagSet.BoolVar(
		&showSize,
		"size",
		false,
		"show file sizes",
	)

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: invalid tree command flags:",
			err,
		)

		return 1
	}

	// Use the current directory when no path is provided.
	rootPath := "."

	remainingArguments := flagSet.Args()

	if len(remainingArguments) > 0 {
		rootPath = remainingArguments[0]
	}

	if maxDepth < 0 {
		fmt.Fprintln(
			app.Writer,
			"Error: depth cannot be negative",
		)

		return 1
	}

	options := treepackage.Options{
		MaxDepth:   maxDepth,
		ShowHidden: showHidden,
		ShowSize:   showSize,
	}

	if app.Verbose {
		fmt.Fprintln(app.Writer, "Tree root:", rootPath)
		app.Logger.Println("Tree root:", rootPath)
	}

	result, err := treepackage.PrintTree(
		rootPath,
		options,
		app.Writer,
	)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Tree error:", err)

		return 1
	}

	fmt.Fprintln(app.Writer)

	fmt.Fprintf(
		app.Writer,
		"Directories: %d\n",
		result.DirectoryCount,
	)

	fmt.Fprintf(
		app.Writer,
		"Files: %d\n",
		result.FileCount,
	)

	return 0
}
