package expand

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mroth/scmpuff/commands/arguments"
	"github.com/spf13/cobra"
)

var expandRelative bool

// CommandExpand generates the command handler for `scmpuff expand`
//
// Allows expansion of numbered shortcuts, ranges of shortcuts, or standard paths.
// Numbered shortcut variables are produced by various commands, such as:
//
//  * scmpuff_status()  - git status implementation
func CommandExpand() *cobra.Command {

	var expandCmd = &cobra.Command{
		Use:   "expand <shortcuts...>",
		Short: "Expands numeric shortcuts",
		Long: `Expands numeric shortcuts to their full filepath.

Takes a list of digits (1 4 5) or numeric ranges (1-5) or even both.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Usage()
			}

			fmt.Print(Process(args))
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

var shellEscaper = regexp.MustCompile("([\\^()\\[\\]<>' \";\\|*])")

// Process expands args and performs all substitution, etc.
//
// Ends up with a final string that is TAB delineated between arguments.
func Process(args []string) string {
	var processedArgs []string
	for _, arg := range arguments.Expand(args) {
		processed := escape(arguments.EvaluateEnvironment(arg, expandRelative))

		// if we still ended up with a totally blank arg, escape it here.
		// we handle this as a special case rather than in expandArg because we
		// don't want it to be subject to normal escaping.
		if processed == "" {
			processed = "''"
		}

		processedArgs = append(processedArgs, processed)
	}

	return strings.Join(processedArgs, "\t")
}

// Escape everything so it can be interpreted once passed along to the shell.
func escape(arg string) string {
	return shellEscaper.ReplaceAllString(arg, "\\$1")
}

