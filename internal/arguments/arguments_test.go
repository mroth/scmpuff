package arguments

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
}

func TestExpand(t *testing.T) {
	for _, tc := range testExpandCases {
		// split here to emulate what Cobra will pass us but still write tests with
		// normal looking strings
		args := strings.Split(tc.args, " ")
		expected := strings.Split(tc.expected, " ")
		actual := Expand(args)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ExpandArgs(%v): expected %v, actual %v", tc.args, expected, actual)
		}
	}
}

// Test that numeric arguments to known git flags are not incorrectly expanded
// to file shortcuts. For example, "git log -n 1" should keep "1" as a literal
// count, not turn it into "$e1". This only kicks in when args[0] matches
// SCMPUFF_GIT_CMD, mirroring how the shell wrapper invokes "scmpuff exec".
//
// Each test case is written as a space-separated command line. The left side is
// the input args, the right side is the expected output after expansion. A bare
// number in the output (e.g. "1") means it was protected from expansion; "$eN"
// means it was expanded as a file shortcut.
var testExpandNumericFlagCases = []struct {
	args, expected string
}{
	// Flags that take a numeric count separated by a space — the count must
	// not be treated as a file shortcut.
	{"git log -n 1", "git log -n 1"},
	{"git log -n 1 2", "git log -n 1 $e2"},
	{"git log --max-count 1 2", "git log --max-count 1 $e2"},
	{"git log --skip 1 2", "git log --skip 1 $e2"},
	{"git log --grep 1 2", "git log --grep 1 $e2"},
	{"git blame -L 1 1", "git blame -L 1 $e1"},
	{"git rebase -C 3 1", "git rebase -C 3 $e1"},
	{"git rebase --onto 713 main topic", "git rebase --onto 713 main topic"},

	// Flags that take a string value that could happen to be all digits.
	{"git commit -m 123", "git commit -m 123"},
	{"git commit --message 456", "git commit --message 456"},
	{"git merge -m 123 topic", "git merge -m 123 topic"},

	// Flags where the value is glued to the flag (no space) — the combined
	// token won't match the digit regex, so expansion is already harmless.
	{"git log -n1 2", "git log -n1 $e2"},
	{"git log -1 1", "git log -1 $e1"},

	// Flags that take a non-numeric argument that could happen to be all
	// digits — e.g. branch names like "713".
	{"git checkout -b 713", "git checkout -b 713"},
	{"git checkout -B 42", "git checkout -B 42"},
	{"git checkout --orphan 99", "git checkout --orphan 99"},

	// Same flag name, different subcommand: -n on "git rm" is --dry-run and
	// takes no value, so the following "1" is a file shortcut.
	{"git rm -n 1", "git rm -n $e1"},
}

func TestExpandNumericFlags(t *testing.T) {
	t.Setenv("SCMPUFF_GIT_CMD", "git")
	for _, tc := range testExpandNumericFlagCases {
		t.Run(tc.args, func(t *testing.T) {
			args := strings.Split(tc.args, " ")
			expected := strings.Split(tc.expected, " ")
			actual := Expand(args)
			if !slices.Equal(actual, expected) {
				t.Errorf("expected %v, actual %v", expected, actual)
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
			t.Fatalf("ExpandArg(%v): expected %v, actual %v", tc.arg, tc.expected, actual)
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
