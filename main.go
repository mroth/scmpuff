package main

import (
	"github.com/mroth/scmpuff/commands/expand"
	"github.com/mroth/scmpuff/commands/inits"
	"github.com/mroth/scmpuff/commands/status"
	"github.com/mroth/scmpuff/commands/version"
	"github.com/spf13/cobra"
)

var puffCmd = &cobra.Command{
	Use:   "scmpuff",
	Short: "scmpuff extends common git commands with numeric filename shortcuts.",
	Long: `scmpuff extends common git commands with numeric filename shortcuts.
  Built by @mroth, with huge props to @nbroadbent for the original.
  `,
}

func main() {

	puffCmd.AddCommand(version.VersionCmd)
	puffCmd.AddCommand(inits.CommandInit())
	puffCmd.AddCommand(expand.CommandExpand())
	puffCmd.AddCommand(status.CommandStatus())

	puffCmd.Execute()
}
