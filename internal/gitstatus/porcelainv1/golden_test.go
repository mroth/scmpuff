package porcelainv1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mroth/scmpuff/internal/gitstatus"
)

func TestProcess(t *testing.T) {
	const (
		// localTestdata contains v1-only test fixtures specific to this parser.
		localTestdata = "testdata"

		// sharedTestdata contains regression fixtures from user-submitted debug
		// dumps, with both v1 and v2 porcelain formats for use across parsers.
		sharedTestdata = "../testdata"
	)

	var testCases = []struct {
		testdata   string
		sampleFile string
		want       *gitstatus.StatusInfo
	}{
		{
			testdata:   localTestdata,
			sampleFile: "process-untracked.porcelain-v1z.bin",
			want: &gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{
					Name:          "main",
					CommitsAhead:  0,
					CommitsBehind: 0,
				},
				Items: []gitstatus.StatusItem{
					{
						Path:       "a.txt",
						ChangeType: gitstatus.ChangeUntracked,
					},
					{
						Path:       "b.txt",
						ChangeType: gitstatus.ChangeUntracked,
					},
				},
			},
		},
		{
			testdata:   localTestdata,
			sampleFile: "process-changes.porcelain-v1z.bin",
			want: &gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{
					Name:          "main",
					CommitsAhead:  0,
					CommitsBehind: 0,
				},
				Items: []gitstatus.StatusItem{
					{
						Path:       "a.txt",
						ChangeType: gitstatus.ChangeStagedModified,
					},
					{
						Path:       "c.txt",
						ChangeType: gitstatus.ChangeStagedNewFile,
					},
					{
						Path:       "file with spaces.txt",
						OrigPath:   "b.txt",
						ChangeType: gitstatus.ChangeStagedRenamed,
					},
				},
			},
		},
		{
			// Regression test for #86: intent-to-add files (git add -N) produce
			// a [ A] status code that was previously unrecognized.
			// Fixture from user-submitted debug dump (jujutsu colocated repo).
			testdata:   sharedTestdata,
			sampleFile: "issue86.porcelain-v1z.bin",
			want: &gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{
					Name:          "main",
					CommitsAhead:  0,
					CommitsBehind: 0,
				},
				Items: []gitstatus.StatusItem{
					{
						Path:       "bar",
						ChangeType: gitstatus.ChangeUnstagedNewFile,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.sampleFile, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(tc.testdata, tc.sampleFile))
			if err != nil {
				t.Fatal(err)
			}

			got, err := Process(data)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Process() mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
