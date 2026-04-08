package status

import (
	"github.com/fatih/color"
	"github.com/mroth/scmpuff/internal/gitstatus"
)

// Color definitions using fatih/color for cross-platform terminal support.
var (
	RedColor           = color.New(color.FgRed)
	GreenColor         = color.New(color.FgGreen)
	YellowColor        = color.New(color.FgYellow)
	BlueColor          = color.New(color.FgBlue)
	MagentaColor       = color.New(color.FgMagenta)
	CyanColor          = color.New(color.FgCyan)
	BoldColor          = color.New(color.Bold)
	DimForegroundColor = color.New(color.Faint)
)

// Semantic color mappings for different change states
var stateColors = map[gitstatus.ChangeState]*color.Color{
	gitstatus.NewState:         YellowColor,
	gitstatus.ModifiedState:    GreenColor,
	gitstatus.DeletedState:     RedColor,
	gitstatus.UntrackedState:   CyanColor,
	gitstatus.RenamedState:     BlueColor,
	gitstatus.CopiedState:      YellowColor,
	gitstatus.TypeChangedState: MagentaColor,
}

// Group color mappings for status groups
var groupColors = map[gitstatus.StatusGroup]*color.Color{
	gitstatus.Staged:    YellowColor,
	gitstatus.Unmerged:  RedColor,
	gitstatus.Unstaged:  GreenColor,
	gitstatus.Untracked: CyanColor,
}

// Bold group colors for headers (arrows)
var groupBoldColors = map[gitstatus.StatusGroup]*color.Color{
	gitstatus.Staged:    color.New(color.FgYellow, color.Bold),
	gitstatus.Unmerged:  color.New(color.FgRed, color.Bold),
	gitstatus.Unstaged:  color.New(color.FgGreen, color.Bold),
	gitstatus.Untracked: color.New(color.FgCyan, color.Bold),
}
