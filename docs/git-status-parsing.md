# Git Status Parsing

scmpuff parses `git status` output into structured data for numbered file display. This document covers the parsing pipeline from raw git output to the `StatusInfo` data structure, and the design of the `ChangeType` abstraction that bridges parsing and rendering. For how that data is rendered to the screen, see the "Status rendering" section in [architecture.md](architecture.md).

## Raw format

scmpuff runs `git status --porcelain=v2 -b -z`:

- **`--porcelain=v2`**: Uses git's structured status format, with typed entries for changed files, renames/copies, unmerged files, untracked files, and structured branch metadata.
- **`-b`** (branch info): Includes branch name plus ahead/behind counts as structured fields rather than a free-form header string.
- **`-z`** (null-delimited): Uses NUL terminators for reliable machine parsing of paths containing spaces and other special characters without shell quoting or escaping issues.

The v2 `-z` format still uses NUL separators, but it avoids a major porcelain v1 quirk: rename and copy entries remain typed records with explicit original and destination paths instead of relying on the v1 short-format `to\0from` field reversal.

## Parsing pipeline

The pipeline has three layers, each with a distinct responsibility.

### Layer 1: Low-level parsing — `github.com/mroth/porcelain/statusv2`

The external `porcelain` library reads the raw null-delimited byte stream and produces structured branch metadata plus typed entry values (`ChangedEntry`, `RenameOrCopyEntry`, `UnmergedEntry`, `UntrackedEntry`, `IgnoredEntry`). This layer is intentionally separate from scmpuff — it does format parsing with no scmpuff-specific rendering semantics, making it reusable for other tooling that needs to parse git porcelain output.

### Layer 2: Semantic conversion — `internal/gitstatus/porcelainv2/process.go`

`Process()` takes the parsed v2 status data from Layer 1 and translates it into scmpuff's domain types. This is where the two-character `XY` status codes from git are decoded into meaningful `ChangeType` values.

The conversion has two parts that run for each entry:

- For **changed** and **rename/copy** entries, `decodeXY()` maps the X column to staged changes and the Y column to unstaged worktree changes.
- For **unmerged** entries, `decodeUnmergedXY()` maps the seven conflict-specific XY combinations into the corresponding unmerged `ChangeType` values.

Because porcelain v2 has distinct entry types, unmerged states and untracked files are handled structurally rather than being inferred from the same short-format record shape as normal tracked changes.

The key insight is that a single git status entry can produce **two** `StatusItem`s. For example, `XY="MM"` means a file has staged modifications *and* unstaged modifications — this fans out into one Staged item and one Unstaged item, each appearing in its own section of the display. This fan-out keeps the renderer simple: every item is uniform and belongs to exactly one group.

The canonical table of XY meanings is in the git status man page under the porcelain v2 format description. scmpuff still uses the same `ChangeType` fan-out model as before; v2 mainly improves how entries and branch information are delivered to that conversion layer.

### Layer 3: Data structures — `internal/gitstatus/gitstatus.go`

```
StatusInfo
├── BranchInfo
│   ├── Name          string
│   ├── CommitsAhead  int
│   └── CommitsBehind int
└── Items []StatusItem
    └── StatusItem
        ├── ChangeType  (enum → Message(), State(), StatusGroup())
        ├── Path        string  (relative to repo root, always forward slashes)
        └── OrigPath    string  (for renames/copies, empty otherwise)
```

## Branch parsing

Porcelain v2 provides branch information as structured fields, so scmpuff no longer needs regex-based parsing of the free-form v1 `##` header. `extractBranch()` maps the parsed v2 branch info into scmpuff's `BranchInfo` type and preserves the existing detached-HEAD display label by translating v2's `(detached)` marker back to `HEAD (no branch)`.

## The ChangeType design

