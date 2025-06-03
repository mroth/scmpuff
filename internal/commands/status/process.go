package status

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

// Process takes the raw output of `git status --porcelain=v1 -b -z` and
// extracts the structured data.
//
// In the output the first segment of the output format is the git branch
// status, and the rest is the git status info.
func Process(gitStatusOutput []byte) (*StatusInfo, error) {
	// NOTE: in the future, we may wish to consume an io.Reader instead of
	// a byte slice, such that we can read from a pipe or other source
	// without needing to buffer the entire output in memory first.  For now,
	// we use a byte slice for reverse compatiblity with the existing tests
	// and architecture, so let's just wrap it in a bytes.Reader so the rest
	// of our code is ready for this change in the future.
	r := bytes.NewReader(gitStatusOutput)

	// parse the first NUL seperated section of the git status output, which contains the branch
	branchBytes, remaining, err := cutFirstSegment(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read first segment of git status output: %w", err)
	}
	branch, err := ExtractBranch(branchBytes)
	if err != nil {
		return nil, err
	}

	// process the remaining NUL-separated sections, which contain the status items
	statuses, err := ProcessChanges(remaining)
	if err != nil {
		return nil, err
	}

	return &StatusInfo{Branch: branch, Items: statuses}, nil
}

// cutFirstSegment returns the first NUL-separated segment from r, and an io.Reader with the remainder of r.
func cutFirstSegment(r io.Reader) ([]byte, io.Reader, error) {
	br := bufio.NewReader(r)

	// read the first section (includes the NUL terminator)
	data, err := br.ReadBytes('\x00')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, err
	}

	// Strip trailing NUL
	data = bytes.TrimRight(data, "\x00")

	// br has already buffered some unread data, so we extract that and prepend it to the rest
	remaining := io.MultiReader(br, r)

	return data, remaining, nil
}

