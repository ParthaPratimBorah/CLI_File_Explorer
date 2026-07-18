package main

import (
	"os"
	"file-explorer/internal/cli"
)

func main() {
	// notun CLI app
	app := cli.NewApp()
	exitCode := app.Run(os.Args[1:])		// os.Args[0] = Program name aru  os.Args[1:] = arguments entered

	// End BY exit code
	os.Exit(exitCode)
}
