package status

import (
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// single test to make sure everything gets stiched together properly, test
// actual cases in more specific methods
func Test_processChange(t *testing.T) {
	chunk := []byte("A  HELLO.md")
	res, err := processChange(chunk, "/tmp", "/tmp")
	actual := res[0]
	if err != nil {
		t.Fatal(err)
	}

	expectedChange := &changeType{
		msg:   "  new file",
		col:   neu,
		group: Staged,
	}

	t.Run("changeset", func(t *testing.T) {
		if actual.col != expectedChange.col ||
			actual.group != expectedChange.group ||
			actual.msg != expectedChange.msg {
			t.Fatal("changes did not match expected")
		}
	})

	t.Run("abspath", func(t *testing.T) {
		if actual.fileAbsPath != filepath.FromSlash("/tmp/HELLO.md") {
			t.Fatal("absolute path did not match expected")
		}
	})

	t.Run("relpath", func(t *testing.T) {
		if actual.fileRelPath == "" {
			t.Fatal("relative path was not present")
		}
	})
}

func Test_extractFile(t *testing.T) {
	var testCases = []struct {
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
			expectedAbs: filepath.FromSlash("/script/benchmark"),
			expectedRel: filepath.FromSlash("script/benchmark"),
		},
		{
			root:        "/tmp",
			wd:          "/tmp",
			chunk:       []byte(" M script/benchmark"),
			expectedAbs: filepath.FromSlash("/tmp/script/benchmark"),
			expectedRel: filepath.FromSlash("script/benchmark"),
		},
		{
			root:        "/tmp/foo/bar//",
			wd:          "/tmp/foo/bar/unicorn",
			chunk:       []byte("?? unicorn/magic/xxx"),
			expectedAbs: filepath.FromSlash("/tmp/foo/bar/unicorn/magic/xxx"),
			expectedRel: filepath.FromSlash("magic/xxx"),
		},
		{
			root:        "/tmp/foo/bar//",
			wd:          "/tmp/foo/bar/unicorn/magic",
			chunk:       []byte("?? narwhal/disco/yyy"),
			expectedAbs: filepath.FromSlash("/tmp/foo/bar/narwhal/disco/yyy"),
			expectedRel: filepath.FromSlash("../../narwhal/disco/yyy"),
		},
		{
			root:        "/tmp/foo",
			wd:          "/tmp/foo",
			chunk:       []byte("R  bar.txt\x00foo.txt"),
			expectedAbs: filepath.FromSlash("/tmp/foo/bar.txt"),
			expectedRel: filepath.FromSlash("foo.txt -> bar.txt"),
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
			expectedAbs: filepath.FromSlash("/tmp/foo/hi there mom.txt"),
			expectedRel: filepath.FromSlash("hi there mom.txt"),
		},
		{
			root:        "/tmp/foo",
			wd:          "/tmp/foo/bar",
			chunk:       []byte(`?? "x.txt`),
			expectedAbs: filepath.FromSlash(`/tmp/foo/"x.txt`),
			expectedRel: filepath.FromSlash(`../"x.txt`),
		},
		{
			root:        "/tmp/foo",
			wd:          "/tmp/foo",
			chunk:       []byte(`?? hi m"o"m.txt`),
			expectedAbs: filepath.FromSlash(`/tmp/foo/hi m"o"m.txt`), //scmbreeze fails these with `hi m"o\`
			expectedRel: filepath.FromSlash(`hi m"o"m.txt`),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("[root:%s],[wd:%s]", tc.root, tc.wd), func(t *testing.T) {
			actualAbs, actualRel, err := extractFile(tc.chunk, tc.root, tc.wd)
			if err != nil {
				t.Fatal(err)
			}

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
		expected []changeType
	}{
		{
			[]byte("A  HELLO.md"),
			[]changeType{
				changeStagedNewFile,
			},
		},
		{
			[]byte(" M script/benchmark"),
			[]changeType{
				changeUnstagedModified,
			},
		},
		{
			[]byte("?? .travis.yml"),
			[]changeType{
				changeUntracked,
			},
		},
		{
			[]byte(" D deleted_file"),
			[]changeType{
				changeUnstagedDeleted,
			},
		},
		{
			[]byte("R  after\x00before"),
			[]changeType{
				changeStagedRenamed,
			},
		},
		{
			[]byte("C  after\x00before"),
			[]changeType{
				changeStagedCopied,
			},
		},
		{
			[]byte("AM added_then_modified_file"),
			[]changeType{
				changeStagedNewFile,
				changeUnstagedModified,
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
		expected BranchInfo
	}{
		{
			[]byte("## master"),
			BranchInfo{name: "master", ahead: 0, behind: 0},
		},
		{
			[]byte("## GetUpGetDown09-11JokeInYoTown"),
			BranchInfo{name: "GetUpGetDown09-11JokeInYoTown", ahead: 0, behind: 0},
		},
		{
			[]byte("## master...origin/master"),
			BranchInfo{name: "master", ahead: 0, behind: 0},
		},
		{
			[]byte("## upstream...upstream/master"),
			BranchInfo{name: "upstream", ahead: 0, behind: 0},
		},
		{
			[]byte("## master...origin/master [ahead 1]"),
			BranchInfo{name: "master", ahead: 1, behind: 0},
		},
		{
			[]byte("## upstream...upstream/master [behind 3]"),
			BranchInfo{name: "upstream", ahead: 0, behind: 3},
		},
		{
			[]byte("## upstream...upstream/master [ahead 5, behind 3]"),
			BranchInfo{name: "upstream", ahead: 5, behind: 3},
		},
		{
			[]byte("## Initial commit on master"),
			BranchInfo{name: "master", ahead: 0, behind: 0},
		},
		{
			[]byte("## No commits yet on master"),
			BranchInfo{name: "master", ahead: 0, behind: 0},
		},
		{
			[]byte("## 3.0...origin/3.0 [ahead 1]"),
			BranchInfo{name: "3.0", ahead: 1, behind: 0},
		},
		{
			[]byte("## HEAD (no branch)"),
			BranchInfo{name: "HEAD (no branch)", ahead: 0, behind: 0},
		},
		{
			[]byte("## HEAD (no branch)UU both_modified.txt"),
			BranchInfo{name: "HEAD (no branch)", ahead: 0, behind: 0},
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

/*
// Test to verify https://github.com/mroth/scmpuff/issues/26.
//
// Leaving commented out since unlikely to encounter this exact issue again in
// future, and I'm not sure about importing the user's datafile into this repo.
//
// If needed again, the data file is attached to the issue as `output.txt`.

func TestBrokenProcessChanges(t *testing.T) {
	dat, err := ioutil.ReadFile("testdata/cjfuller_sample.dat")
	if err != nil {
		t.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(dat))
	s.Split(nulSplitFunc)
	actual, err := ProcessChanges(s, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != 270 { // `git status -s | wc -l` in replicated repo
		t.Errorf("expected %v changes, got %v", 270, len(actual))
	}
}
*/

func Test_calcPaths(t *testing.T) {
	type args struct {
		rootPath string
		root     string
		wd       string
	}
	tests := []struct {
		name        string
		args        args
		wantAbsPath string
		wantRelPath string
	}{
		{
			name: "everything in root",
			args: args{
				rootPath: "a.txt",
				root:     "/tmp/foo",
				wd:       "/tmp/foo",
			},
			wantAbsPath: filepath.FromSlash("/tmp/foo/a.txt"),
			wantRelPath: filepath.FromSlash("a.txt"),
		},
		{
			name: "change in subdir",
			args: args{
				rootPath: "bar/c.txt",
				root:     "/tmp/foo",
				wd:       "/tmp/foo",
			},
			wantAbsPath: filepath.FromSlash("/tmp/foo/bar/c.txt"),
			wantRelPath: filepath.FromSlash("bar/c.txt"),
		},
		{
			name: "change in parent to wd",
			args: args{
				rootPath: "a.txt",
				root:     "/tmp/foo",
				wd:       "/tmp/foo/bar",
			},
			wantAbsPath: filepath.FromSlash("/tmp/foo/a.txt"),
			wantRelPath: filepath.FromSlash("../a.txt"),
		},
		{
			name: "handle trailing slashes",
			args: args{
				rootPath: "a.txt",
				root:     "/tmp/foo/",
				wd:       "/tmp/foo/bar/",
			},
			wantAbsPath: filepath.FromSlash("/tmp/foo/a.txt"),
			wantRelPath: filepath.FromSlash("../a.txt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAbsPath, gotRelPath := calcPaths(tt.args.rootPath, tt.args.root, tt.args.wd)
			if gotAbsPath != tt.wantAbsPath {
				t.Errorf("calcPaths() gotAbsPath = %v, want %v", gotAbsPath, tt.wantAbsPath)
			}
			if gotRelPath != tt.wantRelPath {
				t.Errorf("calcPaths() gotRelPath = %v, want %v", gotRelPath, tt.wantRelPath)
			}
		})
	}
}