// ExtractBranch handles parsing the branch status from `status --porcelain -b`.
//
// Examples of stuff we will want to parse:
//
//	## Initial commit on master
//	## master
//	## master...origin/master
//	## master...origin/master [ahead 1]
func ExtractBranch(bs []byte) (BranchInfo, error) {
	name, err := decodeBranchName(bs)
	if err != nil {
		return BranchInfo{}, err
	}
	a, b := decodeBranchPosition(bs)

	return BranchInfo{
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

/*
ProcessChanges takes `git status --porcelain=v1 -z` output and returns all
status items.

(NOTE: in our case, we actually are using `git status --porcelain=v1 -z` and
removing the branch header when we process it earlier, prior to passing to this
function.)

This is a complicated process because the format is weird. Each line is a
variable length number of columns (2-3), but the separator for 1-2 is a space
(but the content of columns can contain spaces too!), and the separator for 2-3
is a NUL character (ASCII 0), *if* there is a third column. But here's where it
gets wacky: NUL is also the entry terminator (rather than a LF like in normal
porcelain mode)

Thankfully(?), column 1 which contains the status codes is a fixed length of two
bytes, and in theory the status codes contain enough secrets for us to determine
whether we should expect 2 or 3 columns (current hypothesis is we only get the
third column which is PATH2 when there is a "rename" operation). Sooo... we can
just read those two bytes and use that to determine how many NULs to scan to
until we have consumed a full entry.

We put up with this because it means no shell escaping, which should mean better
cross-platform support. Better hope some Windows people end up using it someday!
*/
func ProcessChanges(r io.Reader) ([]StatusItem, error) {
	s := bufio.NewScanner(r)
	s.Split(nulSplitFunc) // custom split function for splitting on NUL

	var results []StatusItem
	for s.Scan() {
		chunk := s.Bytes()
		// ...if chunk represents a rename or copy op, need to append another chunk
		// to get the full change item, with NUL manually reinserted because scanner
		// will extract past it.
		//
		// Note that the underlying slice from previous scanner.Bytes() MAY be
		// overridden by subsequent scans, so need to copy it to a new slice
		// first before scanning to get the next token.
		if chunk[0] == 'R' || chunk[0] == 'C' {
			composite := make([]byte, len(chunk))
			copy(composite, chunk)
			s.Scan()
			composite = append(composite, '\x00')
			composite = append(composite, s.Bytes()...)
			chunk = composite
		}
		statuses, err := processChange(chunk)
		if err != nil {
			return results, err
		}
		results = append(results, statuses...)
	}
	return results, nil
}

// custom split function for splitting on NUL
func nulSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, b := range data {
		if b == '\x00' {
			return i + 1, data[:i], nil
		}
	}
	return 0, nil, nil
}

// processChange for a single item chunk from a `git status --porcelain=v1 -z`.
//
// Note some change items can produce multiple statuses(!), so this returns a slice.
func processChange(chunk []byte) ([]StatusItem, error) {
	var results []StatusItem
	targetPath, origPath, err := extractFilePaths(chunk)
	if err != nil {
		return nil, err
	}

	for _, c := range extractChangeCodes(chunk) {
		r := StatusItem{
			ChangeType: c,
			Path:       targetPath,
			OrigPath:   origPath,
		}
		results = append(results, r)
	}

	if len(results) < 1 {
		return nil, fmt.Errorf(`
Failed to decode git status change code for chunk: [%s]
Please file a bug including this error message as well as the output of:

git status --porcelain

You can file the bug at: https://github.com/mroth/scmpuff/issues/
		`, chunk)
	}
	return results, nil
}

// extractFile extracts the file paths from a status change chunk
// origPath will be empty if the file was not renamed or copied.
func extractFilePaths(chunk []byte) (targetPath, origPath string, err error) {
	filePortion := chunk[3:]                              // file identifier starts at pos4 and continues to EOL
	files := bytes.SplitN(filePortion, []byte{'\x00'}, 2) // files split on NUL (-z option), 2 max

	switch len(files) {
	case 1:
		return string(files[0]), "", nil
	case 2:
		return string(files[0]), string(files[1]), nil
	default:
		return "", "", fmt.Errorf("extractFile: failed processing chunk, unexpected number of file fields: %d", len(files))
	}
}

/*
Extracts a git status "short code" into the proper UI "change" items we will
display in our status output.

Below documentation from git status:

	Ignored files are not listed, unless --ignored option is in effect, in
	which case XY are !!.

	X          Y     Meaning
	-------------------------------------------------
	          [MD]   not updated
	M        [ MD]   updated in index
	A        [ MD]   added to index
	D         [ M]   deleted from index
	R        [ MD]   renamed in index
	C        [ MD]   copied in index
	[MARC]           index and work tree matches
	[ MARC]     M    work tree changed since index
	[ MARC]     D    deleted in work tree
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
func extractChangeCodes(chunk []byte) []ChangeType {
	x := rune(chunk[0])
	y := rune(chunk[1])

	var changes []ChangeType
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
func decodePrimaryChangeCode(x, y rune) (ChangeType, bool) {
	// unmerged cases are simple, only a single change UI is possible
	switch {
	case x == 'D' && y == 'D':
		return ChangeUnmergedDeletedBoth, true
	case x == 'A' && y == 'U':
		return ChangeUnmergedAddedUs, true
	case x == 'U' && y == 'D':
		return ChangeUnmergedDeletedThem, true
	case x == 'U' && y == 'A':
		return ChangeUnmergedAddedThem, true
	case x == 'D' && y == 'U':
		return ChangeUnmergedDeletedUs, true
	case x == 'A' && y == 'A':
		return ChangeUnmergedAddedBoth, true
	case x == 'U' && y == 'U':
		return ChangeUnmergedModifiedBoth, true
	case x == '?' && y == '?':
		return ChangeUntracked, true
	}

	// staged changes are all single X cases
	switch x {
	case 'M':
		return ChangeStagedModified, true
	case 'A':
		return ChangeStagedNewFile, true
	case 'D':
		return ChangeStagedDeleted, true
	case 'R':
		return ChangeStagedRenamed, true
	case 'C':
		return ChangeStagedCopied, true
	case 'T':
		return ChangeStagedType, true
	}

	return -1, false
}

// decodeSecondaryChangeCode returns the secondary change code for a given status,
// or -1, false if it doesn't match any known codes.
func decodeSecondaryChangeCode(x, y rune) (ChangeType, bool) {
	switch {
	case y == 'M': //.M
		return ChangeUnstagedModified, true
	// Don't show deleted 'y' during a merge conflict.
	case y == 'D' && x != 'D' && x != 'U': //[!D!U]D
		return ChangeUnstagedDeleted, true
	case y == 'T': //.T
		return ChangeUnstagedType, true
	}

	return -1, false
}
