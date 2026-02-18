# Testing

Quick reference for make targets:

| Target             | Command                                   | What it runs                   |
|--------------------|-------------------------------------------|--------------------------------|
| `make test`        | `go test -short ./...`                    | Unit tests only                |
| `make integration` | `go test ./internal/cmd -run TestScripts` | Integration tests (testscript) |
| `make lint`        | `golangci-lint run`                       | Linter                         |

## Test Types

### Unit tests

Unit tests cover internal logic and rendering/parsing behavior.

- Framework: Go standard `testing` package.
- Typical style: table-driven tests in `*_test.go` files.
- CI command: `go test -short ./...`.
- Make target: `make test`.

`-short` is the contract for "fast/default" tests and is what CI runs in the unit-test workflow.

### Integration tests (testscript)

Integration tests validate CLI behavior using script fixtures.

- Framework: `github.com/rogpeppe/go-internal/testscript`.
- Primary entrypoint: `internal/cmd/testscript_test.go` (`TestScripts`).
- Script fixtures: `internal/cmd/testdata/script/*.txtar`.
- Make target: `make integration`.

These tests are skipped in short mode (`testing.Short()`), so they do not run in default CI unit test jobs.

## When to Add or Update Tests

Update tests whenever behavior changes in one of these areas:

- Argument expansion and environment evaluation.
- CLI command flags, output, or error behavior.
- Git status parsing/conversion logic.
- Status renderer formatting or ordering.
- Shell init script generation.

General rule: code changes should include test updates in the nearest relevant package.

## Golden Files: When and How to Update

Use golden files for stable, user-visible renderer output.

- If renderer output intentionally changes, update goldens.
- For status renderer tests:
  - Update mode: run tests with `-update`.
  - Optional cleanup: `-clobber` to remove old goldens (see test flags in `render_test.go`).

Example:

```sh
go test ./internal/cmd/status -run TestRenderer_Display -update
```

Only update goldens for intentional output changes. If updates are unexpected, inspect the renderer change first.

## Adding Parsing Regression Tests

`scmpuff debug dump --archive` produces a diagnostic archive containing raw porcelain output in all formats. When a user reports a parsing failure, the raw output from their archive can be added to `internal/gitstatus/porcelainv1/testdata/` as a regression test fixture — similar to the golden file workflow but for the parsing layer.

## Integration Test Style Guide

Use testscript conventions for files under `internal/cmd/testdata/script/*.txtar`.

A nice series of posts about using testscript:

- https://bitfieldconsulting.com/posts/test-scripts
- https://bitfieldconsulting.com/posts/cli-testing
- https://bitfieldconsulting.com/posts/test-scripts-files
- https://bitfieldconsulting.com/posts/conditions-concurrency

When debugging script failures locally, run with `-testwork` to preserve the script work directory:

```sh
go test ./internal/cmd -run TestScripts -count=1 -testwork
```

### testscript setup

The test harness in `internal/cmd/testscript_test.go` registers a `scmpuff` binary command for use in test scripts. The setup function isolates tests from the host git configuration and configures deterministic author/committer identities so that output is reproducible across machines.

### Multi-shell testing pattern

Integration scripts test bash, zsh, and fish using `[exec:shell]` conditions on each line. This means every shell block is repeated almost verbatim — an unfortunate side effect of testscript being a linear DSL with no loops or templating. Attempts to abstract the repetition (custom commands, Go-level shell loops) add indirection that hurts debuggability, so given the relatively small amount of shell-specific testing we need to do, the duplication is an acceptable tradeoff. As a design goal, scmpuff behavior is shell-independent — the core logic lives in Go, and shell-specific code is limited to the thin init/wrapper layer. Keep the test blocks explicit.

**Init boilerplate.** Bash and zsh use the same form:

```
[exec:bash] exec bash -c 'eval "$(scmpuff init -s)"; …'
[exec:zsh]  exec zsh  -c 'eval "$(scmpuff init -s)"; …'
```

Fish uses pipe-to-source:

```
[exec:fish] exec fish -c 'scmpuff init --shell=fish | source; …'
```

**When shells diverge.** Fish has different syntax for exit status (`$status` vs `$?`), variable assignment (`set` vs `=`), and aliases. Keep these as explicit per-shell blocks rather than trying to unify them:

```
[exec:bash] exec bash -c '…; test $? -eq 128'
[exec:zsh]  exec zsh  -c '…; test $? -eq 128'
[exec:fish] exec fish -c '…; test $status -eq 128'
```

**Maintenance rule.** When updating a shell block, update all three. Bash and zsh are usually identical (just swap the shell name in the condition and `exec`); fish needs its own variant for init and any shell-specific syntax.

**State isolation.** Tests that mutate repo state (add, commit, etc.) need per-shell copies of the repo to avoid cross-shell interference:

```
[exec:bash] exec cp -R repo_base repo_bash
[exec:bash] cd repo_bash
…
[exec:zsh]  exec cp -R repo_base repo_zsh
[exec:zsh]  cd repo_zsh
…
[exec:fish] exec cp -R repo_base repo_fish
[exec:fish] cd repo_fish
```

Read-only tests that only inspect output can share a single repo.
