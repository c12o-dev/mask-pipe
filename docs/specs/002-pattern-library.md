# 002 — Built-in Pattern Library

- **Status:** draft
- **Last updated:** 2026-04-16

## Summary

Defines the built-in regex patterns shipped with mask-pipe, the precision requirements each pattern must meet, and the process for adding new patterns.

## Specification

### Design principle: precision over recall

mask-pipe's reputation depends on not destroying normal output. A false positive (masking a non-secret string) is worse than a false negative (missing a real secret). Built-in patterns MUST target **precision ≥99% on realistic input**.

Patterns that cannot meet this bar belong behind an opt-in flag or in a separate "strict" pattern pack, not in the default set.

### Initial pattern set (v1)

| ID | Name | Regex | Rationale |
|---|---|---|---|
| `aws_access_key` | AWS Access Key ID | `AKIA[0-9A-Z]{16}` | 20-char prefix is unique to AWS; virtually no false positives |
| `aws_secret_key` | AWS Secret Access Key | `(?i)aws_secret_access_key\s*[=:]\s*([A-Za-z0-9/+=]{40})` | Context-dependent match (key=value form) to avoid false positives on arbitrary 40-char strings |
| `github_token` | GitHub Token | `\bgh[pousr]_[A-Za-z0-9]{36,}` | GitHub's documented token prefix format; `\b` prevents mid-word FP; body is base62 (no `_`) |
| `github_pat` | GitHub Fine-grained PAT | `\bgithub_pat_[A-Za-z0-9_]{80,}` | Longer prefix for fine-grained PATs; `\b` prevents mid-word FP |
| `stripe_key` | Stripe API Key | `\b[sp]k_(?:live\|test)_[A-Za-z0-9]{24,}` | Stripe's documented prefix scheme; `\b` prevents `disk_test_...` FP class |
| `jwt` | JSON Web Token | `eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+` | Three dot-separated base64 segments starting with JWT header |
| `db_url_password` | DB URL password | `://[^:/\s]+:([^@/\s]+)@` | Captures the password portion only; `/` excluded from capture to avoid `host:port/path@` FP |
| `pem_private_key` | PEM private key block | `-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----` | Multi-line; explicit markers mean no false positives |

### Pattern metadata (required fields)

Each pattern in the codebase MUST declare:

```go
type Pattern struct {
    ID          string         // Stable identifier (e.g., "aws_access_key")
    Name        string         // Human-readable name
    Regex       *regexp.Regexp
    CaptureIdx  int            // Which capture group holds the value to mask (0 = entire match)
    Replacement string         // "****" or custom per-pattern
    Examples    []string       // At least 5 strings that SHOULD match
    NonExamples []string       // At least 5 strings that SHOULD NOT match
    Source      string         // URL to vendor doc or spec justifying the format
    Multiline   bool           // If true, engine uses begin/end marker buffering (ADR 0004)
    BeginMarker *regexp.Regexp // Line that starts a multi-line block (required if Multiline)
    EndMarker   *regexp.Regexp // Line that ends a multi-line block (required if Multiline)
}
```

### Testing requirements

Every pattern MUST have:

- At least 5 positive test cases (`Examples`) based on real-world data (synthetic but realistic)
- At least 5 negative test cases (`NonExamples`) specifically targeting common false-positive scenarios
- A test that runs the pattern against a corpus of "realistic non-secret text" (e.g., code snippets, log lines) and asserts zero matches

Patterns failing any of these tests MUST NOT ship in the default set.

### Mask replacement format

By default, a matched value is replaced with:

```
<first 4 chars><N-8 asterisks><last 4 chars>
```

Example: `AKIAIOSFODNN7EXAMPLE` → `AKIA************MPLE`

Exceptions (full masking):

- Passwords in DB URLs: fully masked (`****`), because first/last chars leak entropy
- Private key blocks: replaced with a single `[REDACTED PRIVATE KEY]` line

The `show_tail` config option controls the number of trailing characters shown (default `4`, `0` = full mask).

### Adding a new pattern

See [CONTRIBUTING.md — Proposing a new built-in pattern](../../CONTRIBUTING.md#proposing-a-new-built-in-pattern). Process:

1. Open a `pattern_proposal` issue with the regex, match/no-match examples, and source reference
2. Open a spec PR updating this file with the pattern added to the table above
3. Open an implementation PR adding the pattern to `patterns/` with tests

### Patterns deliberately excluded from v1

| Pattern | Reason |
|---|---|
| Credit card numbers | High false-positive risk on arbitrary 16-digit strings; PCI compliance is out of scope |
| Generic "password" field | Too broad; moved behind opt-in `generic_password` pattern |
| Email addresses | Not a secret; belongs in a separate PII-focused tool |
| IP addresses | Not a secret; noise |
| API keys without a documented prefix | No way to distinguish from random strings without entropy analysis |

Entropy-based detection (as Gitleaks does) is explicitly rejected for v1 due to false-positive issues documented in competitor analysis.

## Non-goals

- Exhaustive coverage of every secret format (that's TruffleHog's job — 800+ detectors)
- Service-specific verification (calling the API to check if a token is active)
- Historical pattern changes tracking (Git history is sufficient)

## Open questions

- Should we adopt Gitleaks' pattern definitions wholesale as an optional import? Would inherit both their coverage AND their false-positive issues.
- How do we handle patterns that require multi-line matching (e.g., YAML secrets blocks)? Current line-by-line engine is insufficient.
- Should custom patterns be verified against the same 5-match-5-no-match requirement at config load time? Would improve user-defined pattern quality.
