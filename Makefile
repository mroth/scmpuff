# build the binary
# note starting in go1.18 vcs information will be available via go version -m
build:
	$(eval VERSION := $(shell git describe --tags HEAD 2>/dev/null || echo unknown))
	@echo "Building as version: $(VERSION)"
	go build -o bin/scmpuff -mod=readonly -ldflags "-X main.version=$(VERSION)"

# run unit tests
test:
	go test -short ./...

# run lint tests
lint:
	golangci-lint run

# run integration tests (testscript harness)
integration:
	go test ./internal/cmd -run TestScripts

# package as if for distribution
package:
	goreleaser release --clean --skip publish,homebrew

# clean temp files
clean:
	@:

# remove all build artifacts
clobber:
	rm -rf ./dist
	rm -f bin/scmpuff

.PHONY: build install test integration lint package clean clobber
