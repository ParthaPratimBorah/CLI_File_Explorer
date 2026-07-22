package cli

import (
	"flag"
	"fmt"
	"io"
)

// global flag stored here
type GlobalOptions struct {
	Help      bool
	Version   bool
	Verbose   bool
	Recursive bool
	Output    string
	LogFile   string
}

// flags that appear before the command
func parseGlobalFlags(args []string) (GlobalOptions, []string, error) {
	var options GlobalOptions

	flagSet := flag.NewFlagSet("explorer", flag.ContinueOnError)

	// Prevent the flag package from printing its own error message
	flagSet.SetOutput(io.Discard)

	flagSet.BoolVar( &options.Help, "help", false, "show help" )
	flagSet.BoolVar( &options.Help, "h", false, "show help" )
	flagSet.BoolVar( &options.Version, "version", false, "show version" )
	flagSet.BoolVar( &options.Verbose, "verbose", false, "show detailed information" )
	flagSet.BoolVar( &options.Recursive, "recursive", false, "process directories recursively" )
	flagSet.StringVar( &options.Output, "output", "", "write output to a file" )
	flagSet.StringVar( &options.LogFile, "log", "", "write logs to a file" )
	
	err := flagSet.Parse(args)

	if err != nil {
		return options, nil, fmt.Errorf("invalid global flag: %w", err)
	}

	//returns arguments that were not global flags.
	return options, flagSet.Args(), nil
}
