package main

import (
	goversion "github.com/caarlos0/go-version"
	"github.com/mroth/scmpuff/internal/cmd"
)

var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

func main() {
	cmd.Execute(
		buildVersion(version, commit, date, builtBy, treeState),
	)
}

var asciiArt = `                                    ________
   ______________ ___  ____  __  __/ __/ __/
  / ___/ ___/ __ ` + "`" + `__ \/ __ \/ / / / /_/ /_
 (__  ) /__/ / / / / / /_/ / /_/ / __/ __/
/____/\___/_/ /_/ /_/ .___/\__,_/_/ /_/
                   /_/
`

func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("scmpuff", "Git by the numbers.", "https://mroth.github.io/scmpuff/"),
		goversion.WithASCIIName(asciiArt),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}
		},
	)
}
