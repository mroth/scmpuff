package porcelainv1

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mroth/porcelain/statusv1"
	"github.com/mroth/scmpuff/internal/gitstatus"
)

// Process takes the raw output of `git status --porcelain=v1 -b -z` and
// extracts the structured data.
//
// In the output the first segment of the output format is the git branch
// status, and the rest is the git status info.
func Process(gitStatusOutput []byte) (*gitstatus.StatusInfo, error) {
	// NOTE: in the future, we may wish to consume an io.Reader instead of
	// a byte slice, such that we can read from a pipe or other source
	// without needing to buffer the entire output in memory first.  For now,
	// we use a byte slice for reverse compatiblity with the existing tests
	// and architecture, so let's just wrap it in a bytes.Reader so the rest
	// of our code is ready for this change in the future.
	r := bytes.NewReader(gitStatusOutput)

	// parse the status output using the external porcelainv1 package
	status, err := statusv1.ParseZ(r)
	if err != nil {
		return nil, fmt.Errorf("porcelainv1: failed to parse status output: %w", err)
	}

	// we assume the first header is the branch status, based on git status command.
	if len(status.Headers) < 1 {
		return nil, fmt.Errorf("porcelainv1: failed to parse branch header: missing from status output")
	}
	branchHeader := status.Headers[0]
	branch, err := ExtractBranch([]byte(branchHeader))
	if err != nil {
		return nil, fmt.Errorf("porcelainv1: failed to parse branch header: %w", err)
	}

	// convert from porcelainv1 status entries to display status items
	items, err := ConvertEntries(status.Entries)
	if err != nil {
		return nil, fmt.Errorf("porcelainv1: failed to process status entries: %w", err)
	}

	return &gitstatus.StatusInfo{Branch: branch, Items: items}, nil
}

// ExtractBranch handles parsing the branch status from `status --porcelain -b`.
//
// Examples of stuff we will want to parse:
//
//	## Initial commit on master
//	## master
//	## master...origin/master
//	## master...origin/master [ahead 1]
func ExtractBranch(bs []byte) (gitstatus.BranchInfo, error) {
	name, err := decodeBranchName(bs)
	if err != nil {
		return gitstatus.BranchInfo{}, err
	}
	a, b := decodeBranchPosition(bs)

	return gitstatus.BranchInfo{
		Name:          name,
		CommitsAhead:  a,
		CommitsBehind: b,
	}, nil
}

func decodeBranchName(bs []byte) (string, error) {
	branchRegex := regexp.MustCompile(`^## (?:Initial commit on )?(?:No commits yet on )?(\S+?)(?:\.{3}|$)`)
	branchMatch := branchRegex.FindSubmatch(bs)
	if branchMatch != nil {
		return string(branchMatch[1]), nil
	}

	headRegex := regexp.MustCompile(`^## (HEAD \(no branch\))`)
	headMatch := headRegex.FindSubmatch(bs)
	if headMatch != nil {
		return string(headMatch[1]), nil
	}

	return "", fmt.Errorf("failed to parse branch name for output: [%s]", bs)
}

func decodeBranchPosition(bs []byte) (ahead, behind int) {
	reA := regexp.MustCompile(`\[ahead ?(\d+).*\]`)
	reB := regexp.MustCompile(`\[.*behind ?(\d+)\]`)

	mA := reA.FindSubmatch(bs)
	if mA != nil {
		ahead, _ = strconv.Atoi(string(mA[1]))
	}

	mB := reB.FindSubmatch(bs)
	if mB != nil {
		behind, _ = strconv.Atoi(string(mB[1]))
	}

	return
}

func ConvertEntries(entries []statusv1.Entry) ([]gitstatus.StatusItem, error) {
	// initial capacity is len(entries) because we expect at least one StatusItem per entry
	// but some entries can produce multiple StatusItems (e.g. staged + unstaged changes)
	// so we will grow the slice as needed.
	results := make([]gitstatus.StatusItem, 0, len(entries))

	for _, e := range entries {
		si, err := convertEntry(e)
		if err != nil {
			return results, fmt.Errorf("failed to convert entry: %w", err)
		}
		results = append(results, si...)
	}
	return results, nil
}

