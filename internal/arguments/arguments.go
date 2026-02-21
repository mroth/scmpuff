// Package arguments contains shared functions for expanding
// numeric arguments to the associated file names
package arguments

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	expandArgDigitMatcher = regexp.MustCompile("^[0-9]{0,4}$")
	expandArgRangeMatcher = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
	managedEnvVar         = regexp.MustCompile(`^\$e\d+$`)
)

// EvaluateEnvironment evaluates a single arguments and expands environment
// variables.
//
// For scmpuff-managed position variables only (e.g. $e1, etc), the variable is
// expanded into a locatable file path, and expandRelative is true, it will be
// converted into a relative path when possible.
func EvaluateEnvironment(arg string, expandRelative bool) string {
	expanded := os.ExpandEnv(arg)
	wasChanged := (expanded != arg)
	if wasChanged && expandRelative && managedEnvVar.MatchString(arg) {
		relPath, err := convertToRelativeIfFilePath(expanded)
		if err == nil {
			return relPath
		}
	}
	return expanded
}

// For a given arg, try to determine if it represents a file, and if so, convert
// it to a relative filepath.
//
// Otherwise (or if any error conditions occur) return it unmolested.
func convertToRelativeIfFilePath(arg string) (string, error) {
	_, err := os.Stat(arg)
	if err != nil {
		return arg, err
	}
	wd, err := os.Getwd()
	if err != nil {
		return arg, err
	}
	relPath, err := filepath.Rel(wd, arg)
	if err != nil {
		return arg, err
	}
	return relPath, nil
}

// skipExpansion reports whether args[pos] should be left unexpanded.
//
// Normally, Expand treats any bare integer token as a file shortcut (e.g. "1"
// becomes "$e1"). This is wrong when the integer is actually the value of a
// preceding flag — for example, "git log -n 1" means "show one commit", but
// naive expansion turns it into "git log -n $e1" which resolves to a filepath.
// Similarly, "git checkout -b 713" is a branch name, not a file shortcut.
//
// We fix this by maintaining an allowlist of flags (per git subcommand) whose
// space-separated values should never be expanded. Only the subcommands routed
// through "scmpuff exec" by the shell wrapper need coverage, and only flags
// that accept a space-separated value that could be numeric. Flags where the
// value must be glued (e.g. "-U3") are not a problem because the combined
// token doesn't match the digit regex.
//
// gitCmd is the value of SCMPUFF_GIT_CMD; when args[0] matches it, we know
// this is a git command routed through the shell wrapper.
func skipExpansion(args []string, pos int, gitCmd string) bool {
	if pos < 2 || gitCmd == "" || args[0] != gitCmd {
		return false
	}
	prev := args[pos-1]
	switch args[1] {
	case "log":
		switch prev {
		case "-n", "--max-count", "--skip", "--min-parents", "--max-parents":
			return true
		}
	case "checkout":
		switch prev {
		case "-b", "-B", "--orphan":
			return true
		}
	case "blame":
		switch prev {
		case "-L":
			return true
		}
	case "rebase":
		switch prev {
		case "-C":
			return true
		}
	}
	return false
}

// Expand takes the list of arguments received from the command line and expands
// them given our special case rules.
//
// It handles converting numeric file placeholders and range placeholders into
// environment variable symbolic representation,
func Expand(args []string) []string {
	gitCmd := os.Getenv("SCMPUFF_GIT_CMD")
	var results []string
	for i, arg := range args {
		if skipExpansion(args, i, gitCmd) {
			results = append(results, arg)
		} else {
			results = append(results, expandArg(arg)...)
		}
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
