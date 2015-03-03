package expand

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mroth/scmpuff/helpers"
	"github.com/spf13/cobra"
)

var expandRelative bool

// CommandExpand generates the command handler for `scmpuff expand`
//
// Allows expansion of numbered shortcuts, ranges of shortcuts, or standard paths.
// Numbered shortcut variables are produced by various commands, such as:
//
//  * git_status_shortcuts()  - git status implementation
func CommandExpand() *cobra.Command {

	var expandCmd = &cobra.Command{
		Use:   "expand",
		Short: "Expands numbered shortcuts",
		Long: `
LONG DESCRIPTION HERE
    `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Process(args))
		},
	}

	// --relative
	expandCmd.Flags().BoolVarP(
		&expandRelative,
		"relative",
		"r",
		false,
		"try to expand path relative to current working directory",
	)

	return expandCmd
}

var expandArgDigitMatcher = regexp.MustCompile("^[0-9]{0,4}$")
var expandArgRangeMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
var shellEscaper = regexp.MustCompile("([\\^()\\[\\]<>' \"])")

// Process expands args and performs all substitution, etc.
//
// Ends up with a final string that is TAB delineated between arguments.
func Process(args []string) string {
	expandedArgs := evaluateEnvironment(expand(args))
	sequence := strings.Join(expandedArgs, "\t")
	return escape(sequence)
}

// Escape everything so it can be interpreted once passed along to the shell.
func escape(sequence string) string {
	return shellEscaper.ReplaceAllString(sequence, "\\$1")
}

// Evaluates a string of arguments and expands environment variables.
func evaluateEnvironment(args []string) []string {
	var results []string
	for _, arg := range args {
		expandedArg := os.ExpandEnv(arg)
		if expandRelative {
			results = append(results, convertToRelativeIfFilePath(expandedArg))
		} else {
			results = append(results, expandedArg)
		}
	}
	return results
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
		// TODO: check if is numerically named file or directory, dont expand if so
		// return "$e" + dm
		digit, _ := strconv.Atoi(dm)
		result := helpers.IntToEnvVar(digit)
		return []string{result}
	}

	// ...is it a range?
	rm := expandArgRangeMatcher.FindStringSubmatch(arg)
	if rm != nil {
		lo, _ := strconv.Atoi(rm[1])
		hi, _ := strconv.Atoi(rm[2])

		var results []string
		for i := lo; i <= hi; i++ {
			results = append(results, helpers.IntToEnvVar(i))
		}
		return results
	}

	// if it was neither, return as-is
	return []string{arg}
}
