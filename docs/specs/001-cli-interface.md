# 001 — CLI Interface

- **Status:** draft
- **Last updated:** 2026-04-16

## Summary

Defines the command-line interface contract: commands, flags, input/output handling, exit codes, and error messages.

## Specification

### Invocation

```
mask-pipe [flags] [subcommand [flags]]
```

When invoked with no subcommand, mask-pipe operates in **filter mode**: it reads from stdin, applies the configured patterns, and writes the masked output to stdout.

### Filter mode (default)

```
$ <some-command> | mask-pipe [flags]
```

**Input:**
- `stdin` — the raw byte stream to filter
- mask-pipe MUST operate line-by-line to preserve streaming behavior (e.g., `tail -f`)
- **Exception:** multi-line patterns (e.g., PEM private keys) use begin/end marker buffering. Lines between a begin marker and its corresponding end marker are accumulated before matching. A safety limit (100 lines or 64KB) prevents unbounded buffering; if exceeded, buffered lines are flushed unmodified and a warning is written to stderr. See [ADR 0004](../adr/0004-multiline-buffering.md).
- Lines MAY exceed the default `bufio.MaxScanTokenSize` (64KB); mask-pipe MUST handle lines up to at least 1MB without crashing

**Output:**
- `stdout` — the masked byte stream
- `stderr` — diagnostic output only (warnings, dry-run markers). MUST NOT interleave with normal output.
- Output MUST be flushed after each input line to preserve streaming latency

**Exit codes:**
- `0` — normal termination (stdin closed)
- `1` — unrecoverable runtime error (e.g., invalid config)
- `2` — invalid command-line flags or arguments
- `64`–`78` — reserved for future use (following BSD sysexits.h conventions)

### Flags (filter mode)

| Flag | Default | Description |
|---|---|---|
| `--config <path>` | `~/.mask-pipe.toml` or `$XDG_CONFIG_HOME/mask-pipe/config.toml` | Path to TOML config file. If the file does not exist, mask-pipe uses built-in defaults and continues (not a fatal error). |
| `--dry-run` | `false` | Do not replace matches; highlight them in the output (with ANSI color if stdout is a TTY). Useful for previewing what would be masked. |
| `--no-color` | `false` (auto) | Disable ANSI color output. Auto-enabled when stdout is not a TTY, or when `NO_COLOR` env var is set. |
| `--mask-char <char>` | `*` | Override the masking character from config. |
| `--show-tail <N>` | From config (default `4`) | Show the last N characters of each masked value. Set to `0` to mask entirely. |
| `--help` / `-h` | — | Print help and exit `0`. |
| `--version` / `-V` | — | Print version and exit `0`. |
| `--strict` | `false` | (v2) Enable the high-recall pattern set for aggressive matching. Reserved. |

### Subcommands

```
mask-pipe doctor         — Diagnose configuration and patterns
mask-pipe list-patterns  — Print all active built-in and custom patterns
mask-pipe version        — Print version and build info
```

#### `mask-pipe doctor`

Runs self-diagnostic checks and prints a report to stdout. Exit `0` if all checks pass, `1` otherwise.

Checks:
- Config file exists and is parseable (or gracefully absent)
- All custom regex patterns compile
- Built-in patterns are enabled as expected
- stdout is writable

#### `mask-pipe list-patterns`

Prints a table of active patterns: name, source (builtin / custom), regex, sample match. For debugging custom configs.

#### `mask-pipe version`

Prints version, build date, and Git SHA. Machine-readable with `--json` flag.

### Error messages

Error messages MUST:

- Be written to stderr, not stdout
- Include the config file path when the error relates to configuration
- Suggest a next action ("run `mask-pipe doctor` to diagnose", "see docs/specs/003-config-format.md")
- Never include the user's input data (avoid leaking the secrets we're trying to mask)

### TTY behavior

mask-pipe itself does NOT allocate a PTY or wrap child processes in the default filter mode. It is a pure stream filter. This is a deliberate design choice — see [004-shell-integration.md](004-shell-integration.md).

## Non-goals

- Interactive TUI (no `mask-pipe` with a UI)
- Shell integration (no `eval "$(mask-pipe init zsh)"` — users write their own wrappers)
- Automatic stdout wrapping of arbitrary commands (rejected; see [004-shell-integration.md](004-shell-integration.md))
- Config-file hot reload (restart mask-pipe to pick up changes)

## Open questions

- Should `mask-pipe` support JSON input with pattern-matching only on specific fields (as PIMO does)? Tracked as a v2 consideration.
- Should there be a `--allowlist` flag to exempt specific strings from masking? See [002-pattern-library.md §future](002-pattern-library.md).
- Should the binary be renamed to a shorter command like `mp`? Trade-off: shorter typing vs. collision risk with other tools (there's an `mp` Python CLI). Defer until user feedback.
