package exec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mroth/scmpuff/internal/arguments"
	"github.com/spf13/cobra"
)

var expandRelative bool

// ExecCmd expands numeric arguments then executes the command directly.
//
// Allows expansion of numbered shortcuts, ranges of shortcuts, or standard paths.
var ExecCmd = &cobra.Command{
	Use:     "exec [flags] <command> <shortcuts...>",
	Example: "$ scmpuff exec git add 1-4",
	Aliases: []string{"execute"},
	Short:   "Execute cmd with numeric shortcuts",
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

func init() {
	ExecCmd.Flags().BoolVarP(&expandRelative, "relative", "r", false, "make path relative to current working directory")
}

// Process expands args and performs all substitution, then returns the argument array
func Process(args []string) []string {
	var processedArgs []string
	for _, arg := range arguments.Expand(args) {
		processed := arguments.EvaluateEnvironment(arg, expandRelative)
		processedArgs = append(processedArgs, processed)
	}

	return processedArgs
}
