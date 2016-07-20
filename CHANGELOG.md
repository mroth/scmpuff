# Change Log
This project tries to adhere to [Semantic Versioning](http://semver.org/).

## 0.2.1 - In Development
### Fixed
- Fix for expanding filenames containing a `|` character. (#21)


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
