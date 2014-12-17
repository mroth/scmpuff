These scripts are built into the binary by the `go-bindata` program during
compilation.  Thus be sure to run `go generate` for the package directory or
compilation will fail.

Note, go-bindata does not appear to be smart enough to automatically check for
changes and regenerate the files, so if you are making changes to these scripts
you may need to clean the bindata prior to recompilation. (The rake build
process should handle this automatically.)
