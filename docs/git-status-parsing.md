# Git Status Parsing

scmpuff parses `git status` output into structured data for numbered file display. This document covers the parsing pipeline from raw git output to the `StatusInfo` data structure, and the design of the `ChangeType` abstraction that bridges parsing and rendering. For how that data is rendered to the screen, see the "Status rendering" section in [architecture.md](architecture.md).

## Raw format

scmpuff runs `git status -z -b`, which is equivalent to `git status --porcelain=v1 -b -z`:

- **`-z`** (null-delimited): Reliable cross-platform machine parsing. Avoids shell quoting issues and special character escaping. The `-z` flag implies `--porcelain=v1` when no explicit porcelain version is given, which provides backward compatibility with very old git versions.
- **`-b`** (branch info): Adds a header line with branch name and ahead/behind counts.

The `-z` format has quirks compared to the normal porcelain output: the `->` is omitted from rename entries (field order is reversed to `to\0from`), NUL replaces both field separators and line terminators, and filenames are never quoted or escaped.

See [Porcelain v2 migration](#porcelain-v2-migration) below for planned future work.

## Parsing pipeline

The pipeline has three layers, each with a distinct responsibility.

### Layer 1: Low-level parsing — `github.com/mroth/porcelain`

The external `porcelain` library reads the raw null-delimited byte stream and produces structured entry values with XY status codes, file paths, and original paths (for renames). This layer is intentionally separate from scmpuff — it does pure tokenization with no semantic interpretation, making it reusable for other tooling that needs to parse git porcelain output.

### Layer 2: Semantic conversion — `internal/gitstatus/porcelainv1/process.go`

`Process()` takes the tokenized entries from Layer 1 and translates them into scmpuff's domain types. This is where the two-character `XY` status codes from git are decoded into meaningful `ChangeType` values.

The conversion has two parts that run for each entry:

- A **primary decoder** examines both characters together to detect unmerged states (merge conflicts), untracked files, and staged changes (from the X column).
- A **secondary decoder** examines the Y column for unstaged worktree changes (modified, deleted, typechange).

The key insight is that a single git status entry can produce **two** `StatusItem`s. For example, `XY="MM"` means a file has staged modifications *and* unstaged modifications — this fans out into one Staged item and one Unstaged item, each appearing in its own section of the display. This fan-out keeps the renderer simple: every item is uniform and belongs to exactly one group.

The canonical table of all possible XY combinations is in the git status man page under "OUTPUT / Short Format", also archived in the porcelain library at https://github.com/mroth/porcelain/blob/main/docs/git-status.txt.

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

Porcelain v1 does not provide branch data in a structured format — it encodes it as a `##` header line with several variant formats depending on the repository state (normal tracking, no upstream, initial commit, detached HEAD, ahead/behind counts). Since there's no machine-friendly structure to parse, the branch parser uses regex patterns to extract the branch name and ahead/behind counts from these header variants. This is one of the motivations for the eventual [porcelain v2 migration](#porcelain-v2-migration), which provides branch data in a well-defined structured format.

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

- **Porcelain parser tests**: `internal/gitstatus/porcelainv1/` — test files for parsing logic.
- **Debug dump**: `scmpuff debug dump --archive` produces a diagnostic archive containing raw porcelain output in all formats. When a user reports a parsing failure, the archive can be easily integrated into the `porcelainv1/testdata/` directory to add a regression test for their repo's status output.

## Porcelain v2 migration

More recent git versions support `--porcelain=v2`, which provides a cleaner structured format that addresses design mistakes in porcelain v1. The `github.com/mroth/porcelain` library already implements a v2 parser. Migrating scmpuff to v2 is a desired future change. Benefits include:

- **Structured branch information** — v2 provides branch data in a well-defined format, eliminating the regex-based header parsing currently needed for v1.
- **Cleaner entry format** — avoids the `-z` mode quirks (reversed rename field order, NUL delimiter ambiguities).
