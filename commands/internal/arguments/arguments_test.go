package arguments

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Expansion of multiple args at the same time
var testExpandCases = []struct {
	args, expected string
}{
	{"1 3 7", "$e1 $e3 $e7"},
	{"1-3 6", "$e1 $e2 $e3 $e6"},
	{"seven 2-5 1", "seven $e2 $e3 $e4 $e5 $e1"},
	/*
		Testcases for https://github.com/mroth/scmpuff/issues/69, git subcommand
		that can accept a dangling numeric argument. This is annoying because it
		looks like a scmpuff numeric ref, so we'll have to special case those if
		we want to allow this git porcelain CLI parsing behavior.
	*/
	{"git log -1 1", "git log -1 $e1"},                       // was not an issue to begin with
	{"git log -n1 2", "git log -n1 $e2"},                     // was not an issue to begin with
	{"git log -n 1", "git log -n 1"},                         // simple case
	{"git log -n 1 2", "git log -n 1 $e2"},                   // simple case
	{"git log --max-count 1 2", "git log --max-count 1 $e2"}, // other log instances
	{"git log --skip 1 2", "git log --skip 1 $e2"},           // other log instances
	{"git log --min-parents 2 --max-parents 5 1-3", "git log --min-parents 2 --max-parents 5 $e1 $e2 $e3"},
	{"git log -g -n 1 -i 1", "git log -g -n 1 -i $e1"}, // mixed in
	{"git rm -n 1", "git rm -n $e1"},                   // don't just check for -n, its subcommand specific
}

func TestExpand(t *testing.T) {
	t.Setenv("SCMPUFF_GIT_CMD", "git")
	for _, tc := range testExpandCases {
		t.Run(tc.args, func(t *testing.T) {
			// split here to emulate what Cobra will pass us but still write tests with
			// normal looking strings
			args := strings.Split(tc.args, " ")
			expected := strings.Split(tc.expected, " ")
			actual := Expand(args)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("ExpandArgs(%v): expected %v, actual %v", tc.args, expected, actual)
			}
		})
	}
}

// Expansion of a single arg, which might still be a range
var testExpandArgCases = []struct {
	arg      string
	expected []string
}{
	{"1", []string{"$e1"}},                 // single digit
	{"1-3", []string{"$e1", "$e2", "$e3"}}, // range
	{"seven", []string{"seven"}},           // no moleste
}

func TestExpandArg(t *testing.T) {
	for _, tc := range testExpandArgCases {
		actual := expandArg(tc.arg)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("ExpandArg(%v): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}

func TestEvaluateEnvironment(t *testing.T) {
	// It would be wonderful to use fstest.MapFS here and have the function rely
	// upon fs.StatFS, however MapFS currently does not work with absolute paths
	// at all, which makes it useless for our testing here.  Bummer. So we use a
	// testdata fixture instead.
	/* var mockFS = fstest.MapFS{
		"/foo/a.txt":         &fstest.MapFile{Mode: 0644},
		"/foo/b.txt":         &fstest.MapFile{Mode: 0644},
		"/foo/bar/c.txt":     &fstest.MapFile{Mode: 0644},
		"/usr/local/bin/git": &fstest.MapFile{Mode: 0777},
	} */
	wd, err := os.Getwd()
	if err != nil {
		t.Skip("failed to get wd, cannot test")
	}

	fakegitAbsPath := filepath.Join(wd, "testdata", "bin", "fakegit")
	t.Setenv("SCMPUFF_GIT_CMD", fakegitAbsPath)
	t.Setenv("e1", filepath.Join(wd, "testdata", "a.txt"))
	t.Setenv("e2", filepath.Join(wd, "testdata", "b.txt"))
	t.Setenv("FOO_USER", "not_a_file")

	t.Logf("$CWD=%v", wd)
	t.Logf("$SCMPUFF_GIT_CMD=%v", os.Getenv("SCMPUFF_GIT_CMD"))
	t.Logf("$e1=%v", os.Getenv("e1"))
	t.Logf("$e2=%v", os.Getenv("e2"))

	tests := []struct {
		name           string
		arg            string
		expandRelative bool
		want           string
	}{
		{name: "not an env var", arg: "eee", expandRelative: false, want: "eee"},
		{name: "not file absolute", arg: "$FOO_USER", expandRelative: false, want: "not_a_file"},
		{name: "not file relative", arg: "$FOO_USER", expandRelative: true, want: "not_a_file"},
		{name: "absolute file", arg: "$e1", expandRelative: false, want: filepath.Join(wd, "testdata", "a.txt")},
		{name: "relative file", arg: "$e1", expandRelative: true, want: filepath.FromSlash("testdata/a.txt")},
		{name: "path binary dont convert relative - abs", arg: "$SCMPUFF_GIT_CMD", expandRelative: false, want: fakegitAbsPath},
		{name: "path binary dont convert relative - rel", arg: "$SCMPUFF_GIT_CMD", expandRelative: true, want: fakegitAbsPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvaluateEnvironment(tt.arg, tt.expandRelative); got != tt.want {
				t.Errorf("EvaluateEnvironment(%v, %v) = %v, want %v", tt.arg, tt.expandRelative, got, tt.want)
			}
		})
	}
}
