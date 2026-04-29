package ui

import (
	"fmt"

	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

// Confirm prompts the user for confirmation before a destructive action.
// Returns true if confirmed, false if cancelled.
// Skips the prompt and returns true if skipConfirm is set or stdout is not a TTY.
func Confirm(io *iostreams.IOStreams, prompt string, skipConfirm bool) bool {
	if skipConfirm || !io.IsStdoutTTY() {
		return true
	}
	fmt.Fprintf(io.ErrOut, "%s [Y/n] ", prompt)
	var answer string
	_, _ = fmt.Fscanln(io.In, &answer)
	return answer == "" || answer == "y" || answer == "Y" || answer == "yes"
}
