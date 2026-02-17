## Development

If you prefer a containerized environment, use the provided VS Code
`.devcontainer`, which is configured for this workflow.

## Testing

### Test Types

#### Unit tests

Unit tests cover internal logic and rendering/parsing behavior.

- Framework: Go standard `testing` package.
- Typical style: table-driven tests in `*_test.go` files.
- CI command: `go test -short ./...`.
- Make target: `make test`.

`-short` is the contract for “fast/default” tests and is what CI runs in the unit-test workflow.

#### Integration tests (testscript)

Integration tests validate CLI behavior using script fixtures.

- Framework: `github.com/rogpeppe/go-internal/testscript`.
- Primary entrypoint: `internal/cmd/testscript_test.go` (`TestScripts`).
- Script fixtures: `internal/cmd/testdata/script/*.txtar`.
- Make target: `make integration`.

These tests are skipped in short mode (`testing.Short()`), so they do not run in default CI unit test jobs.

### Frameworks and Utilities

- `testing`: core unit tests and subtests.
- `testscript`: integration scenarios with command/script fixtures.
- `go-cmp` (`github.com/google/go-cmp/cmp`): use for structured value comparisons in tests where diffs are clearer than `reflect.DeepEqual` failures.
- Golden-file comparisons: used for status rendering outputs.
  - Example: `internal/cmd/status/render_test.go`
  - Goldens: `internal/cmd/status/testdata/*.golden`
- Policy: do not import additional testing frameworks (for example `testify`); keep tests on `testing`, `testscript`, and approved utilities like `go-cmp`.

### Running Tests

Fast/local unit run:

```sh
go test -short ./...
```

Verbose unit run:

```sh
go test -short -v ./...
```

Integration run:

```sh
go test ./internal/cmd -run TestScripts
```

Or via Make:

```sh
make test
make integration
```

### When to Add or Update Tests

Update tests whenever behavior changes in one of these areas:

- Argument expansion and environment evaluation.
- CLI command flags, output, or error behavior.
- Git status parsing/conversion logic.
- Status renderer formatting or ordering.
- Shell init script generation.

General rule: code changes should include test updates in the nearest relevant package.

### Golden Files: When and How to Update

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

### Practical Guidance

- Prefer deterministic tests (avoid time/network/external dependencies).
- Keep unit tests short and isolated.
- Use integration `testscript` when validating command interactions and shell-facing behavior.
- Keep fixtures focused and readable; split scenarios instead of creating giant cases.

### Testscript Style Guide

Use testscript conventions for files under `internal/cmd/testdata/script/*.txtar`.

- https://bitfieldconsulting.com/posts/test-scripts
- https://bitfieldconsulting.com/posts/cli-testing
- https://bitfieldconsulting.com/posts/test-scripts-files
- https://bitfieldconsulting.com/posts/conditions-concurrency

When debugging script failures locally, run with `-testwork` to preserve the script work directory:

```sh
go test ./internal/cmd -run TestScripts -count=1 -testwork
```

#### Multi-shell testing pattern

Integration scripts test bash, zsh, and fish using `[exec:shell]` conditions on
each line. This means every shell block is repeated almost verbatim — an
unfortunate side effect of testscript being a linear DSL with no loops or
templating. Attempts to abstract the repetition (custom commands, Go-level shell
loops) add indirection that hurts debuggability, so given the relatively small
amount of shell-specific testing we need to do, the duplication is an acceptable
tradeoff. As a design goal, scmpuff behavior is shell-independent — the core
logic lives in Go, and shell-specific code is limited to the thin init/wrapper
layer. Keep the test blocks explicit.

**Init boilerplate.** Bash and zsh use the same form:

```
[exec:bash] exec bash -c 'eval "$(scmpuff init -s)"; …'
[exec:zsh]  exec zsh  -c 'eval "$(scmpuff init -s)"; …'
```

Fish uses pipe-to-source:

```
[exec:fish] exec fish -c 'scmpuff init --shell=fish | source; …'
```

**When shells diverge.** Fish has different syntax for exit status (`$status`
vs `$?`), variable assignment (`set` vs `=`), and aliases. Keep these as
explicit per-shell blocks rather than trying to unify them:

```
[exec:bash] exec bash -c '…; test $? -eq 128'
[exec:zsh]  exec zsh  -c '…; test $? -eq 128'
[exec:fish] exec fish -c '…; test $status -eq 128'
```

**Maintenance rule.** When updating a shell block, update all three. Bash and
zsh are usually identical (just swap the shell name in the condition and `exec`);
fish needs its own variant for init and any shell-specific syntax.

**State isolation.** Tests that mutate repo state (add, commit, etc.) need
per-shell copies of the repo to avoid cross-shell interference:

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
