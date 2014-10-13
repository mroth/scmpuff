package status

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
)

// Process takes the raw output of `git status --porcelain -b` and turns it into
// a well structured data type.
func Process(gitStatusOutput []byte) *StatusList {
	results := NewStatusList()

	if len(gitStatusOutput) > 0 { //TODO: is this check necessary once we added the branch thing?
		// split the status output to get a list of changes as raw bytestrings
		lines := bytes.Split(bytes.Trim(gitStatusOutput, "\n"), []byte{'\n'})

		// branch output is first line
		branchstr := lines[0]
		results.branch = ProcessBranch(branchstr)

		// status changes are everything else
		changes := lines[1:]

		// process each item, and store the results
		for _, change := range changes {
			for _, r := range ProcessChange(change) {
				results.groups[r.group].items = append(results.groups[r.group].items, r)
			}
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

// ProcessChange for a single item from a git status porcelain.
//
// Note some change items can have multiple statuses, so this returns a slice.
func ProcessChange(c []byte) []*StatusItem {
	x := rune(c[0])
	y := rune(c[1])
	file := string(c[3:len(c)])

	var results []*StatusItem
	rp := decodePrimaryChangeCode(x, y, file)
	if rp != nil {
		results = append(results, rp)
	}
	rs := decodeSecondaryChangeCode(x, y, file)
	if rs != nil {
		results = append(results, rs)
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

func decodePrimaryChangeCode(x, y rune, file string) *StatusItem {
	switch {
	case x == 'D' && y == 'D': //DD
		return &StatusItem{
			"   both deleted",
			del,
			Unmerged,
			file,
		}
	case x == 'A' && y == 'U': //AU
		return &StatusItem{
			"    added by us",
			neu,
			Unmerged,
			file,
		}
	case x == 'U' && y == 'D': //UD
		return &StatusItem{
			"deleted by them",
			del,
			Unmerged,
			file,
		}
	case x == 'U' && y == 'A': //UA
		return &StatusItem{
			"  added by them",
			neu,
			Unmerged,
			file,
		}
	case x == 'D' && y == 'U': //DU
		return &StatusItem{
			"  deleted by us",
			del,
			Unmerged,
			file,
		}
	case x == 'A' && y == 'A': //AA
		return &StatusItem{
			"     both added",
			neu,
			Unmerged,
			file,
		}
	case x == 'U' && y == 'U': //UU
		return &StatusItem{
			"  both modified",
			mod,
			Unmerged,
			file,
		}
	case x == 'M': //          //M.
		return &StatusItem{
			"  modified",
			mod,
			Staged,
			file,
		}
	case x == 'A': //          //A.
		return &StatusItem{
			"  new file",
			neu,
			Staged,
			file,
		}
	case x == 'D': //          //D.
		return &StatusItem{
			"   deleted",
			del,
			Staged,
			file,
		}
	case x == 'R': //          //R.
		return &StatusItem{
			"   renamed",
			ren,
			Staged,
			file,
		}
	case x == 'C': //          //C.
		return &StatusItem{
			"    copied",
			cpy,
			Staged,
			file,
		}
	case x == 'T': //          //T.
		return &StatusItem{
			"typechange",
			typ,
			Staged,
			file,
		}
	case x == '?' && y == '?': //??
		return &StatusItem{
			" untracked",
			unt,
			Untracked,
			file,
		}
	}

	return nil
}

func decodeSecondaryChangeCode(x, y rune, file string) *StatusItem {
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
		return &StatusItem{
			"  modified",
			mod,
			Unstaged,
			file,
		}
	case y == 'D' && x != 'D' && x != 'U': //[!D!U]D
		// Don't show deleted 'y' during a merge conflict.
		return &StatusItem{
			"   deleted",
			del,
			Unstaged,
			file,
		}
	case y == 'T': //.T
		return &StatusItem{
			"typechange",
			typ,
			Unstaged,
			file,
		}
	}

	return nil
}
