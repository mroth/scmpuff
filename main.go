package main

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

// version is the default version of the program
// ...in almost all cases this should be overriden by the buildscript.
var version = "0.0.0-development"

var puffCmd = &cobra.Command{
	Use:   "scmpuff",
	Short: "scmpuff extends common git commands with numeric filename shortcuts.",
	Long: `scmpuff extends common git commands with numeric filename shortcuts.

If you are just getting started, try the intro!`,

	// disable default completions introduced in cobra v1.2.0, we will want to
	// customize if we provide these in the future.
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long:  `All software has versions. This is ours.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scmpuff", version)
	},
}

func main() {
	puffCmd.AddCommand(versionCmd)
	puffCmd.AddCommand(intro.IntroCmd)
	puffCmd.AddCommand(debug.DebugCmd)
	puffCmd.AddCommand(exec.ExecCmd)
	puffCmd.AddCommand(expand.ExpandCmd)
	puffCmd.AddCommand(inits.InitCmd)
	puffCmd.AddCommand(status.StatusCmd)

	puffCmd.Execute()
}
