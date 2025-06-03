package status

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// StatusInfo contains the information about git working tree status that is
// necessary to display the status in a user-friendly way (e.g. git status)
// It includes the branch information and a list of status items.
type StatusInfo struct {
	Branch BranchInfo
	Items  []StatusItem
}

// BranchInfo contains all information needed about the active git branch, as
// well as its status relative to upstream commits.
type BranchInfo struct {
	Name          string // name of the active branch
	CommitsAhead  int    // commit position relative to upstream, e.g. +1
	CommitsBehind int    // commit position relative to upstream, e.g. -3
}

// A Renderer formats git status information for display to the screen.
type Renderer struct {
	branch       BranchInfo
	groupedItems map[StatusGroup][]StatusItem // re-organize items by their StatusGroup
	root, cwd    string                       // root and cwd are used to calculate paths for display
}

func NewRenderer(info *StatusInfo, root, cwd string) (*Renderer, error) {
	if info == nil {
		return nil, fmt.Errorf("status info cannot be nil")
	}

	groupedItems := make(map[StatusGroup][]StatusItem)
	for _, item := range info.Items {
		group := item.StatusGroup()
		groupedItems[group] = append(groupedItems[group], item)
	}

	return &Renderer{
		branch:       info.Branch,
		groupedItems: groupedItems,
		root:         root,
		cwd:          cwd,
	}, nil
}

// Add appends a StatusItem to the Renderer, organizing it by its StatusGroup.
func (sl *Renderer) Add(item StatusItem) {
	group := item.StatusGroup()
	sl.groupedItems[group] = append(sl.groupedItems[group], item)
}

// groupOrdering is the hardcoded list of the order StatusGroups should be displayed in
var groupOrdering = []StatusGroup{
	Staged,
	Unmerged,
	Unstaged,
	Untracked,
}

// orderedItems will return a slice of all StatusItems for the list regardless of what
// StatusGroup they belong to.
//
// However, we need to be careful to return them in the same order always.
func (sl *Renderer) orderedItems() (items []StatusItem) {
	for _, g := range groupOrdering {
		if groupItems, ok := sl.groupedItems[g]; ok {
			items = append(items, groupItems...)
		}
	}

	return
}

// numItems returns the count of StatusItems across all groups.
func (sl *Renderer) numItems() int {
	var count int
	for _, g := range sl.groupedItems {
		count += len(g)
	}
	return count
}

// Displays the formatted status list designed for screen output to w.
//
// if `includeParseData` is true, the first line will be a machine parseable
// list of files to be used for environment variable expansion.
func (sl *Renderer) Display(w io.Writer, includeParseData, includeStatusOutput bool) error {
	if includeParseData {
		if _, err := fmt.Fprintln(w, sl.formatParseData()); err != nil {
			return fmt.Errorf("failed to write parse data: %w", err)
		}
	}

	if includeStatusOutput {
		if err := writeDisplayOutput(w, sl); err != nil {
			return fmt.Errorf("failed to write display output: %w", err)
		}
	}

	return nil
}

func writeDisplayOutput(w io.Writer, sl *Renderer) error {
	// buffer writer due to many small writes
	b := bufio.NewWriter(w)

	// print the banner
	fmt.Fprintln(b, sl.formatBranchBanner())

	// iterate through each group in the hardcoded order, for each group print
	// the header, then each item in that group, and finally the footer. For
	// each item, the display number is incremental across the entire list
	// (independent of group), so that the items can be referenced by number in
	// the shell script, with the first item being [1], second being [2], etc.
	itemNumber := 1
	for _, group := range groupOrdering {
		items := sl.groupedItems[group]

		if len(items) > 0 {
			b.WriteString(formatHeaderForGroup(group))

			for _, item := range items {
				b.WriteString(sl.formatStatusItemDisplay(item, itemNumber))
				itemNumber++
			}

			b.WriteString(formatFooterForGroup(group))
		}
	}

	// NOTE: Flush uses the errWriter pattern[1] and will return the first error
	// that was encountered while writing to the buffer, if any.
	//
	// [1]: https://go.dev/blog/errors-are-values
	return b.Flush()
}

