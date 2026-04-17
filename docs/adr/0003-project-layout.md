# 0003 — Project layout: `cmd/` + top-level `patterns/` + `internal/`

- **Status:** accepted
- **Date:** 2026-04-16

## Context

mask-pipe is about to grow from an empty repo into a real Go project. The first PR (Issue #1) creates the skeleton, so we need to pick a directory layout before code lands. Three realistic options:

1. **Flat** — `main.go` + sibling packages at the repo root.
2. **`cmd/mask-pipe/`** — standard Go layout with the entry point under `cmd/` and library packages at the root or under `internal/`.
3. **Full `cmd/` + `pkg/` + `internal/`** — enterprise-style split popular in the Kubernetes ecosystem.

We also need to decide the placement of the `patterns/` package. [Spec 002](../specs/002-pattern-library.md) and [CLAUDE.md](../../CLAUDE.md) both refer to `patterns/` at the top level; a contributor reading either document will expect to find it there.

## Decision

Use the following layout:

```
mask-pipe/
├── cmd/
│   └── mask-pipe/
│       └── main.go           # entry point, flag parsing, wiring
├── patterns/                 # built-in pattern registry (top-level, per spec 002)
│   ├── pattern.go            # Pattern struct + registry API
│   ├── aws_access_key.go     # one file per pattern
│   └── patterns_test.go
├── internal/
│   └── filter/               # line-by-line stdin→stdout filter engine
│       ├── filter.go
│       └── filter_test.go
├── go.mod
└── go.sum
```

Rules:

- **Single binary for now**, but use `cmd/mask-pipe/` anyway so adding a second binary later (e.g., `mask-pipe-benchmark`) costs zero refactor.
- **`patterns/` stays top-level** (not under `internal/`) because spec 002 / CLAUDE.md document it at that path and because the pattern registry is a coherent unit we want findable at a glance. We are NOT promising it as a stable external API — the module path is `github.com/c12o-dev/mask-pipe/patterns` but importers get no semver guarantees until we say so.
- **Everything else goes under `internal/`** by default. A package only leaves `internal/` when we have a specific reason to expose it.
- **No `pkg/` directory.** It's a Go community anti-pattern that adds a directory level without semantic meaning.

## Consequences

**Positive:**

- Contributors find `patterns/` where both the spec and CLAUDE.md say it will be — no surprise.
- `cmd/mask-pipe/` leaves room for a second binary without a painful restructure.
- `internal/` prevents accidental external dependence on implementation details (the compiler enforces it).
- Matches conventions in well-regarded Go CLIs (`gh`, `goreleaser`, `hugo`), lowering onboarding cost for experienced Go contributors.

**Negative / trade-offs:**

- One extra directory level (`cmd/mask-pipe/main.go`) compared to flat layout. Cost is a single `go build ./cmd/mask-pipe` instead of `go build .`. Acceptable.
- `patterns/` being top-level weakly implies API stability. Mitigated by a `doc.go` comment explicitly stating it is not a stable external API in v1.
- If we later decide to publish `patterns/` as a reusable library, we've already chosen the import path — we can't rename without a major version bump.

## Alternatives considered

- **Flat layout.** Simplest for the first 200 lines of code, but the moment we add a second binary or start isolating internals it forces a rename. Picking the destination layout now costs nothing.

- **Full `cmd/` + `pkg/` + `internal/`.** `pkg/` adds a directory level with no semantic meaning; Dave Cheney and the Go team have both argued against it. Skipped.

- **Put `patterns/` under `internal/patterns/`.** Safer in terms of API commitment, but contradicts the documented path in spec 002 and CLAUDE.md. We'd have to update both. The cost of leaving it top-level is only the implicit-API concern, which we mitigate with a `doc.go` note.

## References

- [docs/specs/002-pattern-library.md](../specs/002-pattern-library.md) — references `patterns/` directly
- [Issue #1](https://github.com/c12o-dev/mask-pipe/issues/1) — the MVP vertical slice that this ADR unblocks
- [Standard Go Project Layout (golang-standards)](https://github.com/golang-standards/project-layout) — non-official but widely followed reference
- Dave Cheney on `pkg/`: https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project
