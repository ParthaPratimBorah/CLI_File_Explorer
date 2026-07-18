package cli

import (
	"fmt"
	"io"
	"log"
	"os"
)

// information store koribole
type App struct {
	Version   string
	Verbose   bool
	Recursive bool
	Writer io.Writer
	Logger *log.Logger //log message dibole 

	// these files are stored so that we can close them later
	outputFile *os.File
	logFile    *os.File
}

// creates and returns a new cli app
func NewApp() *App {
	return &App{
		Version: "1.0.0",
		Writer: os.Stdout,		// default output in the terminal
		Logger: log.New(io.Discard, "", log.LstdFlags), 	// default logging is disabled
	}
}

// start koribole func
func (app *App) Run(args []string) int {
	// read the global flags
	globalOptions, remainingArguments, err := parseGlobalFlags(args)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return 1
	}

	app.Verbose = globalOptions.Verbose
	app.Recursive = globalOptions.Recursive

	// set up output file when --output is provided
	err = app.setupOutput(globalOptions.Output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return 1
	}

	// set up log file when --log is provided
	err = app.setupLogger(globalOptions.LogFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		app.closeFiles()
		return 1
	}

	// difer use kori opened file kita close koribole last t
	defer app.closeFiles()

	if globalOptions.Help {
		app.showHelp()
		return 0
	}

	if globalOptions.Version {
		app.showVersion()
		return 0
	}

	// show help jodi eku dia nai
	if len(remainingArguments) == 0 {
		app.showHelp()
		return 0
	}
	command := remainingArguments[0]		// first remaining argument = command
	commandArguments := remainingArguments[1:]		// the rest belongs to commandArguments

	if app.Verbose {
		fmt.Fprintln(app.Writer, "Running command:", command)
		app.Logger.Println("Running command:", command)
	}

	return app.dispatch(command, commandArguments)
}

// setupOutput changes the normal output from terminal to a file
func (app *App) setupOutput(filename string) error {
	if filename == "" {
		return nil
	}

	file, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}

	app.outputFile = file
	app.Writer = file

	return nil
}

// setupLogger creates a log file
func (app *App) setupLogger(filename string) error {
	if filename == "" {
		return nil
	}

	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return fmt.Errorf("could not open log file: %w", err)
	}

	app.logFile = file

	app.Logger = log.New(
		file,
		"FILE-EXPLORER: ",
		log.Ldate|log.Ltime,
	)

	return nil
}

// closeFiles closes output and log files if they were opened
func (app *App) closeFiles() {
	if app.outputFile != nil {
		app.outputFile.Close()
	}

	if app.logFile != nil {
		app.logFile.Close()
	}
}
