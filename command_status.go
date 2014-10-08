package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// CommandStatus processes 'git status --porcelain', and exports numbered
// env variables that contain the path of each affected file.
// Output is also more concise than standard 'git status'.
//
// Call with optional <group> parameter to filter by modification state:
// 1 || Staged,  2 || Unmerged,  3 || Unstaged,  4 || Untracked
func CommandStatus() *cobra.Command {

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Set and display numbered git status",
		Long: `
Processes 'git status --porcelain', and exports numbered env variables that
contain the path of each affected file.
Output is also more concise than standard 'git status'.
    `,
		Run: func(cmd *cobra.Command, args []string) {
			runStatus()
		},
	}

	// --relative
	// statusCmd.Flags().BoolVarP(
	// 	&expandRelative,
	// 	"relative",
	// 	"r",
	// 	false,
	// 	"TODO: DESCRIPTION HERE YO",
	// )

	return statusCmd
}

// StatusGroup encapsulates constants for mapping group status
type StatusGroup int

// constants representing an enum of all possible StatusGroups
const (
	Staged StatusGroup = iota
	Unmerged
	Unstaged
	Untracked
)

// ColorGroup encapsulates constants for mapping color output categories
type ColorGroup int

const (
	rst ColorGroup = iota
	del
	mod
	neu //'new' is reserved in Go
	ren
	cpy
	typ
	unt
	dark
	branch
	header
)

var colorMap = map[ColorGroup]string{
	rst:    "\033[0m",
	del:    "\033[0;31m",
	mod:    "\033[0;32m",
	neu:    "\033[0;33m",
	ren:    "\033[0;34m",
	cpy:    "\033[0;33m",
	typ:    "\033[0;35m",
	unt:    "\033[0;36m",
	dark:   "\033[2;37m",
	branch: "\033[1m",
	header: "\033[0m",
}

var groupColorMap = map[StatusGroup]string{
	Staged:    "33m",
	Unmerged:  "31m",
	Unstaged:  "32m",
	Untracked: "36m",
}

// StatusItem represents a single processed item of change from a 'git status'
type StatusItem struct {
	x, y  rune
	msg   string
	col   ColorGroup
	group StatusGroup
	file  string
}

// StatusList gives us a data structure to store all items of a git status
// organized by what group they fall under.
//
// This is helpful because we want to pull them out by group later, and don't
// want to bear the cost of filtering then.
//
// It also helps us map closer to the program logic of the Ruby code from
// scm_breeze, so hopefully easier to port.
type StatusList struct {
	groups map[StatusGroup]*FileGroup
}

// FileGroup is a bucket of all file StatusItems for a particular StatusGroup
type FileGroup struct {
	group StatusGroup
	desc  string
	items []*StatusItem
}

// NewStatusList is a constructor that initializes a new StatusList so that it's
// ready to use.
func NewStatusList() *StatusList {
	return &StatusList{
		groups: map[StatusGroup]*FileGroup{
			Staged: &FileGroup{
				group: Staged,
				desc:  "Changes to be committed",
				items: make([]*StatusItem, 0),
			},
			Unmerged: &FileGroup{
				group: Unmerged,
				desc:  "Unmerged paths",
				items: make([]*StatusItem, 0),
			},
			Unstaged: &FileGroup{
				group: Unstaged,
				desc:  "Changes not staged for commit",
				items: make([]*StatusItem, 0),
			},
			Untracked: &FileGroup{
				group: Untracked,
				desc:  "Untracked files",
				items: make([]*StatusItem, 0),
			},
		},
	}
}

// Returns the groups of a StatusList in a specific order.
//
// Since you can't range over maps in sequential order, we hard code the order
// here.
//
// We already have the keys as a const enum, so we could replace the map with a
// slice and use the StatsGroup as the index value, but I think it's clearer to
// use a map there even if uneccessary.
//
// If we ever really need to look at the performance of this, it might be worth
// seeing if using arrays is much faster (doubt it will make a difference in our
// case however.)
func (sl StatusList) orderedGroups() []*FileGroup {
	// use number literals rather than const names so that we can define the order
	// via the const definition.
	return []*FileGroup{sl.groups[0], sl.groups[1], sl.groups[2], sl.groups[3]}
}

// Total file change items across *all* groups.
func (sl StatusList) numItems() int {
	var total int
	for _, g := range sl.groups {
		total += len(g.items)
	}
	return total
}

func runStatus() {
	// TODO: fail if not git repo
	// TODO: git clear vars

	// TODO run commands to get status and branch
	gitStatusOutput, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		log.Fatal(err)
	}

	// gitBranchOutput, err := exec.Command("git", "branch", "-v").Output()
	// if err == nil {
	// 	log.Fatal(err)
	// }

	// allocate a StatusList to hold the results
	results := NewStatusList()

	if len(gitStatusOutput) > 0 {
		// split the status output to get a list of changes as raw bytestrings
		changes := bytes.Split(bytes.Trim(gitStatusOutput, "\n"), []byte{'\n'})

		// process each item, and store the results
		for _, change := range changes {
			rs := processChange(change)
			results.groups[rs.group].items = append(results.groups[rs.group].items, rs)
		}
	}

	results.printStatus()
}

func processChange(c []byte) *StatusItem {
	x := rune(c[0])
	y := rune(c[1])
	file := string(c[3:len(c)])
	msg, col, group := decodeChangeCode(x, y, file)

	ccc := StatusItem{
		x:     x,
		y:     y,
		file:  file,
		msg:   msg,
		col:   col,
		group: group,
	}
	return &ccc
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
		return " Untracked", unt, Untracked
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

	panic("Failed to decode git status change code!")
}

func (sl StatusList) printStatus() {
	if sl.numItems() == 0 {
		fmt.Println(outBannerBranch("master", "") + outBannerNoChanges())
	} else {
		for _, fg := range sl.orderedGroups() {
			fg.print()
		}
	}
}

// Make string for first half of the status banner.
// TODO: includes branch name with diff status
func outBannerBranch(branchname, difference string) string {
	return fmt.Sprintf(
		"%s#%s On branch: %s%s%s  %s|  ",
		colorMap[dark], colorMap[rst], colorMap[branch],
		branchname, difference,
		colorMap[dark],
	)
}

// If no changes, just display green no changes message (TODO: ?? and exit here)
func outBannerNoChanges() string {
	return fmt.Sprintf(
		"\033[0;32mNo changes (working directory clean)%s",
		colorMap[rst],
	)
}

// Output an entire filegroup to the screen
// TODO: format me and make me pretty
// TODO: have me return []files or whatever for later env setting
func (fg FileGroup) print() {
	if len(fg.items) > 0 {
		fg.printHeader()

		for _, i := range fg.items {
			i.printItem()
		}
	}
}

func (fg FileGroup) printHeader() {
	// heading := fg.desc
	cArrw := fmt.Sprintf("\033[1;%s", groupColorMap[fg.group])
	cHash := fmt.Sprintf("\033[0;%s", groupColorMap[fg.group])
	fmt.Printf(
		"%s➤%s %s\n%s#%s\n",
		cArrw, colorMap[header], fg.desc, cHash, colorMap[rst],
	)
}

func (si StatusItem) printItem() {
	// TODO: determine padding
	// TODO: find relative path
	// TODO: pretty print
	fmt.Println(si)
}
