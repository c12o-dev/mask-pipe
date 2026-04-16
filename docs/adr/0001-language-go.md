# 0001 — Implementation language: Go

- **Status:** accepted
- **Date:** 2026-04-16

## Context

mask-pipe needs to be distributed as a single binary with zero runtime dependencies, runnable on macOS / Linux / Windows for both amd64 and arm64. It is a small CLI whose hot path is regex matching over stdin.

Previous competitor analysis showed that Teller (Rust) has a `libssl.so.1.1` dynamic-link bug that breaks it on Ubuntu 22.04+. We must not repeat that class of distribution failure.

## Decision

Implement mask-pipe in **Go**.

## Consequences

**Positive:**

- `go build` produces a statically-linked single binary for all supported platforms
- `goreleaser` cross-compiles and publishes Homebrew / Scoop / APT assets with one command
- Go's standard `regexp` package uses RE2, which gives us O(n) linear-time regex guarantees — critical for predictable streaming latency
- Ecosystem includes mature PTY libraries (`creack/pty`) for the eventual v2 shell-wrap feature
- Large contributor pool; Go is familiar to most CLI tool authors

**Negative / trade-offs:**

- Binary size ~5–8MB unstripped (vs. ~2MB for Rust with aggressive optimization). Mitigated by UPX compression in release builds if size becomes an issue.
- Go's `regexp` is slower than Rust's `regex` crate for complex patterns — accepted for v1; re-evaluate only if benchmarks show it matters for realistic pattern counts
- Less aggressive optimization than Rust — accepted; mask-pipe is not CPU-bound for its intended workload

## Alternatives considered

- **Rust** — Excellent performance and the `regex` crate is top-tier, but:
  - Cross-compilation is harder (OpenSSL / libc linkage issues, see Teller case study)
  - Longer build times slow iteration on a small team
  - Fewer contributors than Go in the CLI-tool space

- **Zig** — Produces small static binaries but the ecosystem is too young for a tool we want running on production systems by 2026

- **Python / Node.js** — Rejected immediately: requiring a runtime defeats the "install and forget" UX we're aiming for

- **C** — Rejected for memory-safety reasons in a tool that handles potentially-sensitive data

## References

- [docs/specs/000-vision.md](../specs/000-vision.md) — single binary requirement
- Teller `libssl.so.1.1` bug: https://github.com/tellerops/teller/issues/290
