package porcelainv1

import (
	"reflect"
	"slices"
	"testing"

	"github.com/mroth/scmpuff/internal/commands/status/gitstatus"
)

func Test_extractFilePaths(t *testing.T) {
	var testCases = []struct {
		chunk        []byte
		wantPath     string
		wantOrigPath string
	}{
		{
			chunk:    []byte(" M script/benchmark.sh"),
			wantPath: "script/benchmark.sh",
		},
		{
			chunk:    []byte("?? unicorn/magic/xxx"),
			wantPath: "unicorn/magic/xxx",
		},
		{
			chunk:    []byte("?? file with spaces.txt"),
			wantPath: "file with spaces.txt",
		},
		{
			chunk:        []byte("R  b.txt\x00a.txt"),
			wantPath:     "b.txt",
			wantOrigPath: "a.txt",
		},
		// following examples are ones where scm_breeze strips the escaping that
		// git status --porcelain does in certain cases.  Now that we are using -z
		// we dont have escaped characters in our output (or our tests), so this is
		// fairly redundant as a unit test...
		// (historical note of why we did this: scm_breeze uses the escaped versions
		//  of output and can fail on complex cases of parsing it!)
		{
			chunk:    []byte(`?? "x.txt`),
			wantPath: `"x.txt`,
		},
		{
			chunk:    []byte(`?? hi m"o"m.txt`),
			wantPath: `hi m"o"m.txt`, //scmbreeze fails these with `hi m"o\`
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.chunk), func(t *testing.T) {
			gotPath, gotOrigPath, err := extractFilePaths(tc.chunk)
			if err != nil {
				t.Fatalf("extractFilePaths(%s): unexpected error: %v", tc.chunk, err)
			}
			if gotPath != tc.wantPath {
				t.Errorf("extractFilePaths(%s): expected path %s, got %s", tc.chunk, tc.wantPath, gotPath)
			}
			if gotOrigPath != tc.wantOrigPath {
				t.Errorf("extractFilePaths(%s): expected origPath %s, got %s", tc.chunk, tc.wantOrigPath, gotOrigPath)
			}
		})
	}
}

func Test_extractChangeCodes(t *testing.T) {
	// $ git status --porcelain
	// A  HELLO.md
	//
	//	M script/benchmark
	//
	// ?? .travis.yml
	// ?? commands/status/process_test.go
	var testCases = []struct {
		chunk    []byte
		expected []gitstatus.ChangeType
	}{
		{
			[]byte("A  HELLO.md"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedNewFile,
			},
		},
		{
			[]byte(" M script/benchmark"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedModified,
			},
		},
		{
			[]byte("?? .travis.yml"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUntracked,
			},
		},
		{
			[]byte(" D deleted_file"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedDeleted,
			},
		},
		{
			[]byte("R  after\x00before"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedRenamed,
			},
		},
		{
			[]byte("C  after\x00before"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedCopied,
			},
		},
		{
			[]byte("AM added_then_modified_file"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedNewFile,
				gitstatus.ChangeUnstagedModified,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.chunk[:]), func(t *testing.T) {
			actual := extractChangeCodes(tc.chunk)
			if !slices.Equal(actual, tc.expected) {
				t.Fatalf("processChange('%s'): expected %+v, actual %+v",
					tc.chunk, tc.expected, actual)
			}
		})
	}
}

func TestExtractBranch(t *testing.T) {
	// Examples of stuff we will want to parse:
	//
	//	## Initial commit on master
	//	## master
	//	## master...origin/master
	//	## master...origin/master [ahead 1]
	var testCases = []struct {
		chunk    []byte
		expected gitstatus.BranchInfo
	}{
		{
			[]byte("## master"),
			gitstatus.BranchInfo{Name: "master", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## GetUpGetDown09-11JokeInYoTown"),
			gitstatus.BranchInfo{Name: "GetUpGetDown09-11JokeInYoTown", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## master...origin/master"),
			gitstatus.BranchInfo{Name: "master", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## upstream...upstream/master"),
			gitstatus.BranchInfo{Name: "upstream", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## master...origin/master [ahead 1]"),
			gitstatus.BranchInfo{Name: "master", CommitsAhead: 1, CommitsBehind: 0},
		},
		{
			[]byte("## upstream...upstream/master [behind 3]"),
			gitstatus.BranchInfo{Name: "upstream", CommitsAhead: 0, CommitsBehind: 3},
		},
		{
			[]byte("## upstream...upstream/master [ahead 5, behind 3]"),
			gitstatus.BranchInfo{Name: "upstream", CommitsAhead: 5, CommitsBehind: 3},
		},
		{
			[]byte("## Initial commit on master"),
			gitstatus.BranchInfo{Name: "master", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## No commits yet on master"),
			gitstatus.BranchInfo{Name: "master", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## 3.0...origin/3.0 [ahead 1]"),
			gitstatus.BranchInfo{Name: "3.0", CommitsAhead: 1, CommitsBehind: 0},
		},
		{
			[]byte("## HEAD (no branch)"),
			gitstatus.BranchInfo{Name: "HEAD (no branch)", CommitsAhead: 0, CommitsBehind: 0},
		},
		{
			[]byte("## HEAD (no branch)UU both_modified.txt"),
			gitstatus.BranchInfo{Name: "HEAD (no branch)", CommitsAhead: 0, CommitsBehind: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.chunk[:]), func(t *testing.T) {
			actual, err := ExtractBranch(tc.chunk)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("processBranch('%s'): expected %v, actual %v",
					tc.chunk, tc.expected, actual)
			}
		})
	}
}
