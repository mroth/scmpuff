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
