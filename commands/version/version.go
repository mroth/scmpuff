package version

import (
	"fmt"

	"github.com/mroth/scmpuff/vendor/_nuts/github.com/spf13/cobra"
)

// the name of this software
const NAME string = "scmpuff"

// the version of this software
const VERSION string = "0.0.2"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long:  `All software has versions. This is ours.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(NAME, VERSION)
	},
}
