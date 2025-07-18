package main

import "github.com/mroth/scmpuff/internal/cmd"

// version is the default version of the program
// ...in almost all cases this should be overriden by the buildscript.
var version = "0.0.0-development"

func main() {
	// TODO: instead of passing version via main, change where we embed with build script
	cmd.Version = version
	cmd.Execute()
}
