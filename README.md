# scmpuff :dash:

![miley-puffs](http://media.giphy.com/media/nF8Sgd4X74be/giphy.gif)

Makes working with git from the command line quicker by substituting numeric
shortcuts for files.

[SCREENSHOT]

**scmpuff** is a minimalistic implementation of the core functionality of
[scm_breeze][scmbreeze], without many of the extras.  For true overkill, it's
written in Go and compiled to native binary application so you'll never notice
any speed issues.

[scmbreeze]: https://github.com/ndbroadbent/scm_breeze


## SPEED
(measured on Macbook Air)

Time to evaluate shell setup:

Status on a git repo with ~100 changes:

## Installation

For now, you have to do it from source:

    go get github.com/mroth/scmpuff

(Once this is semi feature complete, will work on binary cross-compilation and
putting a package in Homebrew.)


## Setup
Currently scmpuff supports bash and zsh for all functionality.

To initialize shell functions, add the following to your `~/.bash_profile` or
`~/.zshrc` file:

    eval "$(scmpuff init -s)"

This will define the scmpuff shell functions as well as some git shortcuts.

## Usage


## FAQ

### Wait it doesn't do anything?

Yeah, this is a work in progress, still in development.

### How do you pronounce it?
I like to say "scum puff."

### How does it compare with scm_breeze?

Features we don't have:
 * Design Asset Management
 * Key chording
 * File indexing/search
 * Most everything except var substitution

Antipatterns:
 - Don't set any environment variables that are unnecessary.
 - Do as much as possible in cross-platform compiled code, as little as possible
   in shell script.

### Can I disable or change the default git shortcut alias names?
You can disable them via passing `--aliases=false` to the `scmpuff init` call
in your shell initialization.  Then, if you wish to remap them, simple modify
your default aliases wherever you normally do, but add aliases mapped to the
scmpuff shell functions, e.g. `alias gs='scmpuff_status_shortcuts'`.


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


## TODO

Known issues:

 - doesn't handle getting additional args and passing them along, see test in scm_breeze status_shortcuts_test.sh line 58

 - get rid of the retarded way scmbreeze hardcodes padding into the msg for things,can be better handled with printf padding

 - maybe use chans to make multiple calls to git CLI concurrent? (e.g. branch, status, second status if needed...) would need to bench and see if it actually helps, could maybe even hurt?

 - maybe dont use first line of output trick and instead use STDERR?
