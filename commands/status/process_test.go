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
	{
		[]byte(" D deleted_file"),
		&StatusItem{
			x:     ' ',
			y:     'D',
			file:  "deleted_file",
			msg:   "   deleted",
			col:   del,
			group: Unstaged,
		},
	},
}

func TestProcessChange(t *testing.T) {
	for _, tc := range testCasesProcessChange {
		actual := ProcessChange(tc.arg)
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

// Examples of stuff we will want to parse:
//
// 		## Initial commit on master
// 		## master
// 		## master...origin/master
// 		## master...origin/master [ahead 1]
var testCasesProcessBranch = []struct {
	arg      []byte
	expected *BranchInfo
}{
	{
		[]byte("## master"),
		&BranchInfo{name: "master", ahead: 0, behind: 0},
	},
	{
		[]byte("## GetUpGetDown09-11JokeInYoTown"),
		&BranchInfo{name: "GetUpGetDown09-11JokeInYoTown", ahead: 0, behind: 0},
	},
	{
		[]byte("## master...origin/master"),
		&BranchInfo{name: "master", ahead: 0, behind: 0},
	},
	{
		[]byte("## upstream...upstream/master"),
		&BranchInfo{name: "upstream", ahead: 0, behind: 0},
	},
	{
		[]byte("## master...origin/master [ahead 1]"),
		&BranchInfo{name: "master", ahead: 1, behind: 0},
	},
	{
		[]byte("## upstream...upstream/master [behind 3]"),
		&BranchInfo{name: "upstream", ahead: 0, behind: 3},
	},
	// TODO: test and handle compound up/down status
}

func TestProcessBranch(t *testing.T) {
	for _, tc := range testCasesProcessBranch {
		actual := ProcessBranch(tc.arg)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("processBranch('%s'): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}

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
