package inits

import (
	_ "embed"
	"strings"
)

//go:embed data/status_shortcuts.sh
var scriptStatusShortcuts string

//go:embed data/status_shortcuts.fish
var scriptStatusShortcutsFish string

////go:embed data/status_shortcuts.pwsh
//var scriptStatusShortcutsPwsh string

//go:embed data/aliases.sh
var scriptAliases string

//go:embed data/git_wrapper.sh
var scriptGitWrapper string

//go:embed data/git_wrapper.fish
var scriptGitWrapperFish string

//go:embed data/git_wrapper.ps1
var scriptGitWrapperPwsh string

type scriptCollection struct {
	statusShortcuts string
	gitWrapper      string
	aliases         string
}

var bashCollection = scriptCollection{
	statusShortcuts: scriptStatusShortcuts,
	gitWrapper:      scriptGitWrapper,
	aliases:         scriptAliases,
}

var fishCollection = scriptCollection{
	statusShortcuts: scriptStatusShortcutsFish,
	gitWrapper:      scriptGitWrapperFish,
	aliases:         scriptAliases,
}

var pwshCollection = scriptCollection{
	//statusShortcuts: scriptStatusShortcutsPwsh,
	gitWrapper:      scriptGitWrapperPwsh,
	aliases:         scriptAliases,
}

func (sc scriptCollection) Output(wrapGit, aliases bool) string {
	var b strings.Builder
	b.WriteString(sc.statusShortcuts)
	if wrapGit {
		b.WriteRune('\n')
		b.WriteString(sc.gitWrapper)
	}
	if aliases {
		b.WriteRune('\n')
		b.WriteString(sc.aliases)
	}
	return b.String()
}
