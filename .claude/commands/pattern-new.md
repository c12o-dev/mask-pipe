---
description: Scaffold a new built-in pattern — updates the spec, creates pattern code, and creates a test file
argument-hint: <pattern_id>
allowed-tools: Bash(date *), Read, Write, Edit, Glob, Grep
---

# /pattern-new

Scaffold all the artifacts required to add a new built-in pattern to mask-pipe. This command updates [docs/specs/002-pattern-library.md](../../docs/specs/002-pattern-library.md), creates a Go source file, and creates a matching test file — mirroring the workflow in [CONTRIBUTING.md § Proposing a new built-in pattern](../../CONTRIBUTING.md#proposing-a-new-built-in-pattern).

## Inputs

`$ARGUMENTS` is `<pattern_id>` — a stable kebab-case identifier (e.g., `azure_storage_key`).

If missing, ask via `AskUserQuestion`.

Then gather the following via a single `AskUserQuestion` (multi-question supported):
- **Regex** — Go RE2 syntax
- **Source URL** — vendor documentation defining the format
- **Rationale** — one sentence on why this belongs in the default set

## Protocol

1. **Verify the pattern does not already exist**:
   - Grep `docs/specs/002-pattern-library.md` for the pattern_id. Abort if present.
   - Grep `patterns/` for the pattern_id (if the directory exists). Abort if present.

2. **Remind the user of the 5+5 rule** (do not proceed without acknowledgment):
   - ≥5 positive examples (strings the pattern MUST match)
   - ≥5 negative examples (strings it MUST NOT match, chosen to prevent false positives)
   - These go in the test file; the user fills them in after this command scaffolds the structure.

3. **Update the spec** (`docs/specs/002-pattern-library.md`):
   - Find the "Initial pattern set" table and append a new row
   - Do NOT modify other sections
   - Bump the "Last updated" frontmatter to today's date via `date +%F`

4. **Scaffold the pattern code** at `patterns/<pattern_id>.go`:
   - Follow the `Pattern` struct shape defined in the spec (ID, Name, Regex, CaptureIdx, Replacement, Examples, NonExamples, Source)
   - Leave `Examples` and `NonExamples` as empty slices with a `// TODO: add 5+ entries each` comment
   - Include a `func init()` that registers the pattern in the package-level registry

5. **Scaffold the test** at `patterns/<pattern_id>_test.go`:
   - `TestPattern<Name>_Matches` — runs the regex against every `Examples` entry, asserts match
   - `TestPattern<Name>_NonMatches` — runs against every `NonExamples` entry, asserts no match
   - `TestPattern<Name>_NoFalsePositivesOnRealisticText` — runs against the shared realistic-corpus fixture, asserts zero matches

6. **Report to the user**:
   - The 3 files created/modified
   - The 3 TODO items remaining (fill examples, run tests, open PR)
   - Suggested next step: `gh issue create --template pattern.yml` to file the governance issue referencing this scaffold

## Rules

- **Never fill in examples on behalf of the user.** The 5+5 examples must come from the proposer's research — they are the evidence base for the precision claim.
- **Spec file update is limited to the table row and the `Last updated` date.** Do not touch other rows or sections.
- **No `git add` / `git commit`.** File creation only.
- **If the regex does not compile** (e.g., Go-invalid syntax the user typed), abort with the regex error and do not write any files.

## Gotchas

- The Go package path depends on the project structure. At scaffold time, Read `go.mod` (if present) to get the module path; otherwise fall back to `package patterns` with a TODO.
- The pattern ID is used as both the file name and the struct constant — lowercase with underscores in Go, kebab-case in the spec table. Convert appropriately.
- The "realistic corpus" test requires `patterns/testdata/realistic_corpus.txt` to exist. If it doesn't, add a TODO in the test for the user to create it.
