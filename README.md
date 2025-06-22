# scmpuff 🔢

> Makes working with git from the command line quicker by substituting numeric
shortcuts for files.

<img width=568 src="http://f.cl.ly/items/2726271z170L2y0K3d0b/scmpuff_screenshot.png">

**scmpuff** is a minimalistic reinterpretation of the core functionality of
[SCM Breeze][scmbreeze].

It is focused on simplicity, speed, robustness, and cross-platform support. The
majority of the functionality is contained within a compiled binary, and the
shell integration is under 100 lines of shell script.

**scmpuff** currently has built-in support for `bash`, `zsh`, and `fish`.

[scmbreeze]: https://github.com/ndbroadbent/scm_breeze


## Installation

[Download] the binary for your platform, and copy it to `/usr/local/bin` or
somewhere else in your  default path.

Alternately, if you use [Homebrew], you can just: `brew install scmpuff`.

[download]: https://github.com/mroth/scmpuff/releases/latest
[Homebrew]: https://brew.sh


## Setup

Currently scmpuff supports bash, zsh and fish for all functionality.

To initialize shell functions, add the following to your `~/.bash_profile` or
`~/.zshrc` file:

    eval "$(scmpuff init -s)"

For [fish] shell, add the following to your `~/.config/fish/config.fish` file:

    scmpuff init --shell=fish | source

This will define the scmpuff shell functions as well as some handy shortcuts.

[fish]: https://fishshell.com/


## Usage

**Once things are loaded, the most important function you will want to know
about is `scmpuff_status`, which is aliased to `gs` for short.**

This is a replacement for `git status` that is pretty and shows you numbers next
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

**You can now use these numbers in place of filenames when calling normal git
commands, e.g. `git add 2 3` or `git checkout 1`.**

You can also use numeric ranges, e.g. `git reset 2-4`. Ranges can even be mixed
with normal numeric operands.

Behind the scenes, scmpuff is assigning filenames to sequential environment
variables, e.g. `$e1`, `$e2`, so you can refer to those with other commands too
if needed.

By default, scmpuff will also define a few handy shortcuts to save your fingers,
e.g. `ga`, `gd`, `gco`.  Check your aliases to see what they are.


## FAQ

### How does it compare with SCM Breeze?

The short version: it does less, but is faster and should be more stable and
reliable, especially across different platforms.

The long, detailed version:
https://github.com/mroth/scmpuff/wiki/scmpuff-vs-SCM-Breeze

### Can I disable or change the default git shortcut aliases?
You can disable them via passing `--aliases=false` to the `scmpuff init` call
in your shell initialization.  Then, if you wish to remap them, simple modify
your default aliases wherever you normally do, but add aliases mapped to the
scmpuff shell functions, e.g. `alias gs='scmpuff_status'`.

### I want to use scmpuff in conjunction with [hub][hub] or something else that I've aliased git to, how would I do so?

By default, scmpuff will attempt to utilize the absolute path of whatever `git`
it finds in your system PATH, ignoring existing shell aliases.  If you want to
use a different binary, set `$SCMPUFF_GIT_CMD` in your shell to the path, for
example, `export SCMPUFF_GIT_CMD=/usr/local/bin/hub`.

[hub]: https://github.com/github/hub

