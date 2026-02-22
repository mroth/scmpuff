package inits

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewInitCmd creates and returns the init command
func NewInitCmd() *cobra.Command {
	var (
		shellType      string
		includeAliases bool
		wrapGit        bool
		legacyShow     bool
	)

	initCmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// If someone's using the old ---show flag, opt-in to the newer --shell defaults
			if legacyShow {
				shellType = defaultShellType()
			}

			switch strings.ToLower(shellType) {
			case "":
				cmd.Help()
			case "sh", "bash", "zsh":
				fmt.Fprintln(cmd.OutOrStdout(), bashCollection.Output(wrapGit, includeAliases))
			case "fish":
				fmt.Fprintln(cmd.OutOrStdout(), fishCollection.Output(wrapGit, includeAliases))
			default:
				return fmt.Errorf(`unrecognized shell "%s"`, shellType)
			}

			return nil
		},
		// Watch out for accidental args caused by NoOptDefVal (https://github.com/spf13/cobra/issues/866)
		Args: cobra.NoArgs,
	}

	// --aliases
	initCmd.Flags().BoolVarP(
		&includeAliases,
		"aliases", "a", true,
		"Include short git aliases",
	)

	// --show (deprecated in favor of --shell)
	initCmd.Flags().BoolVar(
		&legacyShow,
		"show", false,
		"Output scmpuff initialization scripts",
	)
	initCmd.Flags().MarkDeprecated("show", "use --shell instead")

	// --wrap
	initCmd.Flags().BoolVarP(
		&wrapGit,
		"wrap", "w", true,
		"Wrap standard git commands",
	)

	// --shell
	initCmd.Flags().StringVarP(
		&shellType,
		"shell", "s", "",
		"Output shell type: sh | bash | zsh | fish",
	)
	initCmd.Flag("shell").NoOptDefVal = defaultShellType()

	return initCmd
}

// defaultShell returns the shellType assumed if user does not specify.
func defaultShellType() string {
	if shellenv, ok := os.LookupEnv("SHELL"); ok {
		base := filepath.Base(shellenv)
		switch base {
		case "sh", "bash", "zsh", "fish":
			return base
		}
	}

	return "sh"
}
