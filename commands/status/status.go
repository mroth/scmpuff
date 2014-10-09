package status

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// CommandStatus processes 'git status --porcelain', and exports numbered
// env variables that contain the path of each affected file.
// Output is also more concise than standard 'git status'.
//
// Call with optional <group> parameter to filter by modification state:
// 1 || Staged,  2 || Unmerged,  3 || Unstaged,  4 || Untracked
func CommandStatus() *cobra.Command {

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Set and display numbered git status",
		Long: `
Processes 'git status --porcelain', and exports numbered env variables that
contain the path of each affected file.
Output is also more concise than standard 'git status'.
    `,
		Run: func(cmd *cobra.Command, args []string) {
			runStatus()
		},
	}

	// --relative
	// statusCmd.Flags().BoolVarP(
	// 	&expandRelative,
	// 	"relative",
	// 	"r",
	// 	false,
	// 	"TODO: DESCRIPTION HERE YO",
	// )

	return statusCmd
}

func runStatus() {
	// TODO: fail if not git repo
	// TODO: git clear vars

	// TODO run commands to get status and branch
	gitStatusOutput, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		log.Fatal(err)
	}

	// gitBranchOutput, err := exec.Command("git", "branch", "-v").Output()
	// if err == nil {
	// 	log.Fatal(err)
	// }

	// allocate a StatusList to hold the results
	results := NewStatusList()

	if len(gitStatusOutput) > 0 {
		// split the status output to get a list of changes as raw bytestrings
		changes := bytes.Split(bytes.Trim(gitStatusOutput, "\n"), []byte{'\n'})

		// process each item, and store the results
		for _, change := range changes {
			rs := processChange(change)
			results.groups[rs.group].items = append(results.groups[rs.group].items, rs)
		}
	}

	results.printStatus()
}
