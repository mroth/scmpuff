package cmd

import (
	"flag"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var updateScripts = flag.Bool("update", false, "update testscript cmp fixtures")

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"scmpuff": runScmpuff,
	})
}

func runScmpuff() {
	root := newRootCmd("test")
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:                 "testdata/script",
		RequireExplicitExec: true,
		UpdateScripts:       *updateScripts,
		Setup: func(e *testscript.Env) error {
			e.Setenv("GIT_AUTHOR_NAME", "SCM Puff")
			e.Setenv("GIT_AUTHOR_EMAIL", "scmpuff@example.com")
			e.Setenv("GIT_COMMITTER_NAME", "SCM Puff")
			e.Setenv("GIT_COMMITTER_EMAIL", "scmpuff@example.com")
			return nil
		},
	})
}
