package exec

import (
	"os"
	"os/exec"

	"github.com/mroth/scmpuff/internal/arguments"
	"github.com/spf13/cobra"
)

var expandRelative bool

// NewExecCmd creates and returns the exec command
func NewExecCmd() *cobra.Command {
	execCmd := &cobra.Command{
		Use:     "exec [flags] <command> <shortcuts...>",
		Example: "$ scmpuff exec -- git add 1-4",
		Aliases: []string{"execute"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Execute cmd with numeric shortcuts",
		Long: `Expands numeric shortcuts to their full filepath and executes the command.

Takes a list of digits (1 4 5) or numeric ranges (1-5) or even both.`,
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			cmd.SilenceUsage = true // silence usage-on-error after args processed

			expandedArgs := Process(inputArgs)
			a := expandedArgs[1:]
			subcmd := exec.Command(expandedArgs[0], a...)
			subcmd.Stdin = os.Stdin
			subcmd.Stdout = os.Stdout
			subcmd.Stderr = os.Stderr

			err := subcmd.Run()
			if err != nil {
				// process exited with a non-zero exit code, we want to just exit
				// directly with that code rather than returning control back to cobra.
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}

				// otherwise, we failed to start execution, return error to cobra
				return err
			}

			// normal case: exec completed successfully.
			return nil
		},
	}

	execCmd.Flags().BoolVarP(&expandRelative, "relative", "r", false, "make path relative to current working directory")
	return execCmd
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
