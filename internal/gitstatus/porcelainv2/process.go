package porcelainv2

import (
	"bytes"
	"fmt"

	"github.com/mroth/porcelain/statusv2"
	"github.com/mroth/scmpuff/internal/gitstatus"
)

// Process takes the raw output of `git status --porcelain=v2 -b -z` and
// extracts the structured data.
//
// Unlike porcelain=v1, the v2 format provides structured branch information
// directly (no regex parsing needed) and uses typed entries that cleanly
// separate changed, renamed/copied, unmerged, and untracked files.
func Process(gitStatusOutput []byte) (*gitstatus.StatusInfo, error) {
	r := bytes.NewReader(gitStatusOutput)

	status, err := statusv2.ParseZ(r)
	if err != nil {
		return nil, fmt.Errorf("porcelainv2: failed to parse status output: %w", err)
	}

	branch, err := extractBranch(status.Branch)
	if err != nil {
		return nil, fmt.Errorf("porcelainv2: failed to extract branch info: %w", err)
	}

	items, err := convertEntries(status.Entries)
	if err != nil {
		return nil, fmt.Errorf("porcelainv2: failed to process status entries: %w", err)
	}

	return &gitstatus.StatusInfo{Branch: branch, Items: items}, nil
}

// extractBranch maps porcelain=v2 BranchInfo to our display BranchInfo.
//
// The v2 format provides branch name, upstream, and ahead/behind counts
// directly as structured fields, eliminating the regex parsing needed for v1.
func extractBranch(b *statusv2.BranchInfo) (gitstatus.BranchInfo, error) {
	if b == nil {
		return gitstatus.BranchInfo{}, fmt.Errorf("missing branch info from status output")
	}

	// v1 used "HEAD (no branch)" to indicate a detached HEAD, but v2 uses
	// "(detached)". Preserve previous display behavior for now - in the future,
	// "HEAD (detached)" might be more user-friendly.
	name := b.Head
	if name == "(detached)" {
		name = "HEAD (no branch)"
	}

	return gitstatus.BranchInfo{
		Name:          name,
		CommitsAhead:  b.Ahead,
		CommitsBehind: b.Behind,
	}, nil
}

// convertEntries converts porcelain=v2 typed entries to display StatusItems.
//
// A single entry can produce multiple StatusItems (e.g. staged + unstaged
// changes for the same file), so the result slice may be longer than the input.
func convertEntries(entries []statusv2.Entry) ([]gitstatus.StatusItem, error) {
	results := make([]gitstatus.StatusItem, 0, len(entries))

	for _, entry := range entries {
		switch e := entry.(type) {
		case statusv2.ChangedEntry:
			changes, err := decodeXY(e.XY)
			if err != nil {
				return nil, err
			}
			for _, c := range changes {
				results = append(results, gitstatus.StatusItem{ChangeType: c, Path: e.Path})
			}

		case statusv2.RenameOrCopyEntry:
			changes, err := decodeXY(e.XY)
			if err != nil {
				return nil, err
			}
			for _, c := range changes {
				results = append(results, gitstatus.StatusItem{ChangeType: c, Path: e.Path, OrigPath: e.Orig})
			}

		case statusv2.UnmergedEntry:
			c, err := decodeUnmergedXY(e.XY)
			if err != nil {
				return nil, err
			}
			results = append(results, gitstatus.StatusItem{ChangeType: c, Path: e.Path})

		case statusv2.UntrackedEntry:
			results = append(results, gitstatus.StatusItem{ChangeType: gitstatus.ChangeUntracked, Path: e.Path})

		case statusv2.IgnoredEntry:
			// Ignored files are not displayed in scmpuff status output.

		default:
			return nil, fmt.Errorf("unknown entry type: %T", e)
		}
	}
	return results, nil
}

