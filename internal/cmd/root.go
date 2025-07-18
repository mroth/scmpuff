package cmd

import (
	"fmt"

	"github.com/mroth/scmpuff/internal/cmd/debug"
	"github.com/mroth/scmpuff/internal/cmd/exec"
	"github.com/mroth/scmpuff/internal/cmd/expand"
	"github.com/mroth/scmpuff/internal/cmd/inits"
	"github.com/mroth/scmpuff/internal/cmd/intro"
	"github.com/mroth/scmpuff/internal/cmd/status"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scmpuff",
	Short: "scmpuff extends common git commands with numeric filename shortcuts.",
	Long: `scmpuff extends common git commands with numeric filename shortcuts.

If you are just getting started, try the intro!`,

	// disable default completions introduced in cobra v1.2.0, we will want to
	// customize if we provide these in the future.
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// Very hacky way to pass the version for now
var Version string = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long:  `All software has versions. This is ours.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scmpuff", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(intro.IntroCmd)
	rootCmd.AddCommand(debug.DebugCmd)
	rootCmd.AddCommand(exec.ExecCmd)
	rootCmd.AddCommand(expand.ExpandCmd)
	rootCmd.AddCommand(inits.InitCmd)
	rootCmd.AddCommand(status.StatusCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
