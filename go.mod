module github.com/mroth/scmpuff

go 1.24

retract v0.6.1 // init -s writes to stderr instead of stdout, breaking eval

require (
	github.com/caarlos0/go-version v0.2.2
	github.com/google/go-cmp v0.7.0
	github.com/mroth/porcelain v0.1.0
	github.com/rogpeppe/go-internal v1.14.1
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
)
