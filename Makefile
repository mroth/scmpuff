# set version in main.go for `go install` and build binary
build:
	$(eval VERSION := $(shell git describe --tags HEAD 2>/dev/null || echo unknown))
	$(shell sed -E -i '' -e "s/^var version = \"[A-Za-z0-9\.-_]+\"/var version = \"$(VERSION)\"/" main.go)
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
