package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"file-explorer/internal/fileops"
)

// runMoveCommand handles the move command.
func (app *App) runMoveCommand(args []string) int {
	var overwrite bool

	flagSet := flag.NewFlagSet("move",flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar(
		&overwrite,
		"overwrite",
		false,
		"overwrite an existing destination",
	)

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: invalid move flags:",
			err,
		)

		return 1
	}

	remainingArguments := flagSet.Args()

	if len(remainingArguments) != 2 {
		fmt.Fprintln(
			app.Writer,
			"Usage: explorer move [flags] <source> <destination>",
		)

		return 1
	}

	sourcePath := remainingArguments[0]
	destinationPath := remainingArguments[1]

	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: could not access source:",
			err,
		)

		return 1
	}

	finalDestination := fileops.PrepareDestination(sourcePath,destinationPath,sourceInfo.IsDir())

	if fileops.PathExists(finalDestination) && !overwrite {
		confirmed := askForConfirmation(
			fmt.Sprintf("%s already exists. Overwrite it?",finalDestination),
		)

		if !confirmed {
			fmt.Fprintln(app.Writer,"Move cancelled.")

			return 0
		}

		overwrite = true
	}

	err = fileops.Move(sourcePath, destinationPath, overwrite)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Move error:", err)

		return 1
	}

	fmt.Fprintln(app.Writer,"Move completed successfully.")

	app.Logger.Printf("Moved %s to %s", sourcePath, destinationPath)

	return 0
}
