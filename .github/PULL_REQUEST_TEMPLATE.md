<!--
Thanks for the contribution. Please fill out the sections below so reviewers can
verify that this PR follows the spec-driven + issue-driven workflow.

See CONTRIBUTING.md for details.
-->

## Linked issue

<!-- Every PR must link to an issue. Use "Fixes #N" or "Refs #N". -->

Fixes #

## What changed

<!-- Summary of the change in 1–3 sentences. Focus on the user-visible behavior difference. -->

## Spec / ADR impact

- [ ] No spec changes (bug fix matching existing spec)
- [ ] Spec updated in this PR (list files):
  - `docs/specs/...`
- [ ] Spec change was already merged in #
- [ ] New ADR added: `docs/adr/NNNN-....md`
- [ ] No ADR needed

## Testing

- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] New patterns include ≥5 match and ≥5 no-match test cases
- [ ] Manually verified on (check all that apply):
  - [ ] macOS
  - [ ] Linux
  - [ ] Windows

## Breaking changes

- [ ] No breaking changes
- [ ] Breaking change — migration notes documented in the spec update and release notes

## Checklist

- [ ] PR title follows the format `<area>: <short description>` (e.g., `patterns: add azure_storage_key`)
- [ ] Commits are focused — unrelated changes are split into separate PRs
- [ ] No new dependencies without an ADR justifying them
- [ ] Documentation (README, spec, help text) updated if user-visible behavior changed
