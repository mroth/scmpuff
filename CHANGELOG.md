# Change Log
This project tries to adhere to [Semantic Versioning](https://semver.org/).

## 0.3.0 - 2019-05-24
Small fixes and no real feature changes, but enough build tooling was
changed that this is being considered a minor release instead of patch.

### Changed
- Build tooling updated to more modern Go toolchain, removing some legacy/deprecated tools so that this actually can be maintained in 2019.
    - Switch from Godep to go modules (and removed vendored modules).
    - Switch release build cross-compile tool from goxc to goreleaser.
    - go-bindata updates.
- CLI library updated to recent version of Cobra.

### Fixed
- Fix processing statuses with very complex changesets (#45)
- scmpuff expand should escape '*' (#44, thx @jdelStrother)
- scmpuff status works on naked zero commit repo (#37, thx @zommerfelds)

## 0.2.1 - 2017-02-17
### Fixed
- Fix for expanding filenames containing a `|` character. (#21)
- Allow scmpuff_status to handle detached git states. (#24, thx @bikezilla)
- Rebuilt binaries with more recent Go to avoid macOS errors. (fixes #23)

## 0.2.0 - 2016-02-27
### Changed
- Updated build toolchain for compiling in Go 1.6 and beyond.

### Fixed
- Fix for truncated branch names containing periods (issue #12)
- Allow for semicolons in commit messages (PR #7, thx @creature)


## 0.1.1 - 2015-08-02
### Fixed
- Allow for passing along empty arguments during expansion.

## 0.1.0 - 2015-04-22
Initial public release.

### Changed
- Cleaned up documentation and website.
- Continuous integration testing via travis-ci.

### Fixed
- Fixed optional flag passing for `git add` wrapper.


## 0.0.3 - 2015-03-16
### Changed
Two build-chain changes that should make it possible for the project to be
compiled by end-users via just `go get` (making the build script only required
for developers):
- Switched to using `nut` for dependency management, which overwrites import
  paths instead of modifying `$GOPATH`.
- Vendor bindata generation.

### Fixed
- Reset ANSI colors properly after "Not a git repository" error.


## 0.0.2 - 2015-03-10
### Changed
Some preliminary work towards robust cross platform support:
- Switched to using `status -z` instead of `status --porcelain` for obtaining
  work tree status.  This adds a bit if parsing complexity, but should be the
  absolute most robust long term way to read things, and should enhance cross
  platform support in the future.
- Use `TAB` as IFS character for file-list instead of `|`. This should still be
  understandable by most shells but significantly less likely to appear in a
  filename.


## 0.0.1 - 2015-03-04
First "ready for daily usage" internal alpha version.
