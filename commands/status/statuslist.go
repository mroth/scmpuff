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
	groups map[StatusGroup]*FileGroup
}

// FileGroup is a bucket of all file StatusItems for a particular StatusGroup
type FileGroup struct {
	group StatusGroup
	desc  string
	items []*StatusItem
}

// StatusItem represents a single processed item of change from a 'git status'
type StatusItem struct {
	x, y  rune
	msg   string
	col   ColorGroup
	group StatusGroup
	file  string
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
	return []*FileGroup{sl.groups[0], sl.groups[1], sl.groups[2], sl.groups[3]}
	// use number literals rather than const names so that we can define the order
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
