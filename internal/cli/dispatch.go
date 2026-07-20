package cli

import "fmt"

// dispatch decides which command should run.
func (app *App) dispatch(command string, args []string) int {
	switch command {
	case "tree":
		return app.runTreeCommand(args)

	case "search":
		return app.runSearchCommand(args)

	case "copy":
		return app.runCopyCommand(args)

	case "move":
		return app.runMoveCommand(args)

	case "rename":
		return app.runRenameCommand(args)

	case "delete":
		return app.runDeleteCommand(args)

	case "info":
		return app.runInfoCommand(args)

	case "duplicate":
		return app.runDuplicateCommand(args)

	case "compare",
		"stats",
		"export",
		"hash",
		"encode",
		"decode":

		fmt.Fprintf(
			app.Writer,
			"The %s command is not implemented yet.\n",
			command,
		)

		return 0

	case "help":
		app.showHelp()
		return 0

	case "version":
		app.showVersion()
		return 0

	default:
		fmt.Fprintf(
			app.Writer,
			"Unknown command: %s\n\n",
			command,
		)

		app.showHelp()

		return 1
	}
}