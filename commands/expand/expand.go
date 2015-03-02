package expand

import (
	"fmt"
	"os"
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
		"TODO: DESCRIPTION HERE YO",
	)

	return expandCmd
}

var expandArgDigitMatcher = regexp.MustCompile("^[0-9]{0,4}$")
var expandArgRangeMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
var shellEscaper = regexp.MustCompile("([\\^()\\[\\]<>' \"])")

// Process expands args and performs all substitution, etc.
func Process(args []string) string {
	return escape(evaluateEnvironment(expand(args)))
}

// Escape everything so it can be interpreted once passed along to the shell.
func escape(sequence string) string {
	return shellEscaper.ReplaceAllString(sequence, "\\$1")
}

// Evaluates a string of arguments and expands environment variables.
func evaluateEnvironment(argstr string) string {
	return os.ExpandEnv(argstr)
}

// Expand takes the list of arguments received from the command line and expands
// them given our special case rules.
//
// It handles converting numeric file placeholders and range placeholders into
// environment variable symbolic representation,
//
// Ends up with a final string that is TAB delineated between arguments.
func expand(args []string) string {
	var results []string
	for _, arg := range args {
		results = append(results, expandArg(arg)...)
	}

	return strings.Join(results, "\t")
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
