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
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func TestScripts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	testscript.Run(t, testscript.Params{
		Dir:                 "testdata/script",
		RequireExplicitExec: true,
		UpdateScripts:       *updateScripts,
		Setup: func(e *testscript.Env) error {
			// Isolate from host git configuration.
			e.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			e.Setenv("GIT_TERMINAL_PROMPT", "0")

			e.Setenv("GIT_AUTHOR_NAME", "SCM Puff")
			e.Setenv("GIT_AUTHOR_EMAIL", "scmpuff@example.com")
			e.Setenv("GIT_COMMITTER_NAME", "SCM Puff")
			e.Setenv("GIT_COMMITTER_EMAIL", "scmpuff@example.com")
			return nil
		},
	})
}
