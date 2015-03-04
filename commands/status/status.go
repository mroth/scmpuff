package status

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// CommandStatus processes 'git status --porcelain', and exports numbered
// env variables that contain the path of each affected file.
// Output is also more concise than standard 'git status'.
//
// TODO: Call with optional <group> parameter to filter by modification state:
// 1 || Staged,  2 || Unmerged,  3 || Unstaged,  4 || Untracked
func CommandStatus() *cobra.Command {
	var optsFilelist bool
	var optsDisplay bool

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Set and display numbered git status",
		Long: `
Processes 'git status --porcelain', and exports numbered env variables that
contain the path of each affected file.

The output is prettier and more concise than standard 'git status'.

In most cases, you won't want to call this directly, but rather will be using
the exported shell-function 'scmpuff_status', which wraps this command and also
sets the environment variables for your shell. (For more information on this,
see 'scmpuff init'.)
    `,
		Run: func(cmd *cobra.Command, args []string) {
			root := gitProjectRoot()
			status := gitStatusOutput()

			results := Process(status, root)
			results.printStatus(optsFilelist, optsDisplay)
		},
	}

	// --filelist, -f
	statusCmd.Flags().BoolVarP(
		&optsFilelist,
		"filelist", "f", false,
		"include machine-parseable filelist",
	)

	// --display
	// allow normal display to be disabled, not really useful unless you know you
	// JUST want the machine parseable file-list for some reason.
	statusCmd.Flags().BoolVarP(
		&optsDisplay,
		"display", "", true,
		"displays the formatted status output",
	)

	// --relative
	// statusCmd.Flags().BoolVarP(
	// 	&expandRelative,
	// 	"relative",
	// 	"r",
	// 	false,
	// 	"TODO: SHOULD THIS BE IMPLEMENTED? OR JUST KEEP ALWAYS ON...",
	// )

	return statusCmd
}

// Runs `git status --porcelain` and returns the results.
//
// If an error is encountered, the process will die fatally.
func gitStatusOutput() []byte {
	gso, err := exec.Command("git", "status", "--porcelain", "-b").Output()
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Trim(gso, "\n")
}

// Returns the root for the git project.
//
// If can't be found, the process will die fatally.
func gitProjectRoot() string {
	gpr, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()

	if err != nil {
		// we want to capture and handle status 128 in a pretty way
		if err.Error() == "exit status 128" {
			msg := "Not a git repository (or any of the parent directories)"
			fmt.Fprintf(os.Stderr, "\033[0;31m"+msg+"\n")
			os.Exit(128)
		}
		// or, some other sort of error?
		log.Fatal(err)
	}
	return string(bytes.TrimSpace(gpr))
}
