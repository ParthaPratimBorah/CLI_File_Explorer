package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"file-explorer/internal/fileops"
)

//handles the copy command
func (app *App) runCopyCommand(args []string) int {
	var recursive bool
	var overwrite bool

	flagSet := flag.NewFlagSet("copy", flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar(&recursive, "recursive", false, "copy a directory recursively")

	flagSet.BoolVar( &overwrite, "overwrite", false, "overwrite an existing destination")

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(app.Writer,"Error: invalid copy flags:",err)
		return 1
	}

	remainingArguments := flagSet.Args()

	if len(remainingArguments) != 2 {
		fmt.Fprintln(app.Writer, "Usage: explorer copy [flags] <source> <destination>")
		return 1
	}

	sourcePath := remainingArguments[0]
	destinationPath := remainingArguments[1]

	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		fmt.Fprintln(app.Writer,"Error: could not access source:",err)

		return 1
	}

	finalDestination := fileops.PrepareDestination(sourcePath,destinationPath,sourceInfo.IsDir())

	// Ask before replacing an existing destination
	if fileops.PathExists(finalDestination) && !overwrite {
		confirmed := askForConfirmation(
			fmt.Sprintf("%s already exists. Overwrite it?",finalDestination),
		)

		if !confirmed {
			fmt.Fprintln(app.Writer,"Copy cancelled.")
			return 0
		}

		overwrite = true
	}
	err = fileops.Copy(sourcePath,destinationPath,recursive,overwrite)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Copy error:", err)

		return 1
	}

	fmt.Fprintln(app.Writer,"Copy completed successfully.")
	app.Logger.Printf("Copied %s to %s",sourcePath,destinationPath)

	return 0
}
