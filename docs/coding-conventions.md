# Coding Conventions

## Formatting and linting

All Go code must be formatted with `gofmt`. Run `make lint` before pushing —
the project uses [golangci-lint](https://golangci-lint.run/) v2 with the
config in `.golangci.yml`.

## Error handling

Wrap errors with context using `fmt.Errorf("…: %w", err)`. Use `errors.Is()`
and `errors.As()` for type checks.

## Naming

Follow standard Go naming conventions. However, favor descriptive names for
longer-lived variables and parameters -- git internals are often complex and
require clear naming to disambiguate their purpose.

## Comments

All exported symbols should have doc comments. Use inline `// NOTE:` comments
to flag non-obvious caveats.

Don't comment obvious code, but do describe with sufficient detail the purpose
of the code when helpful to aid understanding and maintainability, in particular
to assist with onboarding new contributors.

## Switch exhaustiveness

The `exhaustive` linter is enabled for both `switch` and `map` checks. When
switching on a typed enum, handle every case explicitly rather than relying on
`default`. For truly unreachable cases, `panic()` is preferred over a silent
default.

## Tests

- Name tests after the function they cover (e.g. `TestFunctionName` or
  `Test_functionName`). Use sub-tests to nest scope.
- Prefer table-driven tests over multiple test functions. Include only the
  minimum cases needed — each should test a unique aspect, not repeat coverage.
- Never compare against an error's string value. Use `errors.Is` or
  `errors.As`, or define a sentinel error variable/type if one doesn't exist.
- Mark `t.Helper()` on any test helper that can fail or panic.
- Do not use `testify` or other assertion libraries. Use standard `if` checks
  with `t.Error`/`t.Fatal`. For struct comparison output, `go-cmp` is
  acceptable.
- See also [docs/testing.md](testing.md) for project-specific test scaffolding.

## Dependencies

Prefer the standard library over external dependencies.

## Modern idioms

Prefer modern, idiomatic Go code that takes advantage of features available in
the most recent Go version defined by the project's `go.mod` file.
