# 004 — Shell Integration Policy

- **Status:** draft
- **Last updated:** 2026-04-16

## Summary

Defines how mask-pipe is intended to be integrated into a user's shell workflow, and why automatic stdout wrapping of arbitrary commands is rejected as a default.

## Motivation

Users naturally ask: "Can I just set this up once in my `.zshrc` and have it mask everything automatically?" This spec documents why the answer is "no, by design" — and what we recommend instead.

## Specification

### Recommended: explicit pipe (default pattern)

```bash
<command> | mask-pipe
```

This is the primary, documented usage. It is:

- **Safe** — does not wrap processes; does not allocate a PTY
- **Transparent** — users see exactly when masking is applied
- **Composable** — works with any command that writes to stdout
- **Preserves TTY for the upstream command** — since mask-pipe only reads the pipe, the upstream command's own TTY (if any) is unaffected; only mask-pipe's own stdout-to-terminal connection is a pipe from the shell's perspective

### Acceptable: shell function wrappers for non-interactive commands

For commands whose output is strictly non-interactive (logs, dumps), users MAY define shell functions:

```zsh
# Safe to wrap: these produce log output, not interactive sessions
mlogs()  { docker logs "$@"  | mask-pipe; return ${PIPESTATUS[0]}; }
mklogs() { kubectl logs "$@" | mask-pipe; return ${PIPESTATUS[0]}; }
menv()   { env                | mask-pipe; }
```

mask-pipe's documentation MUST include this pattern as an example but MUST also warn that:

- Wrapping interactive commands breaks them (see "Rejected" section)
- `${PIPESTATUS[0]}` is required to preserve the underlying command's exit code
- These wrappers are user-defined, not installed by mask-pipe itself

### Rejected: automatic wrapping of all commands

mask-pipe MUST NOT provide:

- A `mask-pipe init zsh` / `eval "$(mask-pipe init)"` style shell hook
- An `LD_PRELOAD` library that intercepts all `write()` syscalls
- A `preexec` hook that transparently rewrites commands

**Rationale:**

1. **TTY breakage is unavoidable.** Piping any command's output causes `isatty(1)` to return `false` in the child process. This disables:
   - Colored output from tools that check TTY (`ls --color=auto`, `grep --color=auto`, `git`)
   - Interactive prompts (`sudo`, `gh auth login`)
   - TUI applications (`vim`, `less`, `htop`, `lazygit`, any Bubbletea/ncurses app)
   - Progress bars in installers (`npm install`, `cargo build`)

2. **This is exactly why 1Password `op run` fails.** The community has documented that `op run`'s stdout masking destroys TTY behavior, breaks TUI apps, and breaks Ansible integration. We will not repeat that mistake.

3. **Explicit beats implicit for security tools.** Users must know WHEN masking is applied. Silent automatic masking creates false confidence and hides bugs.

See [CLAUDE.md — Workflow](../../CLAUDE.md) and competitor analysis in the design repo for the full reasoning.

### Future: `mask-pipe shell` subcommand (v2, reserved)

A future major version MAY introduce a PTY-proxy mode that starts a wrapped shell session with transparent masking:

```bash
mask-pipe shell
```

This mode:

- Allocates a PTY pair; starts `$SHELL` as a child of the PTY slave
- Copies bytes from the slave → (regex filter) → the user's terminal
- Preserves `isatty()` = true inside the wrapped shell
- Enables Ctrl-C, SIGWINCH, and other TTY-dependent behaviors

**Deferred because:**

- Implementation complexity is high (SIGWINCH propagation, escape-sequence boundary handling, raw-mode management)
- The explicit-pipe default already covers the primary use cases (screen sharing, log review)
- Risk of breaking edge cases (vim, tmux-in-tmux) requires extensive testing

The design is reserved, not committed. It will be specified in a future `005-pty-proxy-shell.md` when implementation begins.

### Commands users MUST NOT wrap

mask-pipe's documentation MUST explicitly warn that these commands break when piped:

| Command | Why |
|---|---|
| `docker run -it ...` | Requires interactive TTY |
| `kubectl exec -it ...` | Requires interactive TTY |
| `vim`, `less`, `nano`, any editor | Full-screen TUI requires TTY |
| `sudo <interactive-command>` | Password prompt requires TTY |
| `git` (with pager enabled) | Auto-pager detection requires TTY |
| `npm install`, `yarn`, `pnpm`, `cargo build` | Progress bars and color output rely on TTY detection |
| `gh pr view`, `gh issue view` | TUI rendering |
| Any shell login (`bash -i`, `zsh`) | Interactivity |

### User education

The README, CONTRIBUTING.md, and any `mask-pipe --help` output MUST:

- Lead with the explicit-pipe pattern (`cmd | mask-pipe`)
- Show shell function wrappers as an opt-in convenience, with the "do not wrap interactive commands" warning prominent
- Reference this spec document when users ask "why not auto-mask everything"

## Non-goals

- "Zero-config auto-mask everything" mode — rejected by design
- Shell-specific plugins (oh-my-zsh, bash-it modules) — users can write their own; we don't maintain ecosystem-specific packages
- Integration with Warp Terminal's display-layer redaction — different layer, different concern

## Open questions

- Should we publish an official `mask-pipe-wrappers.zsh` file with a curated set of safe wrappers, or leave it to users?
- For the future `mask-pipe shell`, should the filter run in a separate goroutine with backpressure, or synchronously?
- Can we detect common unsafe wrappers (user wraps `docker` with `mask-pipe`) and warn at shell startup? Probably out of scope.
