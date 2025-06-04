package status

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// A Renderer formats git status information for display to the screen.
type Renderer struct {
	branch       BranchInfo
	groupedItems map[StatusGroup][]StatusItem // re-organize items by their StatusGroup
	root, cwd    string                       // root and cwd are used to calculate paths for display
}

// NewRenderer creates a new Renderer instance from the provided StatusInfo.
//
// The git repository root and current working directory (cwd) must also be provided
// to correctly format the paths for display.
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
func (r *Renderer) Add(item StatusItem) {
	group := item.StatusGroup()
	r.groupedItems[group] = append(r.groupedItems[group], item)
}

// groupOrdering is the hardcoded list of the order StatusGroups should be displayed in
var groupOrdering = []StatusGroup{
	Staged,
	Unmerged,
	Unstaged,
	Untracked,
}

// orderedItems returns a slice of all StatusItems for the list regardless of what
// StatusGroup they belong to.
//
// However, we need to be careful to return them in the same order always.
func (r *Renderer) orderedItems() []StatusItem {
	var items []StatusItem
	for _, g := range groupOrdering {
		if groupItems, ok := r.groupedItems[g]; ok {
			items = append(items, groupItems...)
		}
	}

	return items
}

// numItems returns the count of StatusItems across all groups.
func (r *Renderer) numItems() int {
	var count int
	for _, g := range r.groupedItems {
		count += len(g)
	}
	return count
}

// Display renders the formatted status list designed for screen output to w.
//
// If includeParseData is true, the first line will be a machine parseable
// list of files to be used for environment variable expansion.
func (r *Renderer) Display(w io.Writer, includeParseData, includeStatusOutput bool) error {
	if includeParseData {
		if _, err := fmt.Fprintln(w, r.formatParseData()); err != nil {
			return fmt.Errorf("failed to write parse data: %w", err)
		}
	}

	if includeStatusOutput {
		if err := writeDisplayOutput(w, r); err != nil {
			return fmt.Errorf("failed to write display output: %w", err)
		}
	}

	return nil
}

func writeDisplayOutput(w io.Writer, r *Renderer) error {
	// Buffer writer due to many small writes
	b := bufio.NewWriter(w)

	// Print the banner
	fmt.Fprintln(b, r.formatBranchBanner())

	// Iterate through each group in the hardcoded order, for each group print
	// the header, then each item in that group, and finally the footer. For
	// each item, the display number is incremental across the entire list
	// (independent of group), so that the items can be referenced by number in
	// the shell script, with the first item being [1], second being [2], etc.
	itemNumber := 1
	for _, group := range groupOrdering {
		items := r.groupedItems[group]

		if len(items) > 0 {
			b.WriteString(formatHeaderForGroup(group))

			for _, item := range items {
				b.WriteString(r.formatStatusItemDisplay(item, itemNumber))
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

// formatParseData returns a machine readable string for environment variable parsing of file list in
// the scmpuff_status() shell script.
//
// Needs to be returned in same order that file lists are outputted to screen,
// otherwise env vars won't match UI.
func (r *Renderer) formatParseData() string {
	items := make([]string, r.numItems())
	for i, si := range r.orderedItems() {
		items[i] = si.AbsPath(r.root)
	}
	return strings.Join(items, "\t")
}

// formatBranchBanner formats the branch banner string to be used for printing.
//
// Banner string contains the branch information, as well as information about
// the branch status relative to upstream.
func (r *Renderer) formatBranchBanner() string {
	prelude := formatBranchBannerPrelude(r.branch)
	if r.numItems() == 0 {
		return prelude + bannerNoChanges()
	}
	return prelude + bannerChangeHeader()
}

// formatBranchBannerPrelude makes string for first half of the status banner.
func formatBranchBannerPrelude(b BranchInfo) string {
	diffStr := formatUpstreamDiffIndicator(b)
	var diffFormatted string
	if diffStr != "" {
		diffFormatted = fmt.Sprintf(
			"  %s|  %s%s%s",
			DimColor, YellowColor, diffStr, ResetColor,
		)
	}

	return fmt.Sprintf(
		"%s#%s On branch: %s%s%s  %s|  ",
		DimColor, ResetColor, BoldColor,
		b.Name, diffFormatted,
		DimColor,
	)
}

// formatUpstreamDiffIndicator formats the +1/-2 ahead/behind diff indicator for a branch relative to upstream
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
		ResetColor, DimColor, ResetColor, DimColor, ResetColor,
	)
}

// bannerNoChanges returns the no changes message when working directory is clean
func bannerNoChanges() string {
	return fmt.Sprintf(
		"%sNo changes (working directory clean)%s",
		GreenColor, ResetColor,
	)
}

// formatHeaderForGroup returns the display header string for a file group.
//
// Colorized version of something like this:
//
//	➤ Changes not staged for commit
//	#
func formatHeaderForGroup(group StatusGroup) string {
	groupColor := groupColors[group]
	groupBoldColor := groupBoldColors[group]
	return fmt.Sprintf(
		"%s➤%s %s\n%s#%s\n",
		groupBoldColor, ResetColor, group.Description(), groupColor, ResetColor,
	)
}

// formatFooterForGroup prints a final "#" for vertical padding
func formatFooterForGroup(group StatusGroup) string {
	groupColor := groupColors[group]
	return fmt.Sprintf("%s#%s\n", groupColor, ResetColor)
}

// formatStatusItemDisplay returns print string for an individual status item for a group.
//
// Colorized version of something like this:
//
//	#       modified: [1] commands/status/constants.go
func (r *Renderer) formatStatusItemDisplay(item StatusItem, displayNum int) string {
	// Get configured colors for the item display based on status group and state.
	groupColor := string(groupColors[item.StatusGroup()])
	stateColor := string(stateColors[item.state()])

	// For reasons lost to time, I originally decided to use a fixed width of 2
	// to pad the display number, so that entries 1-99 would align nicely.
	// scm_breeze uses a variable width of 1 or 2 depending on the number of
	// items in the list, but I went for consistency instead. At some point I
	// should probably look at what the rendering looks like with N>99 items.
	var padding string
	if displayNum < 10 {
		padding = " "
	}

	itemDisplayPath := item.DisplayPath(r.root, r.cwd)

	return fmt.Sprintf(
		"%s#%s     %s%s:%s%s [%s%d%s] %s%s%s\n",
		groupColor, ResetColor, stateColor, item.Message(), padding, DimColor,
		ResetColor, displayNum, DimColor, groupColor, itemDisplayPath, ResetColor,
	)
}
