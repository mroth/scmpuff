package status

import (
	"reflect"
	"testing"
)

// $ git status --porcelain
// A  TODO.md
//  M script/benchmark
// ?? .travis.yml
// ?? commands/status/process_test.go

var testCasesProcessChange = []struct {
	arg      []byte
	expected *StatusItem
}{
	{
		[]byte("A  TODO.md"),
		&StatusItem{
			x:     'A',
			y:     ' ',
			file:  "TODO.md",
			msg:   "  new file",
			col:   neu,
			group: Staged,
		},
	},
  {
    []byte(" M script/benchmark"),
    &StatusItem{
      x:     ' ',
      y:     'M',
      file:  "script/benchmark",
      msg:   "  modified",
      col:   mod,
      group: Unstaged,
    },
  },
  {
    []byte("?? .travis.yml"),
    &StatusItem{
      x:     '?',
      y:     '?',
      file:  ".travis.yml",
      msg:   " untracked",
      col:   unt,
      group: Untracked,
    },
  },
}

func TestProcessChange(t *testing.T) {
	for _, tc := range testCasesProcessChange {
		actual := processChange(tc.arg)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("processChange('%s'): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}

//
// $ git status --porcelain -b
// ## master...origin/master [ahead 1]
// A  TODO.md
//  M script/benchmark
// ?? .travis.yml
// ?? commands/status/process_test.go
//
// $ gs
// # On branch: master  |  +1  |  [*] => $e*
// #
// ➤ Changes to be committed
// #
// #       new file: [1] TODO.md
// #
// ➤ Changes not staged for commit
// #
// #       modified: [2] script/benchmark
// #
// ➤ Untracked files
// #
// #      untracked: [3] .travis.yml
// #      untracked: [4] commands/status/process_test.go
// #
