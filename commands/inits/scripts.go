package inits

import (
	_ "embed"
	"fmt"
)

//go:embed data/status_shortcuts.sh
var scriptStatusShortcuts string

//go:embed data/status_shortcuts.fish
var scriptStatusShortcutsFish string

//go:embed data/aliases.sh
var scriptAliases string

//go:embed data/git_wrapper.sh
var scriptGitWrapper string

//go:embed data/git_wrapper.fish
var scriptGitWrapperFish string

func printScript() {
	if outputScript {
		if shellType == "fish" {
			fmt.Println(scriptStatusShortcutsFish)
		} else {
			fmt.Println(scriptStatusShortcuts)
		}
	}

	if includeAliases {
		fmt.Println(scriptAliases)
	}

	if wrapGit {
		if shellType == "fish" {
			fmt.Println(scriptGitWrapperFish)
		} else {
			fmt.Println(scriptGitWrapper)
		}
	}
}
