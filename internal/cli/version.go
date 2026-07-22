package cli

import "fmt"

//current application version.
func (app *App) showVersion() {
	fmt.Fprintln(app.Writer,"Go File Explorer CLI version",app.Version)
}
