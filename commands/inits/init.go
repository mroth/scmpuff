package inits

import (
	"fmt"
	"os"

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
Outputs the shell initialization script for scmpuff.

Initialize scmpuff by adding the following to your ~/.bash_profile or ~/.zshrc:

    eval "$(scmpuff init --shell=sh)"

For fish shell, add the following to ~/.config/fish/config.fish instead:

    scmpuff init --shell=fish | source

There are a number of flags to customize the shell integration.
    `,
		Run: func(cmd *cobra.Command, args []string) {
			// If someone's using the old ---show flag, opt-in to the newer --shell defaults
			if legacyShow {
				shellType = defaultShellType()
			}
			switch shellType {
			case "":
				cmd.Help()
				os.Exit(0)

			case "sh", "bash", "zsh", "fish":
				printScript()
				os.Exit(0)

			default:
				fmt.Fprintf(os.Stderr, "Unrecognized shell '%s'\n", shellType)
				os.Exit(1)
			}
		},
		// Watch out for accidental args caused by NoOptDefVal (https://github.com/spf13/cobra/issues/866)
		Args: cobra.NoArgs,
	}

	// --aliases
	InitCmd.Flags().BoolVarP(
		&includeAliases,
		"aliases", "a", true,
		"Include short git aliases",
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
		"Output shell type: sh | bash | zsh | fish",
	)
	InitCmd.Flag("shell").NoOptDefVal = defaultShellType()

	return InitCmd
}

// defaultShell returns the shellType assumed if user does not specify.
// in the future, we may wish to customize this based on the $SHELL variable.
func defaultShellType() string {
	return "sh"
}
