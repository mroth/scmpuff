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


## TODO

Known issues:

 - doesn't handle getting additional args and passing them along, see test in scm_breeze status_shortcuts_test.sh line 58

 - switch to buffered output writer for speed! (issue #1)

 - get rid of the retarded way scmbreeze hardcodes padding into the msg for things,can be better handled with printf padding

 - maybe use chans to make multiple calls to git CLI concurrent? (e.g. branch, status, second status if needed...) would need to bench and see if it actually helps, could maybe even hurt?

 - maybe dont use first line of output trick and instead use STDERR? 