// convertEntry takes a single statusv1.Entry and converts it to []gitstatus.StatusItem(s).
// NOTE: A single Entry can produce multiple StatusItems, such as staged + unstaged changes.
func convertEntry(e statusv1.Entry) ([]gitstatus.StatusItem, error) {
	var results []gitstatus.StatusItem

	// we may get multiple change codes for a single entry.
	for _, c := range extractChangeTypes(byte(e.XY.X), byte(e.XY.Y)) {
		r := gitstatus.StatusItem{
			ChangeType: c,
			Path:       e.Path,
			OrigPath:   e.OrigPath,
		}
		results = append(results, r)
	}

	if len(results) < 1 {
		return nil, fmt.Errorf("unknown git status XY code: [%s]", e.XY)
	}
	return results, nil
}

/*
Extracts a git status "short code" into the proper UI "change" items we will
display in our status output.

Below documentation from git status:

	Ignored files are not listed, unless --ignored option is in effect, in
	which case XY are !!.

	X          Y     Meaning
	-------------------------------------------------
	          [AMD]   not updated
	M        [ MTD]   updated in index
	T        [ MTD]   type changed in index
	A        [ MTD]   added to index
	D                 deleted from index
	R        [ MTD]   renamed in index
	C        [ MTD]   copied in index
	[MTARC]           index and work tree matches
	[ MTARC]     M    work tree changed since index
	[ MTARC]     T    type changed in work tree since index
	[ MTARC]     D    deleted in work tree
	             R    renamed in work tree
	             C    copied in work tree
	-------------------------------------------------
	D           D    unmerged, both deleted
	A           U    unmerged, added by us
	U           D    unmerged, deleted by them
	U           A    unmerged, added by them
	D           U    unmerged, deleted by us
	A           A    unmerged, both added
	U           U    unmerged, both modified
	-------------------------------------------------
	?           ?    untracked
	!           !    ignored
	-------------------------------------------------
*/
func extractChangeTypes(x, y byte) []gitstatus.ChangeType {
	var changes []gitstatus.ChangeType
	if p, found := decodePrimaryChangeCode(x, y); found {
		changes = append(changes, p)
	}
	if s, found := decodeSecondaryChangeCode(x, y); found {
		changes = append(changes, s)
	}
	return changes
}

// decodePrimaryChangeCode returns the primary change code for a given status,
// or -1, false if it doesn't match any known codes.
func decodePrimaryChangeCode(x, y byte) (gitstatus.ChangeType, bool) {
	// unmerged cases are simple, only a single change UI is possible
	switch {
	case x == 'D' && y == 'D':
		return gitstatus.ChangeUnmergedDeletedBoth, true
	case x == 'A' && y == 'U':
		return gitstatus.ChangeUnmergedAddedUs, true
	case x == 'U' && y == 'D':
		return gitstatus.ChangeUnmergedDeletedThem, true
	case x == 'U' && y == 'A':
		return gitstatus.ChangeUnmergedAddedThem, true
	case x == 'D' && y == 'U':
		return gitstatus.ChangeUnmergedDeletedUs, true
	case x == 'A' && y == 'A':
		return gitstatus.ChangeUnmergedAddedBoth, true
	case x == 'U' && y == 'U':
		return gitstatus.ChangeUnmergedModifiedBoth, true
	case x == '?' && y == '?':
		return gitstatus.ChangeUntracked, true
	}

	// staged changes are all single X cases
	switch x {
	case 'M':
		return gitstatus.ChangeStagedModified, true
	case 'A':
		return gitstatus.ChangeStagedNewFile, true
	case 'D':
		return gitstatus.ChangeStagedDeleted, true
	case 'R':
		return gitstatus.ChangeStagedRenamed, true
	case 'C':
		return gitstatus.ChangeStagedCopied, true
	case 'T':
		return gitstatus.ChangeStagedType, true
	}

	return -1, false
}

// decodeSecondaryChangeCode returns the secondary change code for a given status,
// or -1, false if it doesn't match any known codes.
func decodeSecondaryChangeCode(x, y byte) (gitstatus.ChangeType, bool) {
	switch {
	case y == 'M': //.M
		return gitstatus.ChangeUnstagedModified, true
	// Don't show deleted 'y' during a merge conflict.
	case y == 'D' && x != 'D' && x != 'U': //[!D!U]D
		return gitstatus.ChangeUnstagedDeleted, true
	case y == 'T': //.T
		return gitstatus.ChangeUnstagedType, true
	// Don't show added 'y' during a merge conflict.
	case y == 'A' && x != 'A' && x != 'U': //[!A!U]A (intent-to-add, via git add -N)
		return gitstatus.ChangeUnstagedNewFile, true
	case y == 'R': //.R
		return gitstatus.ChangeUnstagedRenamed, true
	case y == 'C': //.C
		return gitstatus.ChangeUnstagedCopied, true
	}

	return -1, false
}
