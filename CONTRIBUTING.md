# Contributing to mask-pipe

Thanks for your interest. This project uses a **hybrid spec-driven + issue-driven workflow** — the short version is: specs define *what* we build, issues define *when and why* we do it, and every PR links to both.

## Ground rules

- **Specs are the source of truth.** Behavior, contracts, and built-in patterns live in [`docs/specs/`](docs/specs/). If the code disagrees with the spec, the spec wins — either fix the code or file an issue to update the spec.
- **Every change is tracked by an issue.** No "drive-by" PRs. If you have a fix in mind, open a bug report first so the change has context.
- **Spec changes are themselves issues.** Use the [spec change template](.github/ISSUE_TEMPLATE/spec_change.yml) to propose a behavior change before opening an implementation PR.

## Workflow

### Proposing a new feature or behavior change

```
1. Open a feature issue      → discussion & scope
2. Open a spec PR            → update docs/specs/NNN-*.md
3. After spec merges, open   → implementation PR (Fixes #issue)
   an implementation PR
```

For small, obvious features you can combine steps 2 and 3 into a single PR that updates both the spec and the code.

### Reporting a bug

Use the [bug report template](.github/ISSUE_TEMPLATE/bug.yml). Include:

- The exact `mask-pipe` version (`mask-pipe --version`)
- A minimal input that reproduces the bug
- Expected vs. actual output
- Your OS and terminal emulator

### Proposing a new built-in pattern

Use the [pattern proposal template](.github/ISSUE_TEMPLATE/pattern.yml). A pattern proposal must include:

- The regex
- **5+ real-world examples** the pattern should match
- **5+ examples it should NOT match** (false-positive prevention)
- Rationale for precision-over-recall tuning
- Source reference (e.g., vendor docs) for the format

Patterns with high false-positive rates will be rejected or moved behind an opt-in flag.

### Proposing a spec change

Use the [spec change template](.github/ISSUE_TEMPLATE/spec_change.yml). Describe:

- Which spec file is affected
- The current behavior (quote the relevant line)
- The proposed behavior
- Migration notes if the change is breaking

## Development setup

```bash
git clone https://github.com/c12o-dev/mask-pipe.git
cd mask-pipe
go test ./...
go build -o mask-pipe .
```

### Before submitting a PR

- [ ] Linked issue number in the PR description (`Fixes #N` or `Refs #N`)
- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] `golangci-lint run` is clean (if configured)
- [ ] Spec updated if behavior changed
- [ ] New patterns have 5+ match and 5+ no-match test cases
- [ ] No new dependencies without an ADR justifying them

## Architecture Decision Records

Big decisions (language choice, major dependencies, design philosophies) are recorded as ADRs in [`docs/adr/`](docs/adr/). If your PR includes a decision like this, add a new ADR file — use [`docs/adr/0000-template.md`](docs/adr/0000-template.md) as a starting point.

## Code of conduct

Be kind, be specific, and assume good faith. Security-adjacent projects attract passionate debates about false positives vs. false negatives — those are design tradeoffs, not moral failings. Disagree with the code, not the person.

## Questions?

Open a [Discussion](https://github.com/c12o-dev/mask-pipe/discussions) for anything that isn't a bug, feature, or pattern proposal.
