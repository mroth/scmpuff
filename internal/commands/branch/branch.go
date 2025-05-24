package branch

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// CommandBranch lists git branches with numbered shortcuts.
// The first line of output, when --filelist is provided, will contain
// a TAB separated list of branch names suitable for environment expansion.
func CommandBranch() *cobra.Command {
	var optsFilelist bool

	var branchCmd = &cobra.Command{
		Use:   "branch",
		Short: "Display numbered git branches",
		Run: func(cmd *cobra.Command, args []string) {
			branches := gitBranchOutput()
			numbered, list := process(branches)

			if optsFilelist {
				fmt.Println(strings.Join(list, "\t"))
			}
			fmt.Print(numbered)
		},
	}

	branchCmd.Flags().BoolVarP(
		&optsFilelist,
		"filelist", "f", false,
		"include machine-parseable branch list",
	)

	return branchCmd
}

func gitBranchOutput() []byte {
	out, err := exec.Command("git", "branch").Output()
	if err != nil {
		if err.Error() == "exit status 128" {
			msg := "Not a git repository (or any of the parent directories)"
			fmt.Fprintf(os.Stderr, "\033[0;31m%s\033[0m\n", msg)
			os.Exit(128)
		}
		log.Fatal(err)
	}
	return out
}

// process takes raw `git branch` output and returns numbered output
// along with a slice of branch names in order.
func process(out []byte) (string, []string) {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	var starLine string
	var starBranch string
	var names []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}
		prefix := line[:2]
		name := strings.TrimSpace(line[2:])
		if prefix == "* " {
			starLine = line
			if !strings.HasPrefix(name, "(") {
				starBranch = name
			}
			continue
		}
		names = append(names, name)
	}

	var b strings.Builder
	var result []string
	n := 1
	if starLine != "" {
		if starBranch != "" {
			b.WriteString(fmt.Sprintf("* [%d] %s\n", n, starBranch))
			result = append(result, starBranch)
			n++
		} else {
			b.WriteString(starLine + "\n")
		}
	}
	for _, name := range names {
		b.WriteString(fmt.Sprintf("  [%d] %s\n", n, name))
		result = append(result, name)
		n++
	}
	return b.String(), result
}
