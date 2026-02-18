# Shell Integration

The shell layer is intentionally thin. Go handles all parsing, rendering, and argument expansion logic. The shell is responsible for two things that a subprocess can't do for itself: exporting environment variables into the user's session, and intercepting `git` commands before they run.

## How it works

### The status → environment variable loop

When the user runs `gs` (or `scmpuff_status`), the shell function calls `scmpuff status --filelist` and captures its output. The Go binary does all the real work — running `git status`, parsing the porcelain output, and rendering the numbered display. But it also emits a hidden first line: a tab-delimited list of absolute file paths, in the same order as the numbered display.

The shell function reads this first line, splits on tabs, and exports each path as a numbered environment variable: `$e1`, `$e2`, `$e3`, etc. Then it prints the remaining lines (the colorized status) to the terminal. Before each refresh, all existing `$eN` variables are cleared so stale entries from a previous run don't linger.

These environment variables are the bridge between the two halves of the system. The Go binary sets their values (indirectly, via the shell wrapper), and later reads them back when expanding shortcuts.

### The git wrapper

`scmpuff init` also installs a `git()` shell function that shadows the real git binary. When the user types something like `git add 1 2`, the wrapper intercepts it and routes it through `scmpuff exec`, which expands the numeric arguments. The expansion works by converting `1` → `$e1`, then resolving `$e1` via standard environment variable expansion to get the actual file path that was stored during the last status display.

The wrapper resolves the real git binary path into `$SCMPUFF_GIT_CMD` at init time (via `which git`) and uses that for all actual git invocations, avoiding infinite recursion.

Not all git subcommands need shortcut expansion. The wrapper uses a dispatch table:

| Subcommand(s)                                | Behavior                                                                      |
|----------------------------------------------|-------------------------------------------------------------------------------|
| `commit`, `blame`, `log`, `rebase`, `merge`  | `scmpuff exec -- git <args>` — expands shortcuts to absolute paths            |
| `checkout`, `diff`, `rm`, `reset`, `restore` | `scmpuff exec --relative -- git <args>` — expands shortcuts to relative paths |
| `add`                                        | `scmpuff exec -- git <args>` then auto-refreshes status via `scmpuff_status`  |
| everything else                              | Pass through to real git directly (no expansion)                              |

The `--relative` flag matters for commands like `diff` and `checkout` where git expects paths relative to cwd. The `add` case auto-refreshes status afterward so the numbered shortcuts immediately reflect the new state.

### Aliases

`scmpuff init` optionally installs short aliases (controlled by `--aliases`, default on):

| Alias | Expansion        |
|-------|------------------|
| `gs`  | `scmpuff_status` |
| `ga`  | `git add`        |
| `gd`  | `git diff`       |
| `gl`  | `git log`        |
| `gco` | `git checkout`   |
| `grs` | `git reset`      |

Since `git` is wrapped, `ga 1 2` effectively becomes `scmpuff exec -- git add 1 2` with auto-status-refresh.

## Initialization

Users add `eval "$(scmpuff init -s)"` to their shell profile (or `scmpuff init --shell=fish | source` for fish). The `--shell` flag selects the shell type; if omitted, it's detected from `$SHELL`. The init command emits a script to stdout that installs the `scmpuff_status()` function, the `git()` wrapper (if `--wrap`, default on), and aliases (if `--aliases`, default on).

Shell scripts are embedded in the binary at compile time via `go:embed`. Bash and zsh share the same scripts; fish has its own variants for the status and git wrapper scripts due to syntax differences. The aliases script is shared across all shells.

## Bash/zsh vs fish differences

| Aspect           | Bash/Zsh                        | Fish                                                       |
|------------------|---------------------------------|------------------------------------------------------------|
| Variable export  | `export $var="$value"`          | `set -gx "$var" "$value"`                                  |
| Field splitting  | `IFS=$'\t'; for file in $files` | `string split \t $output[1]`                               |
| Exit status      | `$?`                            | `$status`                                                  |
| Function erase   | `unset -f git`                  | `functions -e git`                                         |
| Output indexing  | `head -n 1` / `tail -n +2`     | `$cmd_output[1]` / `$cmd_output[2..-1]`                    |
| Which command    | `\which git`                    | `which git`                                                |
| Passthrough exec | `"$SCMPUFF_GIT_CMD" "$@"`      | `eval command "$SCMPUFF_GIT_CMD" (string escape -- $argv)` |
