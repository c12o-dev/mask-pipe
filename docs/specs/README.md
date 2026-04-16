# Specifications

This directory contains the canonical specifications for `mask-pipe`. Specs describe **what** the tool does and **why** — the code in this repository is the implementation of these specs.

## Why specs?

- A single source of truth for behavior — the code and the docs can't disagree silently.
- Onboarding new contributors: read the specs first, the code second.
- Changes to behavior are traceable — every breaking change starts with a spec PR.

## Spec index

| # | File | Scope |
|---|---|---|
| 000 | [000-vision.md](000-vision.md) | Product vision, non-goals, positioning |
| 001 | [001-cli-interface.md](001-cli-interface.md) | CLI commands, flags, exit codes, I/O contract |
| 002 | [002-pattern-library.md](002-pattern-library.md) | Built-in patterns, precision requirements |
| 003 | [003-config-format.md](003-config-format.md) | TOML configuration file schema |
| 004 | [004-shell-integration.md](004-shell-integration.md) | Recommended patterns, TTY preservation rules |

## Spec format

Each spec is a standalone Markdown file with:

```markdown
# NNN — Title

- **Status:** draft | accepted | superseded-by-NNN
- **Last updated:** YYYY-MM-DD

## Summary
One paragraph on what this spec covers.

## Motivation
Why this exists. What problem it solves.

## Specification
The actual contract. Use MUST / SHOULD / MAY in the RFC 2119 sense.

## Non-goals
What this spec deliberately does NOT cover.

## Open questions
Unresolved issues that should be tracked as GitHub issues.
```

## Changing a spec

Specs are immutable-ish — they can be changed, but the change must go through the [spec change workflow](../../CONTRIBUTING.md#proposing-a-spec-change):

1. Open a `spec_change` issue describing the proposed change
2. Open a PR updating the spec file (and only the spec, for clarity)
3. After merge, implementation PRs can land that realize the new behavior

Superseded specs are kept in place but marked with `Status: superseded-by-NNN` so history is preserved.
