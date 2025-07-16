package status

import "github.com/kmatt/scmpuff/internal/gitstatus"

// Color represents an ANSI color code
type Color string

// Base color constants with full ANSI escape sequences
const (
	RedColor     Color = "\033[0;31m"
	GreenColor   Color = "\033[0;32m"
	YellowColor  Color = "\033[0;33m"
	BlueColor    Color = "\033[0;34m"
	MagentaColor Color = "\033[0;35m"
	CyanColor    Color = "\033[0;36m"

	BoldColor  Color = "\033[1m"
	DimColor   Color = "\033[2;37m" // dim white
	ResetColor Color = "\033[0m"
)

// Semantic color mappings for different change states
var stateColors = map[gitstatus.ChangeState]Color{
	gitstatus.NewState:         YellowColor,
	gitstatus.ModifiedState:    GreenColor,
	gitstatus.DeletedState:     RedColor,
	gitstatus.UntrackedState:   CyanColor,
	gitstatus.RenamedState:     BlueColor,
	gitstatus.CopiedState:      YellowColor,
	gitstatus.TypeChangedState: MagentaColor,
}

// Group color mappings for status groups
var groupColors = map[gitstatus.StatusGroup]Color{
	gitstatus.Staged:    YellowColor,
	gitstatus.Unmerged:  RedColor,
	gitstatus.Unstaged:  GreenColor,
	gitstatus.Untracked: CyanColor,
}

// Bold group colors for headers (arrows)
var groupBoldColors = map[gitstatus.StatusGroup]Color{
	gitstatus.Staged:    "\033[1;33m", // bold yellow
	gitstatus.Unmerged:  "\033[1;31m", // bold red
	gitstatus.Unstaged:  "\033[1;32m", // bold green
	gitstatus.Untracked: "\033[1;36m", // bold cyan
}
