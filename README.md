# scmpuff :dash:

![miley-puffs](http://media.giphy.com/media/nF8Sgd4X74be/giphy.gif)

Makes working with git from the command line quicker by substituting numeric
shortcuts for files.

[SCREENSHOT]

**scmpuff** is a minimalistic implementation of the core functionality of
[scm_breeze][scmbreeze], without many of the extras.  

It's focus is on simplicity, speed, and cross-platform support. The majority of
the functionality is contained within a compiled binary, and the shell
integration is under 100 lines of shell script.

scmpuff currently functions in `bash` and `zsh` in any *nix-like operating
system. It's written with cross-platform support in mind, so hopefully we'll
have it functioning on Windows soon as well.

[scmbreeze]: https://github.com/ndbroadbent/scm_breeze

## Installation

For now, you have to do it manually.  If you have a Go dev setup, you can do:

    go get github.com/mroth/scmpuff

If you downloaded a binary, place it in `/usr/local/bin` or some other place
on your system's default `PATH`.

Once this is semi-feature complete, I will work on binary cross-compilation and
putting a package in Homebrew.


## Setup

Currently scmpuff supports bash and zsh for all functionality.

To initialize shell functions, add the following to your `~/.bash_profile` or
`~/.zshrc` file:

    eval "$(scmpuff init -s)"

This will define the scmpuff shell functions as well as some git shortcuts.


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

By default, scmpuff will also define a few handy shortcuts to save your fingers,
e.g. `ga`, `gd`, `gco`.  Check your aliases to see what they are.


## FAQ

### How do you pronounce it?

:information_desk_person: I like to say "scum puff." But I'm weird.

### How does it compare with scm_breeze?

There are a number of notable features we don't have (on purpose):

 - Design Asset Management.
 - Keyboard "shortcuts".
 - File indexing/search.
 - Pretty much anything beyond numeric filename substitution.

There are also some philosophical differences that underpin the design:

 - *"Do one thing well, and avoid doing anything else."*

 - *Don't set any environment variables that are unnecessary.* scm_breeze has a
   tendency to pollute the environment with quite a lot of variables, we avoid
   this as much as possible here.

 - *Do as much as possible in cross-platform compiled code, as little as
   possible in shell script.* Shell script can be slow, brittle, and hard to
   maintain, and in particular is problematic for cross-platform utilities due
   to subtle difference between implementations. `scmpuff` contains under 100
   lines of shell script total (compared to ~2000 in scm_breeze), and effort
   is made to keep it as simple as possible.

 - _Rigorously test all interactions with the operating system._ scmpuff has a
   robust (and growing) set of integration tests that aim to verify the behavior
   of the program across different shells and operating systems.


### Can I disable or change the default git shortcut alias names?
You can disable them via passing `--aliases=false` to the `scmpuff init` call
in your shell initialization.  Then, if you wish to remap them, simple modify
your default aliases wherever you normally do, but add aliases mapped to the
scmpuff shell functions, e.g. `alias gs='scmpuff_status'`.

### Okay but really, _why_ did you fork scm_breeze?
Well, you see, I went drinking with some of the IoJS guys and one thing led to
another and when I woke up I had a bad hangover and noticed my text editor was
open...

Okay in all seriousness, scm_breeze is a great tool, and has been an
indispensable part of my daily workflow for years.

That said, it does a lot more than I need, and I have a bit of an obsession with
having my core tools being lean and mean. In addition, the cross-platform
story is not great, and functionality on MacOSX/Darwin can be a bit buggy.  I've
attempted to contribute some patches back to scm_breeze over the years, [with
limited success][patches].

That said, I don't see scmpuff as a replacement or competitor for scm_breeze,
but rather as a complimentary tool for people who may prefer a "lite" version.
So in my mind, it's really not a fork!

[patches]: https://github.com/ndbroadbent/scm_breeze/issues?q=author%3Amroth

## Benchmarks
Some informal benchmarks, measured on Macbook Air.

Time to evaluate shell setup:

Status on a git repo with ~100 changes:

## Development

While the build process itself does not require it, development uses Ruby for
integration testing because of the excellent Cucumber/Aruba package for testing
CLI tools.

Thus, to bootstrap, you will need to have Ruby and bundler installed on your
system.  Do `bundle install; rake bootstrap` to get the dev environment going.
We assume you are both cloned into and have your $GOPATH properly set.

Since we already have Ruby then for tests, we use a Rakefile instead of Makefile
since it offers some niceties.  Do `rake -T` to see available tasks.

`GO_VERSION >= 1.4` is required to build.
