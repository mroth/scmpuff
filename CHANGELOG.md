# Change Log
This project tries to adhere to [Semantic Versioning](http://semver.org/).

## Planned for the future
- Switch to `gonut` instead of `godep` for more reliable builds.

## 0.0.2 - Unreleased
### Changed
- Switched to using `status -z` instead of `status --porcelain` for obtaining
  work tree status.  This adds a bit if parsing complexity, but should be the
  absolute most robust long term way to read things, and should enhance cross
  platform support in the future.
- TODO: Use `TAB` as IFS character for file-list instead of `|`. This should
  still be understandable by most shells but significantly less likely to appear
  in a filename.

## 0.0.1 - 2014-03-04
First "ready for daily usage" internal alpha version.
