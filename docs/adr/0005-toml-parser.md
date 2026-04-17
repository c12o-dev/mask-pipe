# 0005 — TOML parser: BurntSushi/toml

- **Status:** accepted
- **Date:** 2026-04-17

## Context

Spec 003 mandates TOML as the config format. Go's standard library does not include a TOML parser. We need an external dependency — the first one added to this project.

## Decision

Use `github.com/BurntSushi/toml` (v1.x).

## Consequences

**Positive:**

- De facto standard Go TOML parser — written by the TOML spec co-author
- Minimal API (`toml.DecodeFile` is all we need), no transitive dependencies
- Full TOML v1.0 spec compliance
- Active maintenance with a stable v1 API

**Negative / trade-offs:**

- First external dependency — binary size increases ~200KB
- Must track upstream for security fixes (low risk; the library is pure Go with no syscalls)

## Alternatives considered

- **pelletier/go-toml v2** — Good performance but larger API surface and more transitive dependencies. No clear advantage for our simple config.
- **Hand-rolled parser** — TOML v1.0 is deceptively complex (datetime types, inline tables, dotted keys). Not worth the effort or the bugs.
- **Switch to JSON/YAML** — Rejected by spec 003. TOML is the better fit for user-edited config files (comments, readable syntax).

## References

- [docs/specs/003-config-format.md](../specs/003-config-format.md)
- [Issue #18](https://github.com/c12o-dev/mask-pipe/issues/18)
