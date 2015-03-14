# scmpuff :dash:

> Makes working with git from the command line quicker by substituting numeric
shortcuts for files.

<img width=568 src="http://f.cl.ly/items/2726271z170L2y0K3d0b/scmpuff_screenshot.png">

**scmpuff** is a minimalistic implementation of the core functionality of
[SCM Breeze][scmbreeze], without many of the extras.

Its focus is on simplicity, speed, robustness, and cross-platform support. The
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

Once this is semi-feature complete, I will work on getting a package in
Homebrew etc.

[download]: https://github.com/mroth/scmpuff/releases/latest

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

Behind the scenes, scmpuff is assigning filenames to sequential environment
variables, e.g. `$e1`, `$e2`, so you can refer to those with other commands too
if needed.

By default, scmpuff will also define a few handy shortcuts to save your fingers,
e.g. `ga`, `gd`, `gco`.  Check your aliases to see what they are.


## FAQ

### How do you pronounce it?

:information_desk_person: I like to say "scum puff." But I'm weird.

### How does it compare with SCM Breeze?

The short version: it does less, but what it does it should do faster and more
reliably.

The longer version... There are a number of notable features we don't have (on
purpose):

 - Design Asset Management.
 - Keyboard combinations.
 - File indexing/search.
 - ...Pretty much anything beyond numeric filename substitution.

There are also some philosophical differences that underpin the design:

 - **“Do one thing well, and avoid doing anything else.”**

 - **Don't set any environment variables that are unnecessary.** scm_breeze has
   a tendency to pollute the environment with quite a lot of variables, we avoid
   this as much as possible here.

 - **Do as much as possible in cross-platform compiled code, as little as
   possible in shell script.** Shell script can be slow, brittle, and hard to
   maintain, and in particular is problematic for cross-platform utilities due
   to subtle difference between implementations. `scmpuff` contains under 100
   lines of shell script total (compared to ~2000 in scm_breeze), and effort
   is made to keep it as simple as possible.

 - **Rigorously test all interactions with the operating system.** scmpuff has a
   robust (and growing) set of integration tests that aim to verify the behavior
   of the program across different shells and operating systems.


### Can I disable or change the default git shortcut alias names?
You can disable them via passing `--aliases=false` to the `scmpuff init` call
in your shell initialization.  Then, if you wish to remap them, simple modify
your default aliases wherever you normally do, but add aliases mapped to the
scmpuff shell functions, e.g. `alias gs='scmpuff_status'`.

### Okay but really, _why_ did you clone SCM Breeze?
Well, you see, I went drinking with some of the IoJS guys and one thing led to
another and when I woke up I had a bad hangover and noticed my text editor was
open...

Okay in all seriousness, SCM Breeze is a great tool, and has been an
indispensable part of my daily workflow for years.

That said, it does a lot more than I need, and I have a bit of an obsession with
having my core tools being lean and mean. In addition, the cross-platform story
is not great, and functionality on MacOSX/Darwin in particular can be a bit
buggy.  I've attempted to contribute some [patches to SCM Breeze over the
years][patches], but felt like it was time for a rewrite with a different focus
and philosophy.

That said, I don't see scmpuff as a competitor for scm_breeze, but rather as a complimentary tool for people who may prefer a "lite" version.

[patches]: https://github.com/ndbroadbent/scm_breeze/issues?q=author%3Amroth


## Benchmarks

In normal usage I've always found scm_breeze to be acceptably fast, but scmpuff
should be roughly an order of magnitude faster, so here are some informal
benchmarks, measured on my 2011 MacBook Air.

Time to init during shell startup:

    scm_breeze: 0.09sec
    scmpuff:   <0.01sec

Full color numbered status on a complex git repo with ~500 work tree changes:

    scm_breeze: 1.14sec*
    scmpuff:    0.13sec

_Note: scm_breeze normally falls back to normal git status after a configurable
`$gs_max_changes=150`, which I modified here for testing. For comparison, I get
about `0.08sec` for a plain `git status`._


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
