# Contributing to scmpuff

## Development environment

A VS Code `.devcontainer` is provided for a containerized setup. Otherwise, you need:

- **Go+**
- **git**
- Optionally: **bash**, **zsh**, and/or **fish** (for running integration tests locally)

## Building

```sh
make build
```

This produces `bin/scmpuff` with version info injected from git tags.

## Understanding the codebase

- [architecture.md](architecture.md) — overall structure, end-to-end flows, package reference, dependencies
- [git-status-parsing.md](git-status-parsing.md) — parsing pipeline, how XY short codes are handled
- [shell-integration.md](shell-integration.md) — shell scripts, env var protocol, git wrapper dispatch

## Testing

See [testing.md](testing.md) for full details on:

- Unit tests (`make test`) and integration tests (`make integration`)
- Golden file workflow
- Testscript style guide and multi-shell testing patterns

## Code style

See [coding-conventions.md](coding-conventions.md).

## Pull request process

- Include tests for behavior changes (see [testing.md](testing.md) for guidance on when/how).
- If your change affects architecture, parsing, shell integration, or testing workflows, update the corresponding `docs/` files to stay in sync.
- **AI-assisted contributions** must be disclosed. Commits containing AI-generated or AI-assisted code must include a `Co-authored-by` trailer identifying the model used, for example:
  ```
  Co-authored-by: Claude Sonnet 4.5 <noreply@anthropic.com>
  ```
  PRs containing AI-assisted code must include a section from the human author describing what level of human review was applied and how the AI-generated portions were verified.
