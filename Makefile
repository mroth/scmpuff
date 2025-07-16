# build the binary
# note starting in go1.18 vcs information will be available via go version -m
build:
	$(eval VERSION := $(shell git describe --tags HEAD 2>/dev/null || echo unknown))

	@echo "Building as version: $(VERSION)"
	go build -o bin/scmpuff -mod=readonly -ldflags "-X main.version=$(VERSION)"

# run unit tests
test:
	go test ./...

# run integration tests (not including work in progress features)
features: build
	bundle exec cucumber -s --tags='not @wip'

# run integration tests (work in progress features only)
features-wip: build
	bundle exec cucumber -s --tags=@wip

# package as if for distribution
package:
	goreleaser release --clean --skip publish,homebrew

# clean temp files (aruba tmp directory)
clean:
	rm -rf ./tmp

# remove all build artifacts
clobber:
	rm -rf ./dist
	rm -f bin/scmpuff

.PHONY: build install test features features-wip package clean clobber
