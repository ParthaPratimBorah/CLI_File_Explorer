package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"file-explorer/internal/fileops"
)

// runRenameCommand handles the rename command.
func (app *App) runRenameCommand(args []string) int {
	var newName string
	var batch bool
	var prefix string
	var suffix string
	var replaceText string
	var replacement string
	var regexPattern string
	var overwrite bool
	var force bool

	flagSet := flag.NewFlagSet(
		"rename",
		flag.ContinueOnError,
	)

	flagSet.SetOutput(io.Discard)

	flagSet.StringVar(
		&newName,
		"new-name",
		"",
		"new name for one item",
	)

	flagSet.BoolVar(
		&batch,
		"batch",
		false,
		"rename multiple items",
	)

	flagSet.StringVar(
		&prefix,
		"prefix",
		"",
		"add text before names",
	)

	flagSet.StringVar(
		&suffix,
		"suffix",
		"",
		"add text after names",
	)

	flagSet.StringVar(
		&replaceText,
		"replace",
		"",
		"text to replace",
	)

	flagSet.StringVar(
		&replacement,
		"with",
		"",
		"replacement text",
	)

	flagSet.StringVar(
		&regexPattern,
		"regex",
		"",
		"regular expression pattern",
	)

	flagSet.BoolVar(
		&overwrite,
		"overwrite",
		false,
		"overwrite an existing item during single rename",
	)

	flagSet.BoolVar(
		&force,
		"force",
		false,
		"skip confirmation",
	)

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: invalid rename flags:",
			err,
		)

		return 1
	}

	remainingArguments := flagSet.Args()

	if batch {
		return app.runBatchRename(
			remainingArguments,
			prefix,
			suffix,
			replaceText,
			replacement,
			regexPattern,
			force,
		)
	}

	return app.runSingleRename(
		remainingArguments,
		newName,
		overwrite,
	)
}

// runSingleRename renames one file or directory.
func (app *App) runSingleRename(
	args []string,
	newName string,
	overwrite bool,
) int {
	if len(args) != 1 {
		fmt.Fprintln(
			app.Writer,
			"Usage: explorer rename --new-name <name> <path>",
		)

		return 1
	}

	if newName == "" {
		fmt.Fprintln(
			app.Writer,
			"Error: --new-name is required",
		)

		return 1
	}

	sourcePath := args[0]

	_, err := os.Stat(sourcePath)

	if err != nil {
		fmt.Fprintln(
			app.Writer,
			"Error: could not access source:",
			err,
		)

		return 1
	}

	destinationPath := filepath.Join(
		filepath.Dir(sourcePath),
		newName,
	)

	if fileops.PathExists(destinationPath) && !overwrite {
		confirmed := askForConfirmation(
			fmt.Sprintf(
				"%s already exists. Overwrite it?",
				destinationPath,
			),
		)

		if !confirmed {
			fmt.Fprintln(
				app.Writer,
				"Rename cancelled.",
			)

			return 0
		}

		overwrite = true
	}

	err = fileops.Rename(
		sourcePath,
		newName,
		overwrite,
	)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Rename error:", err)

		return 1
	}

	fmt.Fprintln(
		app.Writer,
		"Rename completed successfully.",
	)

	app.Logger.Printf(
		"Renamed %s to %s",
		sourcePath,
		newName,
	)

	return 0
}

// runBatchRename renames all items inside a directory.
func (app *App) runBatchRename(
	args []string,
	prefix string,
	suffix string,
	replaceText string,
	replacement string,
	regexPattern string,
	force bool,
) int {
	if len(args) != 1 {
		fmt.Fprintln(
			app.Writer,
			"Usage: explorer rename --batch [flags] <directory>",
		)

		return 1
	}

	directoryPath := args[0]

	if !force {
		confirmed := askForConfirmation(
			fmt.Sprintf(
				"Rename items inside %s?",
				directoryPath,
			),
		)

		if !confirmed {
			fmt.Fprintln(
				app.Writer,
				"Batch rename cancelled.",
			)

			return 0
		}
	}

	options := fileops.RenameOptions{
		Prefix:       prefix,
		Suffix:       suffix,
		ReplaceText:  replaceText,
		Replacement:  replacement,
		RegexPattern: regexPattern,
	}

	results, err := fileops.BatchRename(
		directoryPath,
		options,
	)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Batch rename error:", err)

		return 1
	}

	if len(results) == 0 {
		fmt.Fprintln(
			app.Writer,
			"No names were changed.",
		)

		return 0
	}

	for _, result := range results {
		fmt.Fprintf(
			app.Writer,
			"Renamed: %s -> %s\n",
			filepath.Base(result.OldPath),
			filepath.Base(result.NewPath),
		)
	}

	fmt.Fprintf(
		app.Writer,
		"Total renamed: %d\n",
		len(results),
	)

	return 0
}
