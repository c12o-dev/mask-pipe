# 000 — Product Vision

- **Status:** draft
- **Last updated:** 2026-04-16

## Summary

`mask-pipe` is a command-line tool that masks secrets in terminal stdout in real-time, intended to be placed in a Unix pipeline. It fills a gap between commit-time secret scanners (secretlint) and post-hoc leak detectors (TruffleHog): **it prevents secrets from hitting your screen in the first place.**

## Motivation

Terminal output is treated as a "trusted space" but frequently leaks through channels engineers do not think about:

- **Screen sharing** during Zoom/Meet/Discord calls can expose `env` and `docker-compose config` to every participant
- **Pair programming** via tmux or VS Code Live Share synchronizes terminal output to collaborators
- **Scrollback buffers** in iTerm2 / Windows Terminal retain thousands of lines with embedded secrets
- **AI assistants** (ChatGPT, Claude, Cursor) receive pasted terminal output that may contain credentials
- **Recorded terminal sessions** (asciinema, `script(1)`) persist secrets on disk

Existing tools do not cover this vector:

| Tool | What it protects | Gap |
|---|---|---|
| secretlint | Pre-commit (files) | Doesn't filter runtime stdout |
| TruffleHog | Post-hoc (repos, logs) | Detects leaks, doesn't prevent screen exposure |
| GitHub `add-mask` | CI logs only | Not usable locally |
| Doppler / Infisical | Secret storage/injection | Doesn't mask at display time |

`mask-pipe` is the missing piece: real-time masking of secrets as they flow through a Unix pipe.

## Specification

### What mask-pipe IS

- A single-binary CLI that reads from stdin, masks patterns matching known secret formats, and writes the masked output to stdout
- Configuration-free by default — ships with high-precision patterns for common secret formats (AWS, GitHub, Stripe, JWT, DB URLs, etc.)
- Opt-in customization via a TOML config file
- Cross-platform: macOS, Linux, Windows (amd64 + arm64)
- Positioned as **complementary** to existing secret-scanning tools, not a replacement

### What mask-pipe is NOT (non-goals)

- **Not a secret manager.** It does not store, generate, or inject secrets. Use Doppler, 1Password, or Vault for that.
- **Not a secret scanner.** It does not detect secrets in codebases or git history. Use TruffleHog or secretlint for that.
- **Not an automatic wrapper for every command.** Auto-piping all output breaks TTY behavior (colors, TUI apps, interactive prompts). See [004-shell-integration.md](004-shell-integration.md) for the explicit-pipe-by-default policy.
- **Not a DLP platform.** It is not a replacement for enterprise data-loss-prevention tools with audit trails, policy enforcement, and compliance reporting.
- **Not a terminal emulator feature.** Warp Terminal does display-level redaction; mask-pipe works in any terminal by filtering the data stream itself.

### Design principles

1. **Precision over recall.** A false positive that shreds normal output is worse than a false negative. Built-in patterns target precision ≥99% on realistic input.
2. **Explicit over implicit.** Users pipe commands through mask-pipe on purpose. Automatic stdout masking (as attempted by 1Password `op run`) breaks TTY detection and is rejected as a default.
3. **Local-first.** No telemetry, no cloud dependency, no account required. The free version runs entirely offline.
4. **Open core.** The masking engine and built-in patterns are open source under MIT. Optional paid features (team pattern sync, audit logs) are separate and do not lock out the core use case.
5. **Single binary.** No runtime dependencies. `curl | tar x` and run.

## Target users

**Primary:** Fully-remote developers (Zoom/Meet frequent), DevOps/SRE engineers reviewing logs, security-conscious individual contributors.

**Secondary:** Engineering teams that want to standardize masking patterns across members via shared TOML config.

**Out of scope (initially):** Large enterprises with DLP compliance requirements — those have different tooling needs.

## Open questions

- Should there be a `mask-pipe shell` subcommand that PTY-wraps a shell session to avoid the explicit-pipe requirement? See [004-shell-integration.md §v2](004-shell-integration.md).
- How aggressive should the default pattern set be? See [002-pattern-library.md](002-pattern-library.md) for the current balance.
- Is there a role for locale-specific pattern packs (e.g., Japanese My Number detection)? Track as a potential extension.
