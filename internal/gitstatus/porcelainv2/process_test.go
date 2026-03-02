package porcelainv2

import (
	"slices"
	"testing"

	"github.com/mroth/porcelain/statusv2"
	"github.com/mroth/scmpuff/internal/gitstatus"
)

func Test_decodeXY(t *testing.T) {
	testCases := []struct {
		xy       statusv2.XYFlag
		want     []gitstatus.ChangeType
		wantErr  bool
	}{
		{
			xy:       statusv2.XYFlag{X: statusv2.Added, Y: statusv2.Unmodified}, // [A.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedNewFile},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Modified}, // [.M]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedModified},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Deleted}, // [.D]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedDeleted},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Renamed, Y: statusv2.Unmodified}, // [R.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedRenamed},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Copied, Y: statusv2.Unmodified}, // [C.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedCopied},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Added}, // [.A]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedNewFile},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Renamed}, // [.R]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedRenamed},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Copied}, // [.C]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedCopied},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Modified, Y: statusv2.Unmodified}, // [M.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedModified},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Deleted, Y: statusv2.Unmodified}, // [D.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedDeleted},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.TypeChanged, Y: statusv2.Unmodified}, // [T.]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedType},
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.TypeChanged}, // [.T]
			want: []gitstatus.ChangeType{gitstatus.ChangeUnstagedType},
		},
		{
			// Compound: staged new file + unstaged modified
			xy:       statusv2.XYFlag{X: statusv2.Added, Y: statusv2.Modified}, // [AM]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedNewFile, gitstatus.ChangeUnstagedModified},
		},
		{
			// Compound: staged modified + unstaged renamed
			xy:       statusv2.XYFlag{X: statusv2.Modified, Y: statusv2.Renamed}, // [MR]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedModified, gitstatus.ChangeUnstagedRenamed},
		},
		{
			// Compound: staged copied + unstaged renamed
			xy:       statusv2.XYFlag{X: statusv2.Copied, Y: statusv2.Renamed}, // [CR]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedCopied, gitstatus.ChangeUnstagedRenamed},
		},
		{
			// Compound: staged modified + unstaged copied
			xy:       statusv2.XYFlag{X: statusv2.Modified, Y: statusv2.Copied}, // [MC]
			want: []gitstatus.ChangeType{gitstatus.ChangeStagedModified, gitstatus.ChangeUnstagedCopied},
		},
		{
			// Both unmodified — should never appear in a ChangedEntry
			xy:      statusv2.XYFlag{X: statusv2.Unmodified, Y: statusv2.Unmodified}, // [..]
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.xy.String(), func(t *testing.T) {
			got, err := decodeXY(tc.xy)
			if (err != nil) != tc.wantErr {
				t.Fatalf("decodeXY(%s): wantErr=%v, got err=%v", tc.xy, tc.wantErr, err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("decodeXY(%s): want %+v, got %+v", tc.xy, tc.want, got)
			}
		})
	}
}

func Test_decodeUnmergedXY(t *testing.T) {
	testCases := []struct {
		xy       statusv2.XYFlag
		want     gitstatus.ChangeType
		wantErr  bool
	}{
		{
			xy:       statusv2.XYFlag{X: statusv2.Deleted, Y: statusv2.Deleted}, // [DD]
			want: gitstatus.ChangeUnmergedDeletedBoth,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Added, Y: statusv2.UpdatedUnmerged}, // [AU]
			want: gitstatus.ChangeUnmergedAddedUs,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.UpdatedUnmerged, Y: statusv2.Deleted}, // [UD]
			want: gitstatus.ChangeUnmergedDeletedThem,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.UpdatedUnmerged, Y: statusv2.Added}, // [UA]
			want: gitstatus.ChangeUnmergedAddedThem,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Deleted, Y: statusv2.UpdatedUnmerged}, // [DU]
			want: gitstatus.ChangeUnmergedDeletedUs,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.Added, Y: statusv2.Added}, // [AA]
			want: gitstatus.ChangeUnmergedAddedBoth,
		},
		{
			xy:       statusv2.XYFlag{X: statusv2.UpdatedUnmerged, Y: statusv2.UpdatedUnmerged}, // [UU]
			want: gitstatus.ChangeUnmergedModifiedBoth,
		},
		{
			// Not a valid unmerged code
			xy:      statusv2.XYFlag{X: statusv2.Modified, Y: statusv2.Modified}, // [MM]
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.xy.String(), func(t *testing.T) {
			got, err := decodeUnmergedXY(tc.xy)
			if (err != nil) != tc.wantErr {
				t.Fatalf("decodeUnmergedXY(%s): wantErr=%v, got err=%v", tc.xy, tc.wantErr, err)
			}
			if got != tc.want {
				t.Errorf("decodeUnmergedXY(%s): want %v, got %v", tc.xy, tc.want, got)
			}
		})
	}
}

func Test_extractBranch(t *testing.T) {
	testCases := []struct {
		name     string
		input    *statusv2.BranchInfo
		want gitstatus.BranchInfo
		wantErr  bool
	}{
		{
			name: "normal branch",
			input: &statusv2.BranchInfo{
				OID:  "abc123",
				Head: "main",
			},
			want: gitstatus.BranchInfo{Name: "main"},
		},
		{
			name: "branch with upstream and ahead/behind",
			input: &statusv2.BranchInfo{
				OID:      "abc123",
				Head:     "feature",
				Upstream: "origin/feature",
				Ahead:    5,
				Behind:   3,
			},
			want: gitstatus.BranchInfo{Name: "feature", CommitsAhead: 5, CommitsBehind: 3},
		},
		{
			name: "detached HEAD",
			input: &statusv2.BranchInfo{
				OID:  "abc123",
				Head: "(detached)",
			},
			want: gitstatus.BranchInfo{Name: "HEAD (no branch)"},
		},
		{
			name:    "nil branch info",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractBranch(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("extractBranch: wantErr=%v, got err=%v", tc.wantErr, err)
			}
			if got != tc.want {
				t.Fatalf("extractBranch: want %+v, got %+v", tc.want, got)
			}
		})
	}
}