`ChangeType` is the central abstraction that bridges parsing and rendering. It's a flat enum with 20 named variants that each capture a specific combination of *where* a change is (staged, unstaged, unmerged, or untracked) and *what kind* of change it is (modified, new, deleted, renamed, etc). Each variant has three derived properties — `Message()`, `State()`, and `StatusGroup()` — backed by a single metadata lookup table, so adding a new variant is a one-line map entry. The renderer uses `StatusGroup()` for section grouping and section-level colors, and `State()` for per-item label colors (so staged and unstaged "modified" share the same label color even though they appear in different sections).

This table is the canonical reference for the 20 variants currently handled. The variant definitions, metadata, and XY decoding logic are spread across multiple source files; this table consolidates them in one place.

| ChangeType                   | Message()         | State()            | StatusGroup() |
|------------------------------|-------------------|--------------------|---------------|
| `ChangeStagedModified`       | `modified`        | `ModifiedState`    | `Staged`      |
| `ChangeStagedNewFile`        | `new file`        | `NewState`         | `Staged`      |
| `ChangeStagedDeleted`        | `deleted`         | `DeletedState`     | `Staged`      |
| `ChangeStagedRenamed`        | `renamed`         | `RenamedState`     | `Staged`      |
| `ChangeStagedCopied`         | `copied`          | `CopiedState`      | `Staged`      |
| `ChangeStagedType`           | `typechange`      | `TypeChangedState` | `Staged`      |
| `ChangeUnmergedDeletedBoth`  | `both deleted`    | `DeletedState`     | `Unmerged`    |
| `ChangeUnmergedAddedUs`      | `added by us`     | `NewState`         | `Unmerged`    |
| `ChangeUnmergedDeletedThem`  | `deleted by them` | `DeletedState`     | `Unmerged`    |
| `ChangeUnmergedAddedThem`    | `added by them`   | `NewState`         | `Unmerged`    |
| `ChangeUnmergedDeletedUs`    | `deleted by us`   | `DeletedState`     | `Unmerged`    |
| `ChangeUnmergedAddedBoth`    | `both added`      | `NewState`         | `Unmerged`    |
| `ChangeUnmergedModifiedBoth` | `both modified`   | `ModifiedState`    | `Unmerged`    |
| `ChangeUnstagedModified`     | `modified`        | `ModifiedState`    | `Unstaged`    |
| `ChangeUnstagedDeleted`      | `deleted`         | `DeletedState`     | `Unstaged`    |
| `ChangeUnstagedType`         | `typechange`      | `TypeChangedState` | `Unstaged`    |
| `ChangeUnstagedNewFile`      | `new file`        | `NewState`         | `Unstaged`    |
| `ChangeUnstagedRenamed`      | `renamed`         | `RenamedState`     | `Unstaged`    |
| `ChangeUnstagedCopied`       | `copied`          | `CopiedState`      | `Unstaged`    |
| `ChangeUntracked`            | `untracked`       | `UntrackedState`   | `Untracked`   |

## Test data and debugging

- **Porcelain parser tests**: `internal/gitstatus/porcelainv2/` — parser-specific tests for the current status conversion layer.
- **Shared regression fixtures**: `internal/gitstatus/testdata/` — debug-dump-derived fixtures that can be reused across parser versions.
- **Debug dump**: `scmpuff debug dump --archive` produces a diagnostic archive containing raw porcelain output in all formats. When a user reports a parsing failure, the archive can be integrated into `internal/gitstatus/testdata/` for shared regression coverage, or into `internal/gitstatus/porcelainv2/testdata/` when the fixture is specific to the v2 conversion layer.

## Porcelain v2 notes

scmpuff now uses `--porcelain=v2`. The main benefits realized by the migration are:

- **Structured branch information** — v2 provides branch data in a well-defined format, eliminating the regex-based header parsing that v1 required.
- **Typed entry records** — changed files, renames/copies, unmerged entries, untracked files, and ignored files are distinct record types instead of being inferred from one overloaded short format.
- **Simpler unmerged handling** — conflict entries are delivered separately, so the normal XY decoder no longer needs to guard against merge-conflict cases.

The older `internal/gitstatus/porcelainv1/` package is still kept in the repository as reference material and for comparison during maintenance, but it is no longer used by the `status` command.
