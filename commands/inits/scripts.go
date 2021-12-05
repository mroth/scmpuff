package inits

import (
	_ "embed"
	"fmt"
)

//go:embed data/status_shortcuts.sh
var scriptStatusShortcuts string

//go:embed data/aliases.sh
var scriptAliases string

//go:embed data/git_wrapper.sh
var scriptGitWrapper string

func printScript() {
	if outputScript {
		fmt.Println(scriptStatusShortcuts)
	}

	if includeAliases {
		fmt.Println(scriptAliases)
	}

	if wrapGit {
		fmt.Println(scriptGitWrapper)
	}
}
