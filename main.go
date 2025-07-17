package main

import (
	"fmt"

	"github.com/kmatt/scmpuff/internal/commands/exec"
	"github.com/kmatt/scmpuff/internal/commands/expand"
	"github.com/kmatt/scmpuff/internal/commands/inits"
	"github.com/kmatt/scmpuff/internal/commands/intro"
	"github.com/kmatt/scmpuff/internal/commands/status"

	"github.com/spf13/cobra"

	_ "embed"
)

var version = "0.0.0"

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
	puffCmd.AddCommand(inits.CommandInit())
	puffCmd.AddCommand(exec.CommandExec())
	puffCmd.AddCommand(expand.CommandExpand())
	puffCmd.AddCommand(status.CommandStatus())

	puffCmd.Execute()
}
