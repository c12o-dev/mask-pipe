# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) — short documents capturing a significant technical decision, the context that drove it, and its consequences.

ADRs are append-only. Don't edit accepted ADRs; if a decision changes, write a new ADR that supersedes the old one.

## Index

| # | Title | Status |
|---|---|---|
| [0001](0001-language-go.md) | Implementation language: Go | accepted |
| [0002](0002-explicit-pipe-default.md) | Explicit pipe as default integration pattern | accepted |
| [0003](0003-project-layout.md) | Project layout: `cmd/` + top-level `patterns/` + `internal/` | accepted |
| [0004](0004-multiline-buffering.md) | Multi-line pattern matching via begin/end marker buffering | accepted |

## When to write an ADR

Write an ADR when the decision:

- Is hard to reverse (language choice, major dependencies, public API shape)
- Affects multiple parts of the system (architectural patterns)
- Is likely to be questioned later ("why did we pick X over Y?")

Skip an ADR for:

- Obvious implementation details
- Temporary workarounds
- Personal code style preferences

## Format

Use [`0000-template.md`](0000-template.md) as the starting point. Keep ADRs short — 1–2 pages. If you need more, you probably need a spec in `docs/specs/` instead.

## Numbering

ADRs are numbered sequentially: `0001`, `0002`, `0003`. Don't skip numbers. If two PRs race and both claim the same number, resolve at merge time.
