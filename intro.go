package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var introCmd = &cobra.Command{
	Use:   "intro",
	Short: "Displays an introduction to scmpuff",
	Long:  `Displays an introduction to using scmpuff.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Hello there!

If you are just getting started, you probably want to make sure scmpuff is
automatically loaded in your shell, in order to do so check out 'scmpuff init',
but for the most part you will be adding a line to ~/.zshrc or ~/.bash_profile
that looks something like this:

  eval "$(scmpuff init -s)"

Once things are loaded, the most important function you will want to know about
is 'scmpuff_status', which is aliased to 'gs' for short.

This is a replacement for 'git status' that is pretty and shows you numbers next
to each filename, for example:

  $ gs
  # On branch: master  |  +1  |  [*] => $e*
  #
  ➤ Changes not staged for commit
  #
  #       modified:  [1] main.go
  #
  ➤ Untracked files
  #
  #      untracked:  [2] HELLO.txt
  #      untracked:  [3] features/shell_aliases.feature
  #      untracked:  [4] mkramdisk.sh
  #

You can now use these numbers in place of filenames when calling normal git
commands, e.g. 'git add 2 3' or 'git checkout 1'.

You can also use numeric ranges, e.g. 'git reset 2-4'. Ranges can even be mixed
with normal numeric operands.

By default, scmpuff will also define a few handy shortcuts to save your fingers,
e.g. 'ga', 'gd', 'gco'.  Check your aliases to see what they are.`)
	},
}
