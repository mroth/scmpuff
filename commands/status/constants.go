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

var groupColorMap = map[StatusGroup]string{
	Staged:    "33m",
	Unmerged:  "31m",
	Unstaged:  "32m",
	Untracked: "36m",
}
