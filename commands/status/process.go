package status

import (
	"log"
	"regexp"
	"strconv"
)

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
func ProcessChange(c []byte) *StatusItem {
	x := rune(c[0])
	y := rune(c[1])
	file := string(c[3:len(c)])
	msg, col, group := decodeChangeCode(x, y, file)

	return &StatusItem{
		x:     x,
		y:     y,
		file:  file,
		msg:   msg,
		col:   col,
		group: group,
	}
}

func decodeChangeCode(x, y rune, file string) (string, ColorGroup, StatusGroup) {
	switch {
	case x == 'D' && y == 'D': //DD
		return "   both deleted", del, Unmerged
	case x == 'A' && y == 'U': //AU
		return "    added by us", neu, Unmerged
	case x == 'U' && y == 'D': //UD
		return "deleted by them", del, Unmerged
	case x == 'U' && y == 'A': //UA
		return "  added by them", neu, Unmerged
	case x == 'D' && y == 'U': //DU
		return "  deleted by us", del, Unmerged
	case x == 'A' && y == 'A': //AA
		return "     both added", neu, Unmerged
	case x == 'U' && y == 'U': //UU
		return "  both modified", mod, Unmerged
	case x == 'M': //          //M.
		return "  modified", mod, Staged
	case x == 'A': //          //A.
		return "  new file", neu, Staged
	case x == 'D': //          //D.
		return "   deleted", del, Staged
	case x == 'R': //          //R.
		return "   renamed", ren, Staged
	case x == 'C': //          //C.
		return "    copied", cpy, Staged
	case x == 'T': //          //T.
		return "typechange", typ, Staged
	case x == '?' && y == '?': //??
		return " untracked", unt, Untracked
	// So here's the thing, below case should never match, because [R.] earlier
	// is going to nab it.  So I'm assuming it's an oversight in the script.
	//
	// it was introduced to scm_breeze in:
	//   https://github.com/ndbroadbent/scm_breeze/pull/145/files
	//
	// case x == 'R' && y == 'M': //RM
	case x != 'R' && y == 'M': //[!R]M
		return "  modified", mod, Unstaged
	case y == 'D' && y != 'D' && y != 'U': //[!D!U]D
		// Don't show deleted 'y' during a merge conflict.
		return "   deleted", del, Unstaged
	case y == 'T': //.T
		return "typechange", typ, Unstaged
	}

	log.Fatalf(`
Failed to decode git status change code for code: [%s]
Please file a bug including this error message as well as the output of:

git status --porcelain

You can file the bug at: https://github.com/mroth/scmpuff/issues/
`, string(x)+string(y))

	panic("...and you may ask yourself, well, how did I get here?")
}