// decodeXY converts a porcelain=v2 XY status code into change types.
//
// X represents staged (index) changes, Y represents unstaged (worktree)
// changes. A single XY code can produce 0–2 change types (e.g. "AM" produces
// both a staged new file and an unstaged modification).
//
// Unlike the v1 decoder, no merge-conflict guards are needed here because v2
// separates UnmergedEntry as a distinct type — unmerged XY codes never reach
// this function.
func decodeXY(xy statusv2.XYFlag) ([]gitstatus.ChangeType, error) {
	var changes []gitstatus.ChangeType

	// Staged (index) change
	switch xy.X {
	case statusv2.Modified:
		changes = append(changes, gitstatus.ChangeStagedModified)
	case statusv2.Added:
		changes = append(changes, gitstatus.ChangeStagedNewFile)
	case statusv2.Deleted:
		changes = append(changes, gitstatus.ChangeStagedDeleted)
	case statusv2.Renamed:
		changes = append(changes, gitstatus.ChangeStagedRenamed)
	case statusv2.Copied:
		changes = append(changes, gitstatus.ChangeStagedCopied)
	case statusv2.TypeChanged:
		changes = append(changes, gitstatus.ChangeStagedType)
	case statusv2.Unmodified:
		// No staged change — normal for single-sided entries like [.M].
	case statusv2.UpdatedUnmerged:
		// Should never appear in a ChangedEntry; v2 routes these to UnmergedEntry.
		return nil, fmt.Errorf("unexpected UpdatedUnmerged in staged change: [%s]", xy)
	}

	// Unstaged (worktree) change
	switch xy.Y {
	case statusv2.Modified:
		changes = append(changes, gitstatus.ChangeUnstagedModified)
	case statusv2.Deleted:
		changes = append(changes, gitstatus.ChangeUnstagedDeleted)
	case statusv2.TypeChanged:
		changes = append(changes, gitstatus.ChangeUnstagedType)
	case statusv2.Added:
		changes = append(changes, gitstatus.ChangeUnstagedNewFile)
	case statusv2.Renamed:
		changes = append(changes, gitstatus.ChangeUnstagedRenamed)
	case statusv2.Copied:
		changes = append(changes, gitstatus.ChangeUnstagedCopied)
	case statusv2.Unmodified:
		// No unstaged change — normal for single-sided entries like [A.].
	case statusv2.UpdatedUnmerged:
		// Should never appear in a ChangedEntry; v2 routes these to UnmergedEntry.
		return nil, fmt.Errorf("unexpected UpdatedUnmerged in unstaged change: [%s]", xy)
	}

	if len(changes) == 0 {
		return nil, fmt.Errorf("unknown git status XY code: [%s]", xy)
	}
	return changes, nil
}

// decodeUnmergedXY maps the 7 possible unmerged XY codes to change types.
func decodeUnmergedXY(xy statusv2.XYFlag) (gitstatus.ChangeType, error) {
	switch {
	case xy.X == statusv2.Deleted && xy.Y == statusv2.Deleted:
		return gitstatus.ChangeUnmergedDeletedBoth, nil
	case xy.X == statusv2.Added && xy.Y == statusv2.UpdatedUnmerged:
		return gitstatus.ChangeUnmergedAddedUs, nil
	case xy.X == statusv2.UpdatedUnmerged && xy.Y == statusv2.Deleted:
		return gitstatus.ChangeUnmergedDeletedThem, nil
	case xy.X == statusv2.UpdatedUnmerged && xy.Y == statusv2.Added:
		return gitstatus.ChangeUnmergedAddedThem, nil
	case xy.X == statusv2.Deleted && xy.Y == statusv2.UpdatedUnmerged:
		return gitstatus.ChangeUnmergedDeletedUs, nil
	case xy.X == statusv2.Added && xy.Y == statusv2.Added:
		return gitstatus.ChangeUnmergedAddedBoth, nil
	case xy.X == statusv2.UpdatedUnmerged && xy.Y == statusv2.UpdatedUnmerged:
		return gitstatus.ChangeUnmergedModifiedBoth, nil
	default:
		return 0, fmt.Errorf("unknown unmerged XY code: [%s]", xy)
	}
}
