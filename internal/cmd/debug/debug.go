package debug

import (
	"github.com/spf13/cobra"
)

// NewDebugCmd creates and returns the debug command
func NewDebugCmd() *cobra.Command {
	debugCmd := &cobra.Command{
		Use:    "debug",
		Short:  "Debug tools for scmpuff",
		Long:   `Low-level debug commands not intended for normal usage.`,
		Hidden: true, // hide from help output
		// No Run function - this is a parent command that lists subcommands
	}

	debugCmd.AddCommand(NewDumpCmd())
	return debugCmd
}
