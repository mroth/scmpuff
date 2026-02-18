# Architecture

scmpuff is a Go CLI with a thin shell integration layer. The Go binary handles parsing git status output, rendering numbered status display, and expanding numeric shortcut arguments. The shell layer (bash/zsh/fish scripts, emitted at init time) exports environment variables and intercepts git commands to wire everything together.

## Directory structure

```
main.go                          Entry point, version info injection, banner embed

internal/
├── arguments/                   Numeric shortcut expansion (1 → $e1, 1-3 → $e1 $e2 $e3)
│
├── cmd/
│   ├── debug/                   `scmpuff debug dump` — diagnostic archive
│   ├── exec/                    `scmpuff exec` — run commands with shortcut expansion
│   ├── expand/                  `scmpuff expand` — expand shortcuts to paths (scripting/debug)
│   ├── inits/                   `scmpuff init` — shell initialization script generation
│   │   └── data/                Embedded shell scripts (bash/zsh/fish)
│   ├── intro/                   `scmpuff intro` — help/getting-started command
│   └── status/                  `scmpuff status` — parsing, rendering, numbering
│
└── gitstatus/
    ├── gitstatus.go             Data structures: StatusInfo, BranchInfo, StatusItem, enums
    └── porcelainv1/             Porcelain v1 parser (raw git output → structured data)
```

## End-to-end flows

### 1. Shell initialization

**Trigger:** `eval "$(scmpuff init -s)"` in user's shell profile.

`scmpuff init` detects the user's shell (from `--shell` flag or `$SHELL`) and emits a script to stdout that the shell evaluates. The script installs three things:

1. **`scmpuff_status()` function** — wraps `scmpuff status --filelist`, captures the machine-readable file list, and exports `$e1`..`$eN` environment variables for each file.
2. **`git()` wrapper function** — intercepts git subcommands and routes them through `scmpuff exec` for numeric shortcut expansion (see [shell-integration.md](shell-integration.md) for the dispatch table).
3. **Short aliases** — `gs`, `ga`, `gd`, `gl`, `gco`, `grs` for common operations.

The wrapper and aliases are each controlled by flags (`--wrap`, `--aliases`, both default on). Shell scripts are embedded in the binary at compile time via `go:embed`.

### 2. Status display

**Trigger:** User types `gs` (alias for `scmpuff_status`).

1. Shell alias `gs` calls `scmpuff_status()` shell function.
2. `scmpuff_status()` shell function runs `scmpuff status --filelist` and captures output.
3. In the scmpuff binary, `scmpuff status` runs `git status` and parses the porcelain output (see [git-status-parsing.md](git-status-parsing.md) for the full pipeline), renders it into a combination output of metadata and display info (see [Status rendering](#status-rendering) below).
4. Back in the shell, `scmpuff_status()` extracts the first line of the output which contains the metadata (tab-delimited file list), parses it, and exports `$e1`, `$e2`, ... `$eN` to the shell as environment variables. Lines 2+ (the colorized display) are printed to the terminal.

### 3. Numeric shortcut expansion

**Trigger:** User types `git add 1 2` or `git diff 1-3`.

1. The shell `git()` wrapper function intercepts the command and, based on the subcommand, dispatches to `scmpuff exec` (see [shell-integration.md](shell-integration.md) for the full dispatch table).
2. `scmpuff exec` expands numeric arguments to environment variable references (`1` → `$e1`, `1-3` → `$e1 $e2 $e3`), then resolves each `$eN` to the actual file path it was set to during the last status display. See [Argument expansion](#argument-expansion) below for details.
3. The fully resolved argument list is used to exec the underlying git command as a subprocess.

## Argument expansion

The `internal/arguments` package handles converting numeric shortcuts into file paths. The pipeline has two stages:

1. **Symbolic expansion** — Numeric tokens become environment variable references: `3` → `$e3`, `1-3` → `$e1 $e2 $e3`. If a file literally named `3` exists on disk, the number is left as-is. Non-numeric arguments pass through unchanged.

2. **Environment resolution** — Each `$eN` reference is resolved to the absolute file path stored during the last status display. For commands that need relative paths (like `git diff`), the absolute path is converted to a path relative to the current working directory.

## Status rendering

After parsing (see [git-status-parsing.md](git-status-parsing.md)), the status renderer produces the display output:

1. **Grouping**: Items are bucketed by `StatusGroup` (derived from each item's `ChangeType`).
2. **Display order**: Groups render in fixed order — Staged → Unmerged → Unstaged → Untracked.
3. **Sequential numbering**: Items are numbered `[1]`, `[2]`, ... sequentially across all groups.
4. **Color mapping**: Each `StatusGroup` has a group color (for the `#` gutter and file path) and each `ChangeState` has a state color (for the change message like "modified"). See `color.go` for the mappings.
5. **Machine-parseable output** (`--filelist`): A tab-delimited line of absolute paths in display order, consumed by the shell function to set `$e1`..`$eN`.

## External dependencies

| Dependency | Import path                       | Purpose                                                       |
|------------|-----------------------------------|---------------------------------------------------------------|
| cobra      | `github.com/spf13/cobra`          | CLI framework                                                 |
| porcelain  | `github.com/mroth/porcelain`      | Low-level git porcelain parser                                |
| go-version | `github.com/caarlos0/go-version`  | Structured version info display                               |
| go-cmp     | `github.com/google/go-cmp`        | Structured comparison in tests                                |
| testscript | `github.com/rogpeppe/go-internal` | Integration test framework (txtar scripts)                    |
