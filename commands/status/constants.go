package status

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

// TODO: Why does scm_breeze define these without the leading escape codes?
// Let's audit our usage and see if always used in same way, if so should be
// integrated into constants here...
//
// Ah, I see, sometimes used with bold, sometimes without.  Might be worth
// extracting that logic into a helper method to improve readability of print
// functions.
var groupColorMap = map[StatusGroup]string{
	Staged:    "33m",
	Unmerged:  "31m",
	Unstaged:  "32m",
	Untracked: "36m",
}
