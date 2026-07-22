package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"file-explorer/internal/fileops"
)

//handles the delete command.
func (app *App) runDeleteCommand(args []string) int {
	var recursive bool
	var force bool

	flagSet := flag.NewFlagSet( "delete", flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar( &recursive, "recursive", false, "delete a directory and all its contents")
	flagSet.BoolVar( &force, "force", false, "delete without confirmation")

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(app.Writer,"Error: invalid delete flags:",err)

		return 1
	}

	remainingArguments := flagSet.Args()

	if len(remainingArguments) != 1 {
		fmt.Fprintln(app.Writer,"Usage: explorer delete [flags] <path>")

		return 1
	}

	targetPath := remainingArguments[0]

	fileInfo, err := os.Stat(targetPath)

	if err != nil {
		fmt.Fprintln(app.Writer,"Error: could not access path:",err)

		return 1
	}

	if fileInfo.IsDir() && !recursive {
		fmt.Fprintln(app.Writer,"Note: a non-empty directory requires --recursive.")
	}

	if !force {
		message := fmt.Sprintf("Are you sure you want to delete %s?",targetPath)

		confirmed := askForConfirmation(message)

		if !confirmed {
			fmt.Fprintln(app.Writer,"Delete cancelled.")

			return 0
		}
	}

	err = fileops.Delete(targetPath,recursive)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Delete error:", err)

		return 1
	}

	fmt.Fprintln(app.Writer,"Delete completed successfully.")

	app.Logger.Println("Deleted:", targetPath)

	return 0
}
