package inits

import "fmt"

// the below go:generate directive will automatically generate a bindata.go file
// which wraps the contents of the data directory so we can include text files
// in our distributed binary directly. (boy, Go can be quite annoying
// sometimes!)

//go:generate go-bindata -o bindata.go -ignore=README* -pkg=inits data

func printScript() {
	fmt.Println(assetString("data/status_shortcuts.sh"))

	if includeAliases {
		fmt.Println(aliasesString())
	}
	// if wrapGit {
	// fmt.Println(gitwrapString()) 	// TODO: enable wrapping git cmds once expand works && add tests
	// }
}

// TODO: check for proper shell version
func helpString() string {
	return `# Initialize scmpuff by adding the following to ~/.bash_profile or ~/.zshrc:

eval "$(scmpuff init -s)"`
}

func gitwrapString() string {
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
alias gs='scmpuff_status_shortcuts'
#alias ga='scmpuff add' #TODO: implement me
  `
}

// returns the string data for an embedded data script
func assetString(file string) string {
	data, err := Asset(file)
	if err != nil {
		// Asset was not found.
	}
	return string(data)
}
