package status

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
var stateColors = map[changeState]Color{
	NewState:         YellowColor,
	ModifiedState:    GreenColor,
	DeletedState:     RedColor,
	UntrackedState:   CyanColor,
	RenamedState:     BlueColor,
	CopiedState:      YellowColor,
	TypeChangedState: MagentaColor,
}

// Group color mappings for status groups
var groupColors = map[StatusGroup]Color{
	Staged:    YellowColor,
	Unmerged:  RedColor,
	Unstaged:  GreenColor,
	Untracked: CyanColor,
}

// Bold group colors for headers (arrows)
var groupBoldColors = map[StatusGroup]Color{
	Staged:    "\033[1;33m", // bold yellow
	Unmerged:  "\033[1;31m", // bold red
	Unstaged:  "\033[1;32m", // bold green
	Untracked: "\033[1;36m", // bold cyan
}