// Machine readable string for environment variable parsing of file list in
// the scmpuff_status() shell script.
//
// Needs to be returned in same order that file lists are outputted to screen,
// otherwise env vars won't match UI.
func (sl *Renderer) formatParseData() string {
	items := make([]string, sl.numItems())
	for i, si := range sl.orderedItems() {
		items[i] = si.AbsPath(sl.root)
	}
	return strings.Join(items, "\t")
}

// Formats the branch banner string to be used for printing.
//
// Banner string contains the branch information, as well as information about
// the branch status relative to upstream.
func (sl Renderer) formatBranchBanner() string {
	if sl.numItems() == 0 {
		return formatBranchBannerPrelude(sl.branch) + bannerNoChanges()
	}
	return formatBranchBannerPrelude(sl.branch) + bannerChangeHeader()
}

// Make string for first half of the status banner.
func formatBranchBannerPrelude(b BranchInfo) string {
	diffStr := formatUpstreamDiffIndicator(b)
	var diffFormatted string
	if diffStr != "" {
		diffFormatted = fmt.Sprintf(
			"  %s|  %s%s%s",
			colorMap[dark], colorMap[neu], diffStr, colorMap[rst],
		)
	}

	return fmt.Sprintf(
		"%s#%s On branch: %s%s%s  %s|  ",
		colorMap[dark], colorMap[rst], colorMap[branch],
		b.Name, diffFormatted,
		colorMap[dark],
	)
}

// formats the +1/-2 ahead/behind diff indicator for a branch relative to upstream
func formatUpstreamDiffIndicator(b BranchInfo) string {
	switch {
	case b.CommitsAhead > 0 && b.CommitsBehind > 0:
		return fmt.Sprintf("+%d/-%d", b.CommitsAhead, b.CommitsBehind)
	case b.CommitsAhead > 0:
		return fmt.Sprintf("+%d", b.CommitsAhead)
	case b.CommitsBehind > 0:
		return fmt.Sprintf("-%d", b.CommitsBehind)
	default:
		return ""
	}
}

func bannerChangeHeader() string {
	return fmt.Sprintf(
		"[%s*%s]%s => $e*\n%s#%s",
		colorMap[rst], colorMap[dark], colorMap[rst], colorMap[dark], colorMap[rst],
	)
}

// If no changes, just display green no changes message
func bannerNoChanges() string {
	return fmt.Sprintf(
		"\033[0;32mNo changes (working directory clean)%s",
		colorMap[rst],
	)
}

// Returns the display header string for a file group.
//
// Colorized version of something like this:
//
//	➤ Changes not staged for commit
//	#
func formatHeaderForGroup(group StatusGroup) string {
	cArrw := fmt.Sprintf("\033[1;%s", groupColorMap[group])
	cHash := fmt.Sprintf("\033[0;%s", groupColorMap[group])
	return fmt.Sprintf(
		"%s➤%s %s\n%s#%s\n",
		cArrw, colorMap[header], group.Description(), cHash, colorMap[rst],
	)
}

// Print a final "#" for vertical padding
func formatFooterForGroup(group StatusGroup) string {
	return fmt.Sprintf("\033[0;%s#%s\n", groupColorMap[group], colorMap[rst])
}

// Returns print string for an individual status item for a group.
//
// Colorized version of something like this:
//
//	#       modified: [1] commands/status/constants.go
func (sl *Renderer) formatStatusItemDisplay(item StatusItem, displayNum int) string {
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

	// find relative path
	relFile := item.DisplayPath(sl.root, sl.cwd)

	// TODO: if some submodules have changed, parse their summaries from long git
	// status the way scm_breeze does this requires a second call to git status,
	// which seems slow so maybe we will skip this for now?
	//
	// note to future self: format would add a final " %s" to output printf to
	// accommodate.

	groupCol := "\033[0;" + groupColorMap[item.StatusGroup()]
	return fmt.Sprintf(
		"%s#%s     %s%s:%s%s [%s%d%s] %s%s%s\n",
		groupCol, colorMap[rst], item.Color(), item.Message(), padding, colorMap[dark],
		colorMap[rst], displayNum, colorMap[dark], groupCol, relFile, colorMap[rst],
	)
}
