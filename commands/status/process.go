package status

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Process takes the raw output of `git status --porcelain -b -z` and turns it
// into a structured data type.
func Process(gitStatusOutput []byte, root string) *StatusList {
	// initialize a statuslist to hold the results
	results := NewStatusList()

	// put the output into a bufferreader+scanner so we can consume it iteratively
	scanner := bufio.NewScanner(bytes.NewReader(gitStatusOutput))

	// the scanner needs a custom split function for splitting on NUL
	scanNul := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i, b := range data {
			if b == '\x00' {
				return i + 1, data[:i], nil
			}
		}
		return 0, nil, nil
	}
	scanner.Split(scanNul)

	// branch output is first line
	if !scanner.Scan() {
		log.Println("Failed to read buffer when expecting branch status")
		log.Fatal(scanner.Err())
	}
	branchBytes := scanner.Bytes()
	results.branch = ExtractBranch(branchBytes)

	// give ProcessChanges the scanner and let it handle the rest
	// (it does complicated stuff so it needs the entire scanner)
	for _, r := range ProcessChanges(scanner, root) {
		results.groups[r.group].items = append(results.groups[r.group].items, r)
	}

	return results
}

// ExtractBranch handles parsing the branch status from `status --porcelain -b`.
//
// Examples of stuff we will want to parse:
//
// 		## Initial commit on master
// 		## master
// 		## master...origin/master
// 		## master...origin/master [ahead 1]
//
func ExtractBranch(bs []byte) *BranchInfo {
	b := BranchInfo{}

	b.name = decodeBranchName(bs)
	b.ahead, b.behind = decodeBranchPosition(bs)

	return &b
}

