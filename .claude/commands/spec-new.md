---
description: Scaffold a new specification document from the standard template in docs/specs/
argument-hint: <NNN> <title>
allowed-tools: Bash(date *), Read, Write, Glob
---

# /spec-new

Create a new specification document in `docs/specs/NNN-<slug>.md` following the format defined in [docs/specs/README.md](../../docs/specs/README.md).

## Inputs

`$ARGUMENTS` must be `<NNN> <title>` where:
- `<NNN>` — zero-padded 3-digit number (e.g., `005`). MUST NOT already exist.
- `<title>` — human-readable title (e.g., "Rate limiting policy")

If `<NNN>` is missing, Glob `docs/specs/*.md` and suggest the next available number.
If `<title>` is missing, ask via `AskUserQuestion` with one concise question.

## Protocol

1. **Resolve date** via `date +%F`.
2. **Check the number is not taken**: Glob `docs/specs/<NNN>-*.md`. If it exists, abort and show the existing file.
3. **Derive slug** from the title: lowercase, kebab-case, ≤30 chars, no trailing punctuation.
4. **Write the file** at `docs/specs/<NNN>-<slug>.md` using the template below.
5. **Update the index** in `docs/specs/README.md`: append a new row to the "Spec index" table.
6. **Report to the user**: the file path and a suggested next action (e.g., "open a spec_change issue linking this draft").

## Template

```markdown
# <NNN> — <Title>

- **Status:** draft
- **Last updated:** <YYYY-MM-DD>

## Summary

One paragraph describing what this spec covers.

## Motivation

Why this exists. What problem it solves. What happens if we don't have it.

## Specification

The actual contract. Use MUST / SHOULD / MAY in the RFC 2119 sense.

## Non-goals

What this spec deliberately does NOT cover.

## Open questions

Unresolved issues that should be tracked as GitHub issues.
```

## Rules

- **Never overwrite an existing spec.** The `<NNN>-*.md` number is a primary key; collisions abort the command.
- **Do not auto-fill content.** The template is intentionally minimal — the user fills in Motivation, Specification, etc.
- **Do not run `git add` or `git commit`.** File creation only — the user controls versioning.
- **Status starts as `draft`.** The user changes it to `accepted` via a spec_change issue + PR.

## Gotchas

- The index table in `docs/specs/README.md` must be kept in sync. Failing to update it breaks spec discoverability.
- If the user provides an existing number (e.g., to supersede), refuse and suggest opening a spec_change issue instead.
