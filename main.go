package main

import "github.com/spf13/cobra"

var puffCmd = &cobra.Command{
	Use:   "scmpuff",
	Short: "scmpuff extends common git commands with numeric filename shortcuts.",
	Long: `scmpuff extends common git commands with numeric filename shortcuts.
  Built by @mroth, with huge props to @nbroadbent for the original.
  `,
}

func main() {

	puffCmd.AddCommand(versionCmd)
	puffCmd.AddCommand(CommandInit())
	puffCmd.AddCommand(CommandExpand())

	puffCmd.Execute()
}
