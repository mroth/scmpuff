package inits

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Since the flags are defined and used in different locations, we need to
// define a variable outside with the correct scope to assign the flag to work
// with.
var includeAliases bool
var outputScript bool

// CommandInit generates the command handler for `scmpuff init`
func CommandInit() *cobra.Command {

	var InitCmd = &cobra.Command{
		Use:   "init",
		Short: "Output initialization script",
		Long: `
Output the bash/zsh initialization script for scmpuff.

This should probably be evaluated in your shell startup.
    `,
		Run: func(cmd *cobra.Command, args []string) {
			if outputScript {
				printScript(includeAliases)
			} else {
				fmt.Println(helpString())
			}
		},
	}

	// --aliases
	InitCmd.Flags().BoolVarP(
		&includeAliases,
		"aliases", "a", true,
		"Define short aliases for convenience",
	)

	// --show
	InitCmd.Flags().BoolVarP(
		&outputScript,
		"show", "s", false,
		"Show outputscript",
	)

	return InitCmd
}

func printScript(includeAliases bool) {
	fmt.Println(scriptString())
	fmt.Println(statusShortcutsString())
	if includeAliases {
		fmt.Println(aliasesString())
	}
}

// TODO: check for proper shell version
func helpString() string {
	return `# Wrap git automatically by adding the following to ~/.zshrc:

eval "$(scmpuff init -s)"`
}

func scriptString() string {
	return `git () {
  case $1 in
    (commit|blame|add|log|rebase|merge) scmpuff expand "$_git_cmd" "$@" ;;
    (checkout|diff|rm|reset) scmpuff expand --relative "$_git_cmd" "$@" ;;
    (branch) _scmb_git_branch_shortcuts "${@:2}" ;;
    (*) "$_git_cmd" "$@" ;;
  esac
}`
}

func aliasesString() string {
	return `
alias gs='scmpuff status'
alias ga='scmpuff add'
  `
}
