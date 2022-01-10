package inits

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Since the flags are defined and used in different locations, we need to
// define a variable outside with the correct scope to assign the flag to work
// with.
var includeAliases bool
var legacyShow bool
var wrapGit bool
var shellType string

// CommandInit generates the command handler for `scmpuff init`
func CommandInit() *cobra.Command {

	var InitCmd = &cobra.Command{
		Use:   "init",
		Short: "Output initialization script",
		Long: `
Outputs the bash/zsh/fish initialization script for scmpuff.

This should probably be evaluated in your shell startup.
    `,
		Run: func(cmd *cobra.Command, args []string) {
			// If someone's using the old -s/--show flag, opt-in to the newer --shell=sh option
			if legacyShow  {
				shellType = "sh"
			}
			if shellType != "" {
				printScript()
			} else {
				fmt.Println(helpString())
			}
		},
		// Watch out for accidental args caused by NoOptDefVal (https://github.com/spf13/cobra/issues/866)
		Args: cobra.NoArgs,
	}

	// --aliases
	InitCmd.Flags().BoolVarP(
		&includeAliases,
		"aliases", "a", true,
		"Include short aliases for convenience",
	)

	// --show (deprecated in favor of --shell)
	InitCmd.Flags().BoolVar(
		&legacyShow,
		"show", false,
		"Output scmpuff initialization scripts",
	)
	InitCmd.Flags().MarkHidden("show")

	// --wrap
	InitCmd.Flags().BoolVarP(
		&wrapGit,
		"wrap", "w", true,
		"Wrap standard git commands",
	)

	// --shell
	InitCmd.Flags().StringVarP(
		&shellType,
		"shell", "s", "",
		"Set shell type - 'sh' (for bash/zsh), or 'fish'",
	)
	InitCmd.Flag("shell").NoOptDefVal = "sh"

	return InitCmd
}

// TODO: check for proper shell version
func helpString() string {
	return `# Initialize scmpuff by adding the following to ~/.bash_profile or ~/.zshrc:

eval "$(scmpuff init --shell=sh)"

# or the following to ~/.config/fish/config.fish:

scmpuff init --shell=fish | source`
}
