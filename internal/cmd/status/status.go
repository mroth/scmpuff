package status

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mroth/scmpuff/internal/gitstatus/porcelainv1"
	"github.com/spf13/cobra"
)

var optsFilelist bool
var optsDisplay bool

// NewStatusCmd creates and returns the status command
func NewStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // silence usage-on-error after args processed

			// Obtain the current working directory (needed to determine git root and relative paths)
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("fatal: failed to retrieve current working directory: %w", err)
			}

			// Determine the git project root
			root, err := gitProjectRoot(wd)
			if err != nil {
				// Error 128 is a common case when running outside of a git repo,
				// and is a common UX situation rather than an actual error, so
				// handle it directly here rather than returning to cobra.
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) && exitErr.ExitCode() == 128 {
					msg := "Not a git repository (or any of the parent directories)"
					fmt.Fprintf(os.Stderr, "%s%s%s\n", RedColor, msg, ResetColor)
					os.Exit(128)
				}

				// or, some sort of an *actual* error
				return fmt.Errorf("fatal: failed to determine git project root: %w", err)
			}

			// Run the git status command to get the porcelain output
			status, err := gitStatusOutput()
			if err != nil {
				return fmt.Errorf("fatal: error running git status command: %w", err)
			}

			// Parse and then process the git status output.
			info, err := porcelainv1.Process(status)
			if err != nil {
				// Currently, this is a "special case" error condition, as it means there was a failure
				// processing the git status output, which is somewhere we definitely want to obtain a debug
				// dump for. We print a message to stderr and exit directly here, rather than returning an
				// error to cobra, as we want to provide special instructions.
				//
				// NOTE: In the future, we might want to consider a more general "alert error" mechanism
				// that can be returned by any process and scmpuff handles the UI.
				fmt.Fprintf(os.Stderr, `Error: failed to process git status output: %v

Please file a bug including this error message as well as the output from:

scmpuff debug dump --archive

You can file the bug at: https://github.com/mroth/scmpuff/issues/`, err)
				os.Exit(1)
			}

			// Render the formatted status output
			renderer, err := NewRenderer(info, root, wd)
			if err != nil {
				return fmt.Errorf("fatal: failed to create status renderer: %w", err)
			}

			if err := renderer.Display(os.Stdout, optsFilelist, optsDisplay); err != nil {
				return fmt.Errorf("fatal: failed to render status: %w", err)
			}

			return nil
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

	return statusCmd
}

// Runs `git status --porcelain=v1 -b -z` and returns the results.
//
// Why z-mode? It lets us do machine parsing in a reliable cross-platform way,
// as per this quote from the git status documentation:
//
//	"There is also an alternate -z format recommended for machine parsing. In
//	that format, the status field is the same, but some other things change.
//	First, the -> is omitted from rename entries and the field order is
//	reversed (e.g from -> to becomes to from). Second, a NUL (ASCII 0) follows
//	each filename, replacing space as a field separator and the terminating
//	newline (but a space still separates the status field from the first
//	filename). Third, filenames containing special characters are not specially
//	formatted; no quoting or backslash-escaping is performed."
//
// Okay, it also introduces some complexity because it wasn't well thought out,
// but it beats dealing with shell escaping and hoping we do it right across
// different platforms and shells, I hope...  see `process.go` for all the parsing
// we do to make sense of it, this just grabs its output.
//
// NOTE: More recent versions of git support `--porcelain=v2`, which is a more
// reasoned structured output format addressing mistakes in git porcelain, but
// we have not yet implemented support for that.
func gitStatusOutput() ([]byte, error) {
	// We actually use `git status -z -b` here, which is the same as `git status
	// --porcelain=v1 -b -z`, as the `-z` flag implies `--porcelain=v1` if not
	// specified otherwise. That way we retain reverse compatiblity with very
	// old versions of git that might not understand the `--porcelain=v1` flag
	// (prior to porcelain v2 support, the flag was just `--porcelain` and was
	// still implied by `-z`).
	return exec.Command("git", "status", "-z", "-b").Output()
}

// Runs git commands to determine the root for the git project.
//
// This handles relative paths within a symlink'd directory correctly,
// which was previously broken as described in:
// https://github.com/mroth/scmpuff/issues/11
//
// Requires knowing the current working directory.
//
// See https://github.com/mroth/scmpuff/pull/94
//
// Note that there is a common 'error' condition when running this command
// outside of a git repository, which is an os/exec.exitError with status code
// 128. Callers of this function should handle that error gracefully, as it is a
// common UX situation, and not an actual error in the program.
func gitProjectRoot(wd string) (string, error) {
	// `--show-cdup` prints the relative path to the Git repository root,
	// which we then join with the current working directory.
	cdup, err := exec.Command("git", "rev-parse", "--show-cdup").Output()
	if err != nil {
		return "", err
	}

	absPath := filepath.Join(wd, string(bytes.TrimSpace(cdup)))
	return filepath.Clean(absPath), nil
}
