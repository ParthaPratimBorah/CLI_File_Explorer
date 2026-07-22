package cli

import (
	"flag"
	"fmt"
	"io"

	comparepackage "file-explorer/internal/compare"
)

// handles the compare command.
func (app *App) runCompareCommand(args []string) int {
	var recursive bool
	var showHidden bool

	flagSet := flag.NewFlagSet(
		"compare",
		flag.ContinueOnError,
	)

	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar(
		&recursive,
		"recursive",
		true,
		"compare files inside subdirectories",
	)

	flagSet.BoolVar(
		&showHidden,
		"hidden",
		false,
		"include hidden files and directories",
	)

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: invalid compare flags:",
			err,
		)

		return 1
	}

	remainingArguments := flagSet.Args()

	if len(remainingArguments) != 2 {
		fmt.Fprintln(
			app.Writer,
			"Usage: explorer compare [flags] <first-folder> <second-folder>",
		)

		return 1
	}

	firstPath := remainingArguments[0]
	secondPath := remainingArguments[1]

	options := comparepackage.Options{
		FirstPath: firstPath,
		SecondPath: secondPath,
		Recursive: recursive,
		ShowHidden: showHidden,
	}

	if app.Verbose {
		fmt.Fprintln(
			app.Writer,
			"First directory:",
			firstPath,
		)

		fmt.Fprintln(
			app.Writer,
			"Second directory:",
			secondPath,
		)

		app.Logger.Println(
			"Comparing directories:",
			firstPath,
			secondPath,
		)
	}

	result, err := comparepackage.CompareDirectories(
		options,
	)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println(
			"Directory comparison error:",
			err,
		)

		return 1
	}

	app.printComparisonResult(result)

	return 0
}

//displays the comparison
func (app *App) printComparisonResult(
	result comparepackage.Result,
) {
	fmt.Fprintln(
		app.Writer,
		"Directory Comparison Report",
	)


	fmt.Fprintln(app.Writer)

	app.printMissingFiles(result.MissingFiles)
	app.printExtraFiles(result.ExtraFiles)
	app.printModifiedFiles(result.ModifiedFiles)
	app.printSizeDifferences(result.DifferentSizes)
	app.printHashDifferences(result.DifferentHashes)

	totalDifferences :=
		len(result.MissingFiles) +
			len(result.ExtraFiles) +
			len(result.ModifiedFiles)

	fmt.Fprintln(app.Writer)

	if totalDifferences == 0 {
		fmt.Fprintln(app.Writer, "The directories contain matching files.")
		return
	}

	fmt.Fprintf(
		app.Writer,
		"Missing files:  %d\n",
		len(result.MissingFiles),
	)

	fmt.Fprintf(
		app.Writer,
		"Extra files:    %d\n",
		len(result.ExtraFiles),
	)

	fmt.Fprintf(
		app.Writer,
		"Modified files: %d\n",
		len(result.ModifiedFiles),
	)
}

//displays missing files.
func (app *App) printMissingFiles(
	files []comparepackage.FileDetails,
) {
	fmt.Fprintln(app.Writer, "Missing Files")

	if len(files) == 0 {
		fmt.Fprintln(app.Writer, "None")
		fmt.Fprintln(app.Writer)
		return
	}

	for _, file := range files {
		fmt.Fprintf(
			app.Writer,
			"- %s\n",
			file.RelativePath,
		)
	}

	fmt.Fprintln(app.Writer)
}

//displays extra files.
func (app *App) printExtraFiles(
	files []comparepackage.FileDetails,
) {
	fmt.Fprintln(app.Writer, "Extra Files")

	if len(files) == 0 {
		fmt.Fprintln(app.Writer, "None")
		fmt.Fprintln(app.Writer)
		return
	}

	for _, file := range files {
		fmt.Fprintf(
			app.Writer,
			"- %s\n",
			file.RelativePath,
		)
	}

	fmt.Fprintln(app.Writer)
}

//displays modified files
func (app *App) printModifiedFiles(files []string) {
	fmt.Fprintln(app.Writer, "Modified Files")

	if len(files) == 0 {
		fmt.Fprintln(app.Writer, "None")
		fmt.Fprintln(app.Writer)
		return
	}

	for _, file := range files {
		fmt.Fprintf(
			app.Writer,
			"- %s\n",
			file,
		)
	}

	fmt.Fprintln(app.Writer)
}

//displays size changes
func (app *App) printSizeDifferences(
	differences []comparepackage.SizeDifference,
) {
	fmt.Fprintln(app.Writer, "Different Sizes")

	if len(differences) == 0 {
		fmt.Fprintln(app.Writer, "None")
		fmt.Fprintln(app.Writer)
		return
	}

	for _, difference := range differences {
		fmt.Fprintf(app.Writer, "- %s\n",difference.RelativePath)

		fmt.Fprintf(
			app.Writer,
			"  First:  %d bytes (%s)\n",
			difference.FirstSize,
			comparepackage.FormatSize(difference.FirstSize),
		)

		fmt.Fprintf(
			app.Writer,
			"  Second: %d bytes (%s)\n",
			difference.SecondSize,
			comparepackage.FormatSize(difference.SecondSize),
		)
	}

	fmt.Fprintln(app.Writer)
}

// displays hash changes
func (app *App) printHashDifferences(
	differences []comparepackage.HashDifference,
) {
	fmt.Fprintln(app.Writer, "Different Hashes")

	if len(differences) == 0 {
		fmt.Fprintln(app.Writer, "None")
		fmt.Fprintln(app.Writer)
		return
	}

	for _, difference := range differences {
		fmt.Fprintf(
			app.Writer,
			"- %s\n",
			difference.RelativePath,
		)

		fmt.Fprintf(app.Writer, "  First SHA256:  %s\n", difference.FirstHash)
		fmt.Fprintf(app.Writer, "  Second SHA256: %s\n", difference.SecondHash)
	}

	fmt.Fprintln(app.Writer)
}