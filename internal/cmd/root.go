package cmd

import (
	"fmt"
	"os"

	"github.com/mroth/scmpuff/internal/cmd/debug"
	"github.com/mroth/scmpuff/internal/cmd/exec"
	"github.com/mroth/scmpuff/internal/cmd/expand"
	"github.com/mroth/scmpuff/internal/cmd/inits"
	"github.com/mroth/scmpuff/internal/cmd/intro"
	"github.com/mroth/scmpuff/internal/cmd/status"
	"github.com/spf13/cobra"
)

// newRootCmd creates and returns the root command
func newRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "scmpuff",
		Short: "scmpuff extends common git commands with numeric filename shortcuts.",
		Long: `scmpuff extends common git commands with numeric filename shortcuts.

If you are just getting started, try the intro!`,

		// disable default completions introduced in cobra v1.2.0, we will want to
		// customize if we provide these in the future.
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

		// don't print usage on error, just the error message
		// many commands will print usage themselves if needed
		SilenceUsage: true,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number",
		Long:  `All software has versions. This is ours.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("scmpuff", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(intro.NewIntroCmd())
	rootCmd.AddCommand(debug.NewDebugCmd())
	rootCmd.AddCommand(exec.NewExecCmd())
	rootCmd.AddCommand(expand.NewExpandCmd())
	rootCmd.AddCommand(inits.NewInitCmd())
	rootCmd.AddCommand(status.NewStatusCmd())

	// For now, add a command that always returns an error for testing purposes
	errCmd := &cobra.Command{
		Use:    "xerror",
		Short:  "A command that always returns an error",
		Long:   `This command is for testing error handling.`,
		Hidden: true, // hide from help output
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this is a test error")
		},
	}
	rootCmd.AddCommand(errCmd)

	return rootCmd
}

// Execute executes the root command.
func Execute(version string) {
	cmd := newRootCmd(version)
	if err := cmd.Execute(); err != nil {
		// not all commands will return, currently many exit directly on their own
		// Cobra already prints the error, so we just need to exit
		os.Exit(1)
	}
}
