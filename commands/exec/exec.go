package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var expandRelative bool

// CommandExec expands numeric arguments then execs the command
//
// Allows expansion of numbered shortcuts, ranges of shortcuts, or standard paths.
// Numbered shortcut variables are produced by various commands, such as:
//
//  * scmpuff_status()  - git status implementation
func CommandExec() *cobra.Command {

	var expandCmd = &cobra.Command{
		Use:   "exec <shortcuts...>",
		Short: "Execute cmd with numeric shortcuts",
		Long: `Expands numeric shortcuts to their full filepath and executes the command.

Takes a list of digits (1 4 5) or numeric ranges (1-5) or even both.`,
		Run: func(cmd *cobra.Command, inputArgs []string) {
			if len(inputArgs) < 1 {
				cmd.Usage()
				os.Exit(1)
			}

			expandedArgs := Process(inputArgs)
			a := expandedArgs[1:]
			subcmd := exec.Command(expandedArgs[0], a...)
			subcmd.Stdin = os.Stdin
			subcmd.Stdout = os.Stdout
			subcmd.Stderr = os.Stderr
			err := subcmd.Run()
			if err == nil {
				os.Exit(0)
			}
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	// --relative
	expandCmd.Flags().BoolVarP(
		&expandRelative,
		"relative",
		"r",
		false,
		"make path relative to current working directory",
	)

	return expandCmd
}

var expandArgDigitMatcher = regexp.MustCompile("^[0-9]{0,4}$")
var expandArgRangeMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")

// Process expands args and performs all substitution, then returns the argument array
func Process(args []string) []string {
	var processedArgs []string
	for _, arg := range expand(args) {
		processed := evaluateEnvironment(arg)
		processedArgs = append(processedArgs, processed)
	}

	return processedArgs
}

// Evaluates a string of arguments and expands environment variables.
func evaluateEnvironment(arg string) string {
	expandedArg := os.ExpandEnv(arg)
	if expandRelative {
		return convertToRelativeIfFilePath(expandedArg)
	}
	return expandedArg
}

// For a given arg, try to determine if it represents a file, and if so, convert
// it to a relative filepath.
//
// Otherwise (or if any error conditions occur) return it unmolested.
func convertToRelativeIfFilePath(arg string) string {
	if _, err := os.Stat(arg); err == nil {
		wd, err1 := os.Getwd()
		relPath, err2 := filepath.Rel(wd, arg)
		if err1 == nil && err2 == nil {
			return relPath
		}
	}
	return arg
}

// Expand takes the list of arguments received from the command line and expands
// them given our special case rules.
//
// It handles converting numeric file placeholders and range placeholders into
// environment variable symbolic representation,
func expand(args []string) []string {
	var results []string
	for _, arg := range args {
		results = append(results, expandArg(arg)...)
	}
	return results
}

// expandArg "expands" a single argument we received on the command line.
//
// It's possible that argument represents a numeric file placeholder, in which
// case we will replace it with the syntax to represent the environment variable
// that it will be held in (e.g. "$e1").
//
// It's also possible that argument may represent a range, in which case it will
// return multiple instances of environment variable placeholders.
func expandArg(arg string) []string {

	// ...is it a single digit?
	dm := expandArgDigitMatcher.FindString(arg)
	if dm != "" {
		// dont expand if its actually a numerically named file or directory!
		if _, err := os.Stat(dm); err == nil {
			return []string{arg} //return as-is
		}

		result := "$e" + dm
		return []string{result}
	}

	// ...is it a range?
	rm := expandArgRangeMatcher.FindStringSubmatch(arg)
	if rm != nil {
		lo, _ := strconv.Atoi(rm[1])
		hi, _ := strconv.Atoi(rm[2])

		var results []string
		for i := lo; i <= hi; i++ {
			results = append(results, "$e"+strconv.Itoa(i))
		}
		return results
	}

	// if it was neither, return as-is
	return []string{arg}
}
