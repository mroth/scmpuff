package main

import "github.com/mroth/scmpuff/internal/cmd"

// version is the default version of the program
// ...in almost all cases this should be overriden by the buildscript.
var version = "0.0.0-development"

func main() {
	cmd.Execute(version)
}
