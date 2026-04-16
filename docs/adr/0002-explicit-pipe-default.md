# 0002 — Explicit pipe as default integration pattern

- **Status:** accepted
- **Date:** 2026-04-16

## Context

Users naturally want a "set it and forget it" install — add a line to `.zshrc`, have all terminal output masked automatically. 1Password's `op run` tried exactly this (transparent stdout masking of child processes) and the consequences are well-documented:

- TUI applications break (charmbracelet/crush issue #1415)
- Colored output from `ls`, `grep`, `git` disappears
- Ansible integration breaks because `op run` sets up pipes that make `isatty()` return false
- The only workaround is `--no-masking`, which defeats the purpose

The core issue is that inserting a filter process into the stdout pipeline of a child process causes that child to see a pipe (non-TTY) instead of a terminal. This is a fundamental Unix constraint, not an implementation bug.

## Decision

mask-pipe's **primary and default** integration pattern is the explicit pipe:

```bash
<command> | mask-pipe
```

Users must consciously pipe commands through mask-pipe. Automatic wrapping of arbitrary commands is rejected for v1.

Shell function wrappers (`mlogs() { docker logs "$@" | mask-pipe; }`) are documented as an opt-in convenience for non-interactive commands only, with prominent warnings against wrapping interactive or TUI commands.

A PTY-proxy mode (`mask-pipe shell`) is reserved for v2 but not committed.

## Consequences

**Positive:**

- mask-pipe does not break TUI apps, interactive prompts, or color output by default
- Users see clearly WHEN masking is applied; no silent failures
- The default install has zero footprint on the user's shell behavior
- Positions mask-pipe favorably against `op run`'s documented issues — "we intentionally avoided the pitfall"

**Negative / trade-offs:**

- Users must remember to add `| mask-pipe` to commands — one more keystroke
- Some users will be disappointed there's no auto-mode; we must educate via README and [docs/specs/004-shell-integration.md](../specs/004-shell-integration.md)
- "Set and forget" demand is real; we defer it to v2 `mask-pipe shell` with a PTY proxy, which is a significant future implementation effort

## Alternatives considered

- **Auto-wrap via `LD_PRELOAD`** — Intercept `write()` syscalls at the libc level. Rejected because:
  - Static Go binaries (docker, kubectl) don't use libc for I/O, breaking the mechanism
  - Security-sensitive environments often disable `LD_PRELOAD`
  - Multi-threading correctness is hard to get right

- **Auto-wrap via zsh `preexec` hook** — Rewrite commands before execution. Rejected because:
  - Unreliable (subshells, `exec`, and explicit paths bypass it)
  - zsh `exec | pipe` redirects don't work cleanly; requires fragile `coproc` machinery
  - Same TTY-breaking problem as `op run` once the pipe is inserted

- **Ship `mask-pipe init zsh` that installs wrappers for common commands** — Partial solution. Deferred: we document the pattern in README but don't bundle a hook-install step in v1. Users who want it can copy the examples.

- **Run under `unbuffer`/`script`** — Preserves TTY for the wrapped command but pipe output is still passed through another filter that loses TTY downstream. Not a complete solution.

## References

- [docs/specs/004-shell-integration.md](../specs/004-shell-integration.md) — full policy
- 1Password `op run` TTY issues: https://www.1password.community/discussions/developers/op-run-changes-stdout-and-stderr-to-not-be-ttys-when-masking/26040
- charmbracelet/crush TUI breakage: https://github.com/charmbracelet/crush/issues/1415
