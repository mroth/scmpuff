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


## Installation

For now, you have to do it from source:

    go install github.com/mroth/scmpuff

Once this is semi feature complete, will work on binary cross-compilation and
putting a package in Homebrew.

## FAQ

### Wait it doesn't do anything?

Yeah, this is a work in progress, still in development.

### How do you pronounce it?
I like to say "scum puff."

### How does it compare with scm_breeze?

Features we don't have
 * Design Asset Management
 * Most everything except var substitution

Antipatterns
 - Don't set any environment variables that are unnecessary.
