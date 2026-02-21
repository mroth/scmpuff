package porcelainv1

import (
	"reflect"
	"slices"
	"testing"

	"github.com/mroth/scmpuff/internal/gitstatus"
)

func Test_extractChangeTypes(t *testing.T) {
	var testCases = []struct {
		xy       []byte
		expected []gitstatus.ChangeType
	}{
		{
			[]byte("A "), //[]byte("A  HELLO.md"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedNewFile,
			},
		},
		{
			[]byte(" M"), //[]byte(" M script/benchmark"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedModified,
			},
		},
		{
			[]byte("??"), //[]byte("?? .travis.yml"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUntracked,
			},
		},
		{
			[]byte(" D"), //[]byte(" D deleted_file"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedDeleted,
			},
		},
		{
			[]byte("R "), //[]byte("R  after\x00before"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedRenamed,
			},
		},
		{
			[]byte("C "), //[]byte("C  after\x00before"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedCopied,
			},
		},
		{
			[]byte(" A"), //[]byte(" A intent_to_add_file"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedNewFile,
			},
		},
		{
			// Verify UA produces only the unmerged type, not a spurious
			// ChangeUnstagedNewFile from the Y='A' secondary decoder.
			[]byte("UA"), //[]byte("UA added_by_them"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnmergedAddedThem,
			},
		},
		{
			// Verify AA produces only the unmerged type, not a spurious
			// ChangeUnstagedNewFile from the Y='A' secondary decoder.
			[]byte("AA"), //[]byte("AA both_added"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnmergedAddedBoth,
			},
		},
		{
			[]byte(" R"), //[]byte(" R renamed_in_worktree"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedRenamed,
			},
		},
		{
			[]byte(" C"), //[]byte(" C copied_in_worktree"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeUnstagedCopied,
			},
		},
		{
			[]byte("AM"), //[]byte("AM added_then_modified_file"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedNewFile,
				gitstatus.ChangeUnstagedModified,
			},
		},
		{
			// Compound code: staged modified + unstaged renamed
			[]byte("MR"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedModified,
				gitstatus.ChangeUnstagedRenamed,
			},
		},
		{
			// Compound code: staged copied + unstaged renamed
			[]byte("CR"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedCopied,
				gitstatus.ChangeUnstagedRenamed,
			},
		},
		{
			// Compound code: staged modified + unstaged copied
			[]byte("MC"),
			[]gitstatus.ChangeType{
				gitstatus.ChangeStagedModified,
				gitstatus.ChangeUnstagedCopied,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.xy), func(t *testing.T) {
			if len(tc.xy) != 2 {
				t.Fatalf("invalid test case: expected 2 bytes, actual %d", len(tc.xy))
			}

			actual := extractChangeTypes(tc.xy[0], tc.xy[1])
			if !slices.Equal(actual, tc.expected) {
				t.Fatalf("processChange('%s'): expected %+v, actual %+v",
					tc.xy, tc.expected, actual)
			}
		})
	}
}

func TestExtractBranch(t *testing.T) {
	var testCases = []struct {
		chunk    []byte
		expected gitstatus.BranchInfo
	}{
		{
			[]byte("## master"),
			gitstatus.BranchInfo{Name: "master"},
		},
		{
			[]byte("## feature/JIRA-1234_add-login+oauth2_support@2025-07-15"),
			gitstatus.BranchInfo{Name: "feature/JIRA-1234_add-login+oauth2_support@2025-07-15"},
		},
		{
			[]byte("## master...origin/master"),
			gitstatus.BranchInfo{Name: "master"},
		},
		{
			[]byte("## upstream...upstream/master"),
			gitstatus.BranchInfo{Name: "upstream"},
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
			gitstatus.BranchInfo{Name: "master"},
		},
		{
			[]byte("## No commits yet on master"),
			gitstatus.BranchInfo{Name: "master"},
		},
		{
			[]byte("## 3.0...origin/3.0 [ahead 1]"),
			gitstatus.BranchInfo{Name: "3.0", CommitsAhead: 1, CommitsBehind: 0},
		},
		{
			[]byte("## HEAD (no branch)"),
			gitstatus.BranchInfo{Name: "HEAD (no branch)"},
		},
		{
			// malformed header, missing LF and containing trailing entry
			[]byte("## HEAD (no branch)UU both_modified.txt"),
			gitstatus.BranchInfo{Name: "HEAD (no branch)"},
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
