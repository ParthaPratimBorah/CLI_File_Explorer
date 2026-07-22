package cli

import (
	"flag"
	"fmt"
	"io"

	"file-explorer/internal/fileinfo"
)

// handles the info command
func (app *App) runInfoCommand(args []string) int {
	flagSet := flag.NewFlagSet("info", flag.ContinueOnError)

	flagSet.SetOutput(io.Discard)

	err := flagSet.Parse(args)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error: invalid info flags:", err)

		return 1
	}

	remainingArguments := flagSet.Args()

	if len(remainingArguments) != 1 {
		fmt.Fprintln(app.Writer, "Usage: explorer info <path>")

		return 1
	}

	path := remainingArguments[0]

	result, err := fileinfo.GetInfo(path)

	if err != nil {
		fmt.Fprintln(app.Writer, "Error:", err)
		app.Logger.Println("Info error:", err)

		return 1
	}

	fmt.Fprintln( app.Writer, "File Information" )

	fmt.Fprintf( app.Writer, "File Name:      %s\n", result.FileName )

	if result.Extension == "" {
		fmt.Fprintln(app.Writer,"Extension: No extension")
	} else {
		fmt.Fprintf(app.Writer,"Extension: %s\n",result.Extension)
	}

	fmt.Fprintf( app.Writer, "Absolute Path: %s\n", result.AbsolutePath)
	fmt.Fprintf( app.Writer, "Relative Path: %s\n", result.RelativePath)
	fmt.Fprintf( app.Writer, "Directory: %s\n", result.Directory )
	fmt.Fprintf( app.Writer, "Size: %d bytes\n", result.Size )
	fmt.Fprintf( app.Writer, "Readable Size: %s\n", result.ReadableSize)
	fmt.Fprintf( app.Writer, "Permissions: %s\n", result.Permissions)

	if result.CreatedKnown {
		fmt.Fprintf( app.Writer, "Created Time: %s\n", result.CreatedTime.Format("2006-01-02 15:04:05") )
	} else {
		fmt.Fprintln( app.Writer, "Created Time:   Not available" )
	}

	fmt.Fprintf( app.Writer, "Modified Time: %s\n", result.ModifiedTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf( app.Writer, "File Type: %s\n", result.FileType )
	fmt.Fprintf( app.Writer, "Hidden: %t\n", result.Hidden )
	fmt.Fprintf( app.Writer, "Checksum: %s\n", result.Checksum)

	return 0
}
