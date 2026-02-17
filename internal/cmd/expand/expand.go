package expand

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mroth/scmpuff/internal/arguments"
	"github.com/spf13/cobra"
)

// NewExpandCmd creates and returns the expand command
func NewExpandCmd() *cobra.Command {
	var expandRelative bool

	expandCmd := &cobra.Command{
		Use:   "expand [flags] <shortcuts...>",
		Short: "Expands numeric shortcuts",
		Long: `Expands numeric shortcuts to their full filepath.

Takes a list of digits (1 4 5) or numeric ranges (1-5) or even both.`,
		Example: "$ scmpuff expand 1-2\n/tmp/foo.txt    /tmp/bar.txt",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // silence usage-on-error after args processed
			fmt.Print(Process(args, expandRelative))
			return nil
		},
	}

	expandCmd.Flags().BoolVarP(&expandRelative, "relative", "r", false, "make path relative to current working directory")
	return expandCmd
}

var shellEscaper = regexp.MustCompile("([\\^()\\[\\]<>' \";\\|*])")

// Process expands args and performs all substitution, etc.
//
// Ends up with a final string that is TAB delineated between arguments.
func Process(args []string, relative bool) string {
	var processedArgs []string
	for _, arg := range arguments.Expand(args) {
		processed := escape(arguments.EvaluateEnvironment(arg, relative))

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
