---
description: Scaffold a new Architecture Decision Record (ADR) from the template in docs/adr/
argument-hint: <NNNN> <title>
allowed-tools: Bash(date *), Read, Write, Glob
---

# /adr-new

Create a new Architecture Decision Record in `docs/adr/NNNN-<slug>.md` based on [docs/adr/0000-template.md](../../docs/adr/0000-template.md).

## Inputs

`$ARGUMENTS` must be `<NNNN> <title>` where:
- `<NNNN>` — zero-padded 4-digit number (e.g., `0003`). MUST NOT already exist.
- `<title>` — short title (e.g., "Ship separate arm64 binaries")

If `<NNNN>` is missing, Glob `docs/adr/*.md` and suggest the next available number (ignoring `0000-template.md` and `README.md`).
If `<title>` is missing, ask via `AskUserQuestion`.

## Protocol

1. **Resolve date** via `date +%F`.
2. **Check the number is not taken**: Glob `docs/adr/<NNNN>-*.md`. Abort if it exists.
3. **Derive slug** from the title: lowercase, kebab-case, ≤40 chars.
4. **Read** `docs/adr/0000-template.md`.
5. **Write the file** at `docs/adr/<NNNN>-<slug>.md`, substituting:
   - Title placeholder → `<NNNN> — <Title>`
   - `YYYY-MM-DD` → today's date
   - Leave `Deciders`, `Context`, `Decision`, `Consequences`, `Alternatives considered`, `References` empty for the user to fill in
   - Initial status: `proposed`
6. **Update the index** in `docs/adr/README.md`: append a new row to the "Index" table with status `proposed`.
7. **Report** the file path and suggest that the user fill in Context before merging.

## Rules

- **Use the template file as the source of truth.** If it changes, this command picks up the new structure automatically.
- **Never skip ADR numbers.** If `0003` is taken and the user asks for `0005`, refuse and say `0004` is next.
- **Initial status is `proposed`.** It becomes `accepted` when the PR introducing it merges.
- **No code changes in ADRs.** ADRs are decisions and rationale; code lives elsewhere.

## Gotchas

- If you accidentally pick a number that's in-flight on an open PR, the merge will conflict. When unsure, check open PRs via `gh pr list --search "docs/adr"`.
- Keep ADRs under 2 pages. If you need more, you probably need a spec in `docs/specs/` instead.
