package debug

import (
	"github.com/spf13/cobra"
)

// DebugCmd provides a debugging tools menu for scmpuff
var DebugCmd = &cobra.Command{
	Use:    "debug",
	Short:  "Debug tools for scmpuff",
	Long:   `Low-level debug commands not intended for normal usage.`,
	Hidden: true, // hide from help output
	// No Run function - this is a parent command that lists subcommands
}

func init() {
	DebugCmd.AddCommand(DumpCmd)
}
