package status

import (
	"reflect"
	"testing"
)

// single test to make sure everything gets stiched together properly, test
// actual cases in more specific methods
func TestProcessChange(t *testing.T) {
	chunk := []byte("A  HELLO.md")
	actual := processChange(chunk, "/tmp", "/tmp")[0]

	expectedChange := &change{
		msg:   "  new file",
		col:   neu,
		group: Staged,
	}
	if actual.col != expectedChange.col ||
		actual.group != expectedChange.group ||
		actual.msg != expectedChange.msg {
		t.Fatal("changes did not match expected")
	}

	if actual.fileAbsPath != "/tmp/HELLO.md" {
		t.Fatal("absolute path did not match expected")
	}

	if actual.fileRelPath == "" {
		t.Fatal("relative path was not present")
	}
}

var testCasesExtractFile = []struct {
	root        string
	wd          string
	chunk       []byte
	expectedAbs string
	expectedRel string
}{
	{
		root:        "/",
		wd:          "/",
		chunk:       []byte(" M script/benchmark"),
		expectedAbs: "/script/benchmark",
		expectedRel: "script/benchmark",
	},
	{
		root:        "/tmp",
		wd:          "/tmp",
		chunk:       []byte(" M script/benchmark"),
		expectedAbs: "/tmp/script/benchmark",
		expectedRel: "script/benchmark",
	},
	{
		root:        "/tmp/foo/bar//",
		wd:          "/tmp/foo/bar/unicorn",
		chunk:       []byte("?? unicorn/magic/xxx"),
		expectedAbs: "/tmp/foo/bar/unicorn/magic/xxx",
		expectedRel: "magic/xxx",
	},
	{
		root:        "/tmp/foo/bar//",
		wd:          "/tmp/foo/bar/unicorn/magic",
		chunk:       []byte("?? narwhal/disco/yyy"),
		expectedAbs: "/tmp/foo/bar/narwhal/disco/yyy",
		expectedRel: "../../narwhal/disco/yyy",
	},
	{
		root:        "/tmp/foo",
		wd:          "/tmp/foo",
		chunk:       []byte("R  bar.txt\x00foo.txt"),
		expectedAbs: "/tmp/foo/bar.txt",
		expectedRel: "foo.txt -> bar.txt",
	},
	// following examples are ones where scm_breeze strips the escaping that
	// git status --porcelain does in certain cases.  Now that we are using -z
	// we dont have escaped characters in our output (or our tests), so this is
	// fairly redundant as a unit test...
	// (historical note of why we did this: scm_breeze uses the escaped versions
	//  of output and can fail on complex cases of parsing it!)
	{
		root:        "/tmp/foo",
		wd:          "/tmp/foo",
		chunk:       []byte(`A  hi there mom.txt`),
		expectedAbs: "/tmp/foo/hi there mom.txt",
		expectedRel: "hi there mom.txt",
	},
	{
		root:        "/tmp/foo",
		wd:          "/tmp/foo/bar",
		chunk:       []byte(`?? "x.txt`),
		expectedAbs: `/tmp/foo/"x.txt`,
		expectedRel: `../"x.txt`,
	},
	{
		root:        "/tmp/foo",
		wd:          "/tmp/foo",
		chunk:       []byte(`?? hi m"o"m.txt`),
		expectedAbs: `/tmp/foo/hi m"o"m.txt`, //scmbreeze fails these with `hi m"o\`
		expectedRel: `hi m"o"m.txt`,
	},
}

func TestExtractFile(t *testing.T) {
	for _, tc := range testCasesExtractFile {
		actualAbs, actualRel := extractFile(tc.chunk, tc.root, tc.wd)

		if actualAbs != tc.expectedAbs {
			t.Fatalf(
				"extractFile(%s)/absPath:\nexpect\t%v\nactual\t%v",
				tc.chunk, tc.expectedAbs, actualAbs)
		}
		if actualRel != tc.expectedRel {
			t.Fatalf(
				"extractFile(%s)/relPath:\nexpect\t%v\nactual\t%v",
				tc.chunk, tc.expectedRel, actualRel)
		}
	}
}

// $ git status --porcelain
// A  HELLO.md
//  M script/benchmark
// ?? .travis.yml
// ?? commands/status/process_test.go
var testCasesExtractChangeCodes = []struct {
	chunk    []byte
	expected []*change
}{
	{
		[]byte("A  HELLO.md"),
		[]*change{
			&change{msg: "  new file", col: neu, group: Staged},
		},
	},
	{
		[]byte(" M script/benchmark"),
		[]*change{
			&change{msg: "  modified", col: mod, group: Unstaged},
		},
	},
	{
		[]byte("?? .travis.yml"),
		[]*change{
			&change{msg: " untracked", col: unt, group: Untracked},
		},
	},
	{
		[]byte(" D deleted_file"),
		[]*change{
			&change{msg: "   deleted", col: del, group: Unstaged},
		},
	},
	{
		[]byte("R  after\x00before"),
		[]*change{
			&change{msg: "   renamed", col: ren, group: Staged},
		},
	},
	{
		[]byte("C  after\x00before"),
		[]*change{
			&change{msg: "    copied", col: cpy, group: Staged},
		},
	},
	{
		[]byte("AM added_then_modified_file"),
		[]*change{
			&change{msg: "  new file", col: neu, group: Staged},
			&change{msg: "  modified", col: mod, group: Unstaged},
		},
	},
}

func TestExtractChangeCodes(t *testing.T) {
	for _, tc := range testCasesExtractChangeCodes {
		actual := extractChangeCodes(tc.chunk)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("processChange('%s'): expected %+v, actual %+v",
				tc.chunk, tc.expected, actual)
		}
	}
}

// Examples of stuff we will want to parse:
//
// 		## Initial commit on master
// 		## master
// 		## master...origin/master
// 		## master...origin/master [ahead 1]
var testCasesExtractBranch = []struct {
	chunk    []byte
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
	{
		[]byte("## upstream...upstream/master [ahead 5, behind 3]"),
		&BranchInfo{name: "upstream", ahead: 5, behind: 3},
	},
	{
		[]byte("## 3.0...origin/3.0 [ahead 1]"),
		&BranchInfo{name: "3.0", ahead: 1, behind: 0},
	},
}

func TestExtractBranch(t *testing.T) {
	for _, tc := range testCasesExtractBranch {
		actual := ExtractBranch(tc.chunk)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Fatalf("processBranch('%s'): expected %v, actual %v",
				tc.chunk, tc.expected, actual)
		}
	}
}
