package status

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Process takes the raw output of `git status --porcelain -b` and turns it into
// a well structured data type.
func Process(gitStatusOutput []byte, root string) *StatusList {
	results := NewStatusList()

	// split the status output to get a list of changes as raw bytestrings
	lines := bytes.Split(bytes.Trim(gitStatusOutput, "\n"), []byte{'\n'})

	// branch output is first line
	branchstr := lines[0]
	results.branch = ProcessBranch(branchstr)

	// status changes are everything else
	changes := lines[1:]

	// process each item, and store the results
	for _, change := range changes {
		for _, r := range ProcessChange(change, root) {
			results.groups[r.group].items = append(results.groups[r.group].items, r)
		}
	}

	return results
}

// ProcessBranch handles parsing the branch status from git status porcelain.
//
// Examples of stuff we will want to parse:
//
// 		## Initial commit on master
// 		## master
// 		## master...origin/master
// 		## master...origin/master [ahead 1]
//
func ProcessBranch(bs []byte) *BranchInfo {
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

// ProcessChange for a single item from a git status porcelain.
//
// Note some change items can have multiple statuses, so this returns a slice.
func ProcessChange(c []byte, root string) (results []*StatusItem) {
	x := rune(c[0])
	y := rune(c[1])

	absolutePath, relativePath := processFile(c, root)

	processedChanges := []*change{
		decodePrimaryChangeCode(x, y),
		decodeSecondaryChangeCode(x, y),
	}

	for _, c := range processedChanges {
		if c != nil {
			result := &StatusItem{
				msg:         c.msg,
				col:         c.col,
				group:       c.group,
				fileAbsPath: absolutePath,
				fileRelPath: relativePath,
			}
			results = append(results, result)
		}
	}

	if len(results) < 1 {
		log.Fatalf(`
Failed to decode git status change code for code: [%s]
Please file a bug including this error message as well as the output of:

git status --porcelain

You can file the bug at: https://github.com/mroth/scmpuff/issues/
		`, string(x)+string(y))
	}

	return results
}

/*
ProcessFile extracts the filename from a status change, and determines the
absolute and relative paths.

Parameters:
 - c: the raw bytes representing a status change
 - root: the absolute path to the git working tree
*/
func processFile(c []byte, root string) (absolutePath, relativePath string) {
	file := string(c[3:len(c)])
	absolutePath = filepath.Join(root, file)

	// get the current working directory and calculate relative path for file.
	// if for some reason this fails, fallback to absolute path.
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("DEBUG: ***wtf i didnt get working directory??***")
		relativePath = absolutePath
	} else {
		relativePath, err = filepath.Rel(wd, absolutePath)
		if err != nil {
			fmt.Println("DEBUG: ***wtf i couldn't calculate relative path??***")
			relativePath = absolutePath
		}
	}

	return
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
