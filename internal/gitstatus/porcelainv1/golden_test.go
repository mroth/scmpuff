package porcelainv1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kmatt/scmpuff/internal/gitstatus"
)

func TestProcess(t *testing.T) {
	var testCases = []struct {
		sampleFile string
		want       *gitstatus.StatusInfo
	}{
		{
			sampleFile: "process-untracked.porcelain-v1z.txt",
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
			sampleFile: "process-changes.porcelain-v1z.txt",
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
	}

	for _, tc := range testCases {
		t.Run(tc.sampleFile, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.sampleFile))
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
