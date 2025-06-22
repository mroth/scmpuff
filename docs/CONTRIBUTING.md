## Development

While the build process itself does not require it, development currently uses
Ruby for integration testing because of the excellent Cucumber/Aruba package for
testing CLI tools.

Thus, to run integration tests, you will need to have Ruby and bundler installed
on your system.  Do `bundle install` to get the test environment going.

If you don't have a local Ruby environment, there is a VSCode .devcontainer
provided to make things simpler.