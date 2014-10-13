package status

import "fmt"

// StatusList gives us a data structure to store all items of a git status
// organized by what group they fall under.
//
// This is helpful because we want to pull them out by group later, and don't
// want to bear the cost of filtering then.
//
// It also helps us map closer to the program logic of the Ruby code from
// scm_breeze, so hopefully easier to port.
type StatusList struct {
	branch *BranchInfo
	groups map[StatusGroup]*FileGroup
}

// BranchInfo contains all information needed about the active git branch, as
// well as its status relative to upstream commits.
type BranchInfo struct {
	name   string // name of the active branch
	ahead  int    // commit position relative to upstream, e.g. +1
	behind int    // commit position relative to upstream, e.g. -3
}

// FileGroup is a bucket of all file StatusItems for a particular StatusGroup
type FileGroup struct {
	group StatusGroup
	desc  string
	items []*StatusItem
}

// StatusItem represents a single processed item of change from a 'git status'
type StatusItem struct {
	msg   string
	col   ColorGroup
	group StatusGroup
	file  string
}

// NewStatusList is a constructor that initializes a new StatusList
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
	return []*FileGroup{sl.groups[0], sl.groups[1], sl.groups[2], sl.groups[3]}
	// uses number literals rather than const names so we can define the order
	// via the const definition.
}

// Total file change items across *all* groups.
func (sl StatusList) numItems() int {
	var total int
	for _, g := range sl.groups {
		total += len(g.items)
	}
	return total
}

func (sl StatusList) printStatus() {
	sl.printBanner()

	if sl.numItems() >= 1 {
		startNum := 1
		for _, fg := range sl.orderedGroups() {
			fg.print(startNum)
			startNum += len(fg.items)
		}
	}
}

func (sl StatusList) printBanner() {
	if sl.numItems() == 0 {
		fmt.Println(bannerBranch(sl.branch.name, "") + bannerNoChanges())
	} else {
		fmt.Println(bannerBranch(sl.branch.name, "") + bannerChangeHeader())
	}
}

// Make string for first half of the status banner.
// TODO: includes branch name with diff status
func bannerBranch(branchname, difference string) string {
	return fmt.Sprintf(
		"%s#%s On branch: %s%s%s  %s|  ",
		colorMap[dark], colorMap[rst], colorMap[branch],
		branchname, difference,
		colorMap[dark],
	)
}

func bannerChangeHeader() string {
	return fmt.Sprintf(
		"[%s*%s]%s => $e*\n%s#%s",
		colorMap[rst], colorMap[dark], colorMap[rst], colorMap[dark], colorMap[rst],
	)
}

// If no changes, just display green no changes message (TODO: ?? and exit here)
func bannerNoChanges() string {
	return fmt.Sprintf(
		"\033[0;32mNo changes (working directory clean)%s",
		colorMap[rst],
	)
}

// Output an entire filegroup to the screen
//
// The startNum argument tells us what number to start the listings at, it
// should probably be N+1 where N was the last number displayed (from previous
// outputted groups, that is.)
//
// TODO: have me return []files or whatever for later env setting?
func (fg FileGroup) print(startNum int) {
	if len(fg.items) > 0 {
		fg.printHeader()

		for n, i := range fg.items {
			i.printItem(startNum + n)
		}

		fg.printFooter()
	}
}

// Print the display header for a file group.
//
// Colorized version of something like this:
//
// 		➤ Changes not staged for commit
// 		#
//
func (fg FileGroup) printHeader() {
	cArrw := fmt.Sprintf("\033[1;%s", groupColorMap[fg.group])
	cHash := fmt.Sprintf("\033[0;%s", groupColorMap[fg.group])
	fmt.Printf(
		"%s➤%s %s\n%s#%s\n",
		cArrw, colorMap[header], fg.desc, cHash, colorMap[rst],
	)
}

// Print a final "#" for vertical padding
func (fg FileGroup) printFooter() {
	fmt.Printf("\033[0;%s#%s\n", groupColorMap[fg.group], colorMap[rst])
}

// Print an individual status item for a group.
//
// Colorized version of something like this:
//
//		#       modified: [1] commands/status/constants.go
//
// Arguments
// ---------
// displayNumber - the display number for the item, which should correspond to
//   the environment variable that will get set for it later ($eN).
//
func (si StatusItem) printItem(displayNum int) {

	// Determine padding size
	// scm_breeze does the following (Ruby code):
	//
	// 		padding = (@e < 10 && @changes.size >= 10) ? " " : ""
	//
	// instead of scm_breeze method, let's just fix the width at 2, so the output
	// is consistently spaced for e<=99, really we don't need to worry about the
	// one lost extra space when max(e)<10, I'd rather the spacing just be the
	// same.
	var padding string
	if displayNum < 10 {
		padding = " "
	}

	// TODO: find relative path
	relFile := si.file

	// TODO: if some submodules have changed, parse their summaries from long git
	// status the way scm_breeze does this requires a second call to git status,
	// which seems slow so maybe we will skip this for now?
	//
	// note to future self: format would add a final " %s" to output printf to
	// accomodate.

	groupCol := "\033[0;" + groupColorMap[si.group]
	fmt.Printf(
		"%s#%s     %s%s:%s%s [%s%d%s] %s%s%s\n",
		groupCol, colorMap[rst], colorMap[si.col], si.msg, padding, colorMap[dark],
		colorMap[rst], displayNum, colorMap[dark], groupCol, relFile, colorMap[rst],
	)
}
