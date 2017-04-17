package main

import (
	"fmt"

	"github.com/mroth/scmpuff/commands/expand"
	"github.com/mroth/scmpuff/commands/inits"
	"github.com/mroth/scmpuff/commands/status"

	"github.com/spf13/cobra"
)

var NAME = "scmpuff"
var VERSION = "0.?.? (not using buildscript)" // overridden via build scripts

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
	puffCmd.AddCommand(expand.CommandExpand())
	puffCmd.AddCommand(status.CommandStatus())

	puffCmd.Execute()
}
