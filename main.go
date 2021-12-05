package main

import (
	"fmt"

	"github.com/mroth/scmpuff/commands/exec"
	"github.com/mroth/scmpuff/commands/expand"
	"github.com/mroth/scmpuff/commands/inits"
	"github.com/mroth/scmpuff/commands/status"

	"github.com/spf13/cobra"
)

// NAME of the program, hardcoded for consistency
var NAME = "scmpuff"

// VERSION is the default version of the program
// ...in almost all cases this should be overriden by the buildscript.
var VERSION = "0.?.? (not using buildscript)"

var puffCmd = &cobra.Command{
	Use:   "scmpuff",
	Short: "scmpuff extends common git commands with numeric filename shortcuts.",
	Long: `scmpuff extends common git commands with numeric filename shortcuts.

If you are just getting started, try the intro!`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long:  `All software has versions. This is ours.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(NAME, VERSION)
	},
}

func main() {
	puffCmd.AddCommand(introCmd)
	puffCmd.AddCommand(versionCmd)
	puffCmd.AddCommand(inits.CommandInit())
	puffCmd.AddCommand(exec.CommandExec())
	puffCmd.AddCommand(expand.CommandExpand())
	puffCmd.AddCommand(status.CommandStatus())

	puffCmd.Execute()
}
