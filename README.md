# scmpuff :dash:

> Makes working with git from the command line quicker by substituting numeric
shortcuts for files.

[![Build Status](https://travis-ci.org/mroth/scmpuff.svg?branch=master)](https://travis-ci.org/mroth/scmpuff)

<img width=568 src="http://f.cl.ly/items/2726271z170L2y0K3d0b/scmpuff_screenshot.png">

**scmpuff** is a minimalistic reinterpretation of the core functionality of
[SCM Breeze][scmbreeze], without many of the extras.

It is focused on simplicity, speed, robustness, and cross-platform support. The
majority of the functionality is contained within a compiled binary, and the
shell integration is under 100 lines of shell script.

**scmpuff** currently functions in `bash` and `zsh` in any *nix-like operating
system. It's written with cross-platform support in mind, so hopefully we'll
have it functioning on Windows soon as well.

**scmpuff** is fully compatible with the most-excellent [Hub][hub].

[scmbreeze]: https://github.com/ndbroadbent/scm_breeze
[hub]: https://github.com/github/hub

## Installation

[Download] the binary for your platform, and copy it to `/usr/local/bin` or
somewhere else in your  default `$PATH`.

Alternately, if you use [homebrew], you can just: `brew install scmpuff`. :beer:

[download]: https://github.com/mroth/scmpuff/releases/latest
[homebrew]: http://brew.sh

## Setup

Currently scmpuff supports bash and zsh for all functionality.

To initialize shell functions, add the following to your `~/.bash_profile` or
`~/.zshrc` file:

    eval "$(scmpuff init -s)"

This will define the scmpuff shell functions as well as some handy shortcuts.


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

### How do you pronounce it?

:information_desk_person: I like to say "scum puff." But I'm weird.

### How does it compare with SCM Breeze?

The short version: it does less, but is faster and should be significantly more
stable and reliable, especially across different platforms.

The long, detailed version:
https://github.com/mroth/scmpuff/wiki/scmpuff-vs-SCM-Breeze

### Can I disable or change the default git shortcut aliases?
You can disable them via passing `--aliases=false` to the `scmpuff init` call
in your shell initialization.  Then, if you wish to remap them, simple modify
your default aliases wherever you normally do, but add aliases mapped to the
scmpuff shell functions, e.g. `alias gs='scmpuff_status'`.


## Development

While the build process itself does not require it, development uses Ruby for
integration testing because of the excellent Cucumber/Aruba package for testing
CLI tools.

Thus, to bootstrap, you will need to have Ruby and bundler installed on your
system.  Do `bundle install; rake bootstrap` to get the dev environment going.
We assume you are both cloned into and have your $GOPATH properly set.

Since we already have Ruby then for tests, we use a Rakefile instead of Makefile
since it offers some niceties.  Do `rake -T` to see available tasks.

`GO_VERSION >= 1.6` is required to build (or 1.5 with `GO15VENDOREXPERIMENT=1`).
