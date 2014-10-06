package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var expandRelative bool

// CommandExpand generates the command handler for `scmpuff expand`
//
// Allows expansion of numbered shortcuts, ranges of shortcuts, or standard paths.
// Numbered shortcut variables are produced by various commands, such as:
//
//  * git_status_shortcuts()  - git status implementation
//  * git_show_affected_files() - shows files affected by a given SHA1, etc.
func CommandExpand() *cobra.Command {

	var expandCmd = &cobra.Command{
		Use:   "expand",
		Short: "Expands numbered shortcuts",
		Long: `LONG DESCRIPTION HERE
    `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(EvaluateArgs(ExpandArgs(args)))
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

// Evaluates the args and substitutes in environment variables
func EvaluateArgs(argstr string) string {
	return os.ExpandEnv(argstr)
}

// ExpandArgs provides expansion of numbered shortcuts, ranges of shortcuts,
// or standard paths.
func ExpandArgs(args []string) string {
	// split on spaces to loop over expanded range
	// args := strings.Split(argstr, " ")

	var results []string
	for _, arg := range args {
		results = append(results, ExpandArg(arg))
	}

	return strings.Join(results, " ")
}

// TODO: figure out how to turn these into precompiled constants?
var expandArgDigitMatcher = regexp.MustCompile("^[0-9]{0,4}$")
var expandArgRangeMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")

// ExpandArg expands a single argument.
func ExpandArg(arg string) string {
	// ...is it a single digit?
	dm := expandArgDigitMatcher.FindString(arg)
	if dm != "" {
		// TODO: check if is numerically named file or directory, dont expand if so
		return "$e" + dm
	}

	// is it a range?
	rm := expandArgRangeMatcher.FindStringSubmatch(arg)
	if rm != nil {
		// convert to ints, shouldn't be able to fail if regex worked properly
		lo, _ := strconv.Atoi(rm[1])
		hi, _ := strconv.Atoi(rm[2])

		var results []string
		for i := lo; i <= hi; i++ {
			results = append(results, intToEnvVar(i))
		}

		return strings.Join(results, " ")
	}

	// if can do nothing else, return it as is
	return arg
}

func intToEnvVar(num int) string {
	return "$e" + strconv.Itoa(num)
}
