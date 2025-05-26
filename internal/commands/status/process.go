package status

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// temporary structure until we rationalize StatusList which is a bit of a mess...
type statusInfo struct {
	branch BranchInfo
	items  []StatusItem
}

// Process takes the raw output of `git status --porcelain -b -z` and turns it
// into a structured data type.
//
// In the output of `git status --porcelain -b -z` the first segment of the output
// format is the git branch status, and the rest is the git status info.
func Process(gitStatusOutput []byte, root string) (*statusInfo, error) {
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
	statuses, err := ProcessChanges(remaining, root)
	if err != nil {
		return nil, err
	}

	return &statusInfo{branch: *branch, items: statuses}, nil
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
func ExtractBranch(bs []byte) (*BranchInfo, error) {
	name, err := decodeBranchName(bs)
	if err != nil {
		return nil, err
	}
	a, b := decodeBranchPosition(bs)

	return &BranchInfo{
		name:   name,
		ahead:  a,
		behind: b,
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
ProcessChanges takes `git status -z` output and returns all status items.

(Note: in our case, we actually use `git status -bz` and remove branch header
when we process it earlier, but the results are binary identical.)

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
func ProcessChanges(r io.Reader, root string) ([]StatusItem, error) {
	// Before we process any changes, get the Current Working Directory.
	// We're going to need use to calculate absolute and relative filepaths for
	// every change, so we get it once now and pass it along.
	// If for some reason this fails (?), fallback to the git worktree root.
	wd, err := os.Getwd()
	if err != nil {
		wd = root
	}

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
		statuses, err := processChange(chunk, wd, root)
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

// process change for a single item from a `git status -z`.
//
// Takes raw bytes representing status change from `git status --porcelain -z`,
// assumes that it has already been properly split away from the rest of the
// changes.
//
// See ProcessChanges (plural) for more details on that process.
//
// Note some change items can have multiple statuses, so this returns a slice.
func processChange(chunk []byte, wd, root string) ([]StatusItem, error) {
	var results []StatusItem
	absolutePath, relativePath, err := extractFile(chunk, root, wd)
	if err != nil {
		return nil, err
	}

	for _, c := range extractChangeCodes(chunk) {
		r := StatusItem{
			changeType:  c,
			fileAbsPath: absolutePath,
			fileRelPath: relativePath,
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

/*
extractFile extracts the filename from a status change, and determines the
absolute and display paths.

  - root: the absolute path to the git working tree
  - wd: current working directory path
*/
func extractFile(chunk []byte, root, wd string) (absPath, relPath string, err error) {
	// file identifier starts at pos4 and continues to EOL
	filePortion := chunk[3:]
	files := bytes.SplitN(filePortion, []byte{'\x00'}, 2)

	n := len(files)
	switch {
	case n < 1:
		err = errors.New("tried to process a change chunk with no file")
	case n > 1:
		toFile, fromFile := files[0], files[1]
		var toRelPath, fromRelPath string

		absPath, toRelPath = calcPaths(toFile, root, wd)
		_, fromRelPath = calcPaths(fromFile, root, wd)

		relPath = fmt.Sprintf("%s -> %s", fromRelPath, toRelPath)
	default:
		absPath, relPath = calcPaths(files[0], root, wd)
	}

	return
}

// given path of a file relative to git root, git root, and working directory,
// calculate the absolute path of the file on the system, and attempt to figure
// out its relative path to $CWD (if can't, fallback to absolute for both).
func calcPaths(rootPath []byte, root, wd string) (absPath, relPath string) {
	file := rootPath
	absPath = filepath.Join(root, string(file))
	relPath, err := filepath.Rel(wd, absPath)
	if err != nil {
		relPath = absPath
	}
	return
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
func extractChangeCodes(chunk []byte) []changeType {
	x := rune(chunk[0])
	y := rune(chunk[1])

	var changes []changeType
	if p := decodePrimaryChangeCode(x, y); p != nil {
		changes = append(changes, *p)
	}
	if s := decodeSecondaryChangeCode(x, y); s != nil {
		changes = append(changes, *s)
	}
	return changes
}

// decodePrimaryChangeCode returns the primary change code for a given status,
// or nil if it doesn't match any known codes.
func decodePrimaryChangeCode(x, y rune) *changeType {
	xy := string(x) + string(y)

	// unmerged cases are simple, only a single change UI is possible
	switch xy {
	case "DD":
		return &changeUnmergedDeletedBoth
	case "AU":
		return &changeUmmergedAddedUs
	case "UD":
		return &changeUnmergedDeletedThem
	case "UA":
		return &changeUnmergedAddedThem
	case "DU":
		return &changeUnmergedDeletedUs
	case "AA":
		return &changeUnmergedAddedBoth
	case "UU":
		return &changeUnmergedModifiedBoth
	case "??":
		return &changeUntracked
	}

	// staged changes are all single X cases
	switch x {
	case 'M':
		return &changeStagedModified
	case 'A':
		return &changeStagedNewFile
	case 'D':
		return &changeStagedDeleted
	case 'R':
		return &changeStagedRenamed
	case 'C':
		return &changeStagedCopied
	case 'T':
		return &changeStagedType
	}

	return nil
}

// decodeSecondaryChangeCode returns the secondary change code for a given status,
// or nil if it doesn't match any known codes.
func decodeSecondaryChangeCode(x, y rune) *changeType {
	switch {
	case y == 'M': //.M
		return &changeUnstagedModified
	// Don't show deleted 'y' during a merge conflict.
	case y == 'D' && x != 'D' && x != 'U': //[!D!U]D
		return &changeUnstagedDeleted
	case y == 'T': //.T
		return &changeUnstagedType
	}

	return nil
}

var (
	changeUnmergedDeletedBoth  = changeType{"   both deleted", del, Unmerged}
	changeUmmergedAddedUs      = changeType{"    added by us", neu, Unmerged}
	changeUnmergedDeletedThem  = changeType{"deleted by them", del, Unmerged}
	changeUnmergedAddedThem    = changeType{"  added by them", neu, Unmerged}
	changeUnmergedDeletedUs    = changeType{"  deleted by us", del, Unmerged}
	changeUnmergedAddedBoth    = changeType{"     both added", neu, Unmerged}
	changeUnmergedModifiedBoth = changeType{"  both modified", mod, Unmerged}
	changeUntracked            = changeType{" untracked", unt, Untracked}
	changeStagedModified       = changeType{"  modified", mod, Staged}
	changeStagedNewFile        = changeType{"  new file", neu, Staged}
	changeStagedDeleted        = changeType{"   deleted", del, Staged}
	changeStagedRenamed        = changeType{"   renamed", ren, Staged}
	changeStagedCopied         = changeType{"    copied", cpy, Staged}
	changeStagedType           = changeType{"typechange", typ, Staged}
	changeUnstagedModified     = changeType{"  modified", mod, Unstaged}
	changeUnstagedDeleted      = changeType{"   deleted", del, Unstaged}
	changeUnstagedType         = changeType{"typechange", typ, Unstaged}
)
