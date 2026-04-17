# 0004 — Multi-line pattern matching via begin/end marker buffering

- **Status:** accepted
- **Date:** 2026-04-17

## Context

The filter engine processes stdin line-by-line for streaming latency (spec 001: "Output MUST be flushed after each input line"). PEM private key blocks span 20–50 lines. Matching them requires seeing the entire block before deciding whether to mask.

Three approaches were considered (see Alternatives). The key constraint is preserving streaming latency for the 99%+ of lines that don't contain multi-line secrets.

## Decision

Add optional **begin/end marker buffering** to the filter engine:

1. The `Pattern` struct gains three optional fields: `Multiline bool`, `BeginMarker *regexp.Regexp`, `EndMarker *regexp.Regexp`.
2. When the engine encounters a line matching a multiline pattern's `BeginMarker`, it enters **buffering mode**: lines are accumulated instead of flushed.
3. Buffering continues until `EndMarker` is matched or a safety limit (100 lines or 64KB, whichever comes first) is exceeded.
4. On `EndMarker`: the buffered block is joined and the full `Regex` is applied. If matched, the block is replaced with `Replacement` (e.g. `[REDACTED PRIVATE KEY]`).
5. On safety limit exceeded: the buffered lines are flushed **unmodified** (fail-open). A warning is written to stderr.
6. Single-line patterns (the default) are unaffected — `Multiline` defaults to `false`.

## Consequences

**Positive:**

- Streaming latency is preserved for all single-line patterns (no buffering overhead)
- PEM private keys are fully masked with zero false positives (explicit markers)
- Fail-open on safety limit means malformed input never causes data loss
- The same mechanism can later support YAML secret blocks or multi-line JSON if needed

**Negative / trade-offs:**

- Lines inside a PEM block experience latency equal to the block's total read time (typically <10ms for a 2KB key, acceptable)
- If a `BEGIN` marker appears without a matching `END`, up to 100 lines are buffered before flushing — a brief delay for malformed input
- The Pattern struct gains 3 fields that are nil/false for most patterns — minor ergonomic cost

## Alternatives considered

- **Full-input accumulation.** Read all of stdin into memory, then match multi-line patterns. Rejected: breaks streaming guarantee (spec 001). A `tail -f | mask-pipe` pipeline would hang until EOF.

- **Sliding window.** Keep a rolling buffer of N lines, apply multi-line regex over the window. Rejected: PEM blocks vary widely in size (RSA-2048 = ~27 lines, RSA-4096 = ~52 lines). A fixed window either misses large keys or wastes memory on every line.

- **Defer to v2.** Skip `pem_private_key` entirely. Rejected: PEM keys have zero false positives (explicit markers) and are one of the most commonly leaked secret types. Omitting them weakens mask-pipe's value proposition.

## References

- [Issue #16](https://github.com/c12o-dev/mask-pipe/issues/16) — spec_change issue
- [Issue #8](https://github.com/c12o-dev/mask-pipe/issues/8) — pem_private_key pattern
- [docs/specs/001-cli-interface.md](../specs/001-cli-interface.md) — streaming latency requirement
- [docs/specs/002-pattern-library.md](../specs/002-pattern-library.md) — PEM pattern definition
