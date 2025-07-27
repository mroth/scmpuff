package cmd

import (
	"fmt"
	"os"

	goversion "github.com/caarlos0/go-version"
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
		Version: version,
		Args:    cobra.NoArgs,

		// disable default completions introduced in cobra v1.2.0, we will want to
		// customize if we provide these in the future.
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.SetVersionTemplate("{{.Version}}")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number",
		Long:  `All software has versions. This is ours.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
		Deprecated: "use `scmpuff --version` instead",
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
		Args:   cobra.NoArgs,
		Hidden: true, // hide from help output
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // silence usage-on-error after args processed
			return fmt.Errorf("this is a test error")
		},
	}
	rootCmd.AddCommand(errCmd)

	return rootCmd
}

// Execute executes the root command.
func Execute(version goversion.Info) {
	cmd := newRootCmd(version.String())
	if err := cmd.Execute(); err != nil {
		// Cobra already prints the error, so we just need to exit.
		//
		// Currently there are a few places we exit directly in the rather than
		// returning control to here, typically to enforce a specific error code
		// of UX convention.
		//
		// Set cmd.SilenceUsage on erroring commands to avoid printing usage if
		// desired. Note that this can be done in the command's RunE function
		// after args are processed, in order to print usage only on arg parsing
		// errors.
		os.Exit(1)
	}
}
