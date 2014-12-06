package status

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var porcelainFiles bool

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

	// TODO
	// statusCmd.Flags().BoolVarP()
	// --aliases
	statusCmd.Flags().BoolVarP(
		&porcelainFiles,
		"filelist", "f", false,
		"include parseable filelist as first line of output",
	)

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
	// TODO: git clear vars (needs to be done in shellscript)

	// TODO run commands to get status and branch
	gitStatusOutput, err := exec.Command("git", "status", "--porcelain", "-b").Output()

	if err != nil {
		if err.Error() == "exit status 128" {
			fmt.Println("\033[0;31mNot a git repository (or any of the parent directories)")
			os.Exit(128)
		}
		// or, some other sort of error?
		log.Fatal(err)
	}

	results := Process(gitStatusOutput)
	results.printStatus(porcelainFiles)
}