func decodeBranchName(bs []byte) string {
	re := regexp.MustCompile(`^## (?:Initial commit on )?([^ \.]+)`)
	m := re.FindSubmatch(bs)
	if m == nil {
		log.Fatalf("Failed to parse branch name for output: [%s]", bs)
	}

	return string(m[1])
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

// basically a StatusItem minus the file information, for now just being
// used to get results from the change code processing...
// This could probably be encapsulated in StatusItem itself, but wary of adding
// more nesting...
type change struct {
	msg   string
	col   ColorGroup
	group StatusGroup
}

/*
ProcessChanges takes `git status -z` output and returns all status items.

(Note: in our case, we actually use `git status -bz` and remove branch header
when we process it earlier, but the results are binary identical.)

This is a complicated process because the format is weird. Each line is a
variable length number of columns (2-3), but the separator for 1-2 is a space
(but the content of columns can contain spaces too!), and the seperator for 2-3
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
func ProcessChanges(s *bufio.Scanner, root string) (results []*StatusItem) {

	// Before we process any changes, get the Current Working Directory.
	// We're going to need use to calculate absolute and relative filepaths for
	// every change, so we get it once now and pass it along.
	// If for some reason this fails (?), fallback to the git worktree root.
	wd, err := os.Getwd()
	if err != nil {
		wd = root
	}

	for s.Scan() {
		chunk := s.Bytes()
		// ...if chunk represents a rename or copy op, need to append another chunk
		// to get the full change item, with NUL manually reinserted because scanner
		// will extract past it.
		if (chunk[0] == 'R' || chunk[0] == 'C') && s.Scan() {
			chunk = append(chunk, '\x00')
			chunk = append(chunk, s.Bytes()...)
		}
		results = append(results, processChange(chunk, wd, root)...)
	}

	return
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
func processChange(chunk []byte, wd, root string) (results []*StatusItem) {

	absolutePath, relativePath := extractFile(chunk, root, wd)

	for _, c := range extractChangeCodes(chunk) {
		result := &StatusItem{
			msg:         c.msg,
			col:         c.col,
			group:       c.group,
			fileAbsPath: absolutePath,
			fileRelPath: relativePath,
		}
		results = append(results, result)
	}

	if len(results) < 1 {
		log.Fatalf(`
Failed to decode git status change code for chunk: [%s]
Please file a bug including this error message as well as the output of:

git status --porcelain

You can file the bug at: https://github.com/mroth/scmpuff/issues/
		`, chunk)
	}
	return results
}

/*
extractFile extracts the filename from a status change, and determines the
absolute and display paths.

 - root: the absolute path to the git working tree
 - wd: current working directory path
*/
func extractFile(chunk []byte, root, wd string) (absPath, relPath string) {
	// file identifier starts at pos4 and continues to EOL
	filePortion := chunk[3:len(chunk)]
	files := bytes.SplitN(filePortion, []byte{'\x00'}, 2)

	n := len(files)
	switch {
	case n < 1:
		log.Fatalf("tried to process a change chunk with no file")
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
	 TODO: REPLICATE THIS LOGIC, INSTEAD OF STUFF PORTED FROM SCM_BREEZE.

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
func extractChangeCodes(chunk []byte) []*change {
	x := rune(chunk[0])
	y := rune(chunk[1])

	var changes []*change
	if p := decodePrimaryChangeCode(x, y); p != nil {
		changes = append(changes, p)
	}
	if s := decodeSecondaryChangeCode(x, y); s != nil {
		changes = append(changes, s)
	}
	return changes
}

func decodePrimaryChangeCode(x, y rune) *change {
	switch {
	case x == 'D' && y == 'D': //DD
		return &change{
			"   both deleted",
			del,
			Unmerged,
		}
	case x == 'A' && y == 'U': //AU
		return &change{
			"    added by us",
			neu,
			Unmerged,
		}
	case x == 'U' && y == 'D': //UD
		return &change{
			"deleted by them",
			del,
			Unmerged,
		}
	case x == 'U' && y == 'A': //UA
		return &change{
			"  added by them",
			neu,
			Unmerged,
		}
	case x == 'D' && y == 'U': //DU
		return &change{
			"  deleted by us",
			del,
			Unmerged,
		}
	case x == 'A' && y == 'A': //AA
		return &change{
			"     both added",
			neu,
			Unmerged,
		}
	case x == 'U' && y == 'U': //UU
		return &change{
			"  both modified",
			mod,
			Unmerged,
		}
	case x == 'M': //          //M.
		return &change{
			"  modified",
			mod,
			Staged,
		}
	case x == 'A': //          //A.
		return &change{
			"  new file",
			neu,
			Staged,
		}
	case x == 'D': //          //D.
		return &change{
			"   deleted",
			del,
			Staged,
		}
	case x == 'R': //          //R.
		return &change{
			"   renamed",
			ren,
			Staged,
		}
	case x == 'C': //          //C.
		return &change{
			"    copied",
			cpy,
			Staged,
		}
	case x == 'T': //          //T.
		return &change{
			"typechange",
			typ,
			Staged,
		}
	case x == '?' && y == '?': //??
		return &change{
			" untracked",
			unt,
			Untracked,
		}
	}

	return nil
}

func decodeSecondaryChangeCode(x, y rune) *change {
	switch {
	// TODO: fix the below and restore now that my cluelessness about these being
	// seperate statuses is reflected.
	//
	// So here's the thing, below case should never match, because [R.] earlier
	// is going to nab it.  So I'm assuming it's an oversight in the script.
	//
	// it was introduced to scm_breeze in:
	//   https://github.com/ndbroadbent/scm_breeze/pull/145/files
	//
	// case x == 'R' && y == 'M': //RM
	case x != 'R' && y == 'M': //[!R]M
		return &change{
			"  modified",
			mod,
			Unstaged,
		}
	case y == 'D' && x != 'D' && x != 'U': //[!D!U]D
		// Don't show deleted 'y' during a merge conflict.
		return &change{
			"   deleted",
			del,
			Unstaged,
		}
	case y == 'T': //.T
		return &change{
			"typechange",
			typ,
			Unstaged,
		}
	}

	return nil
}
