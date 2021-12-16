## Development

While the build process itself does not require it, development uses Ruby for
integration testing because of the excellent Cucumber/Aruba package for testing
CLI tools.

Thus, to bootstrap, you will need to have Ruby and bundler installed on your
system.  Do `bundle install` to get the dev environment going.

Since we already have Ruby then for tests, we use a Rakefile instead of Makefile
since it offers some niceties.  Do `rake -T` to see available tasks.
