# 003 — Configuration File Format

- **Status:** draft
- **Last updated:** 2026-04-16

## Summary

Defines the TOML configuration file format for mask-pipe, including schema, default location, and loading semantics.

## Specification

### File location

mask-pipe searches for a config file in this order, using the first match:

1. Path from `--config <path>` flag (explicit)
2. `$MASK_PIPE_CONFIG` environment variable
3. `$XDG_CONFIG_HOME/mask-pipe/config.toml` (typically `~/.config/mask-pipe/config.toml`)
4. `~/.mask-pipe.toml` (legacy / convenient location)

If none exist, mask-pipe runs with built-in defaults — this is NOT an error.

### Schema

```toml
# ~/.mask-pipe.toml

# Toggle individual built-in patterns. All enabled by default.
# Pattern IDs are stable; see docs/specs/002-pattern-library.md.
[builtin]
aws_access_key   = true
aws_secret_key   = true
github_token     = true
github_pat       = true
stripe_key       = true
jwt              = true
db_url_password  = true
pem_private_key  = true

# Custom patterns. Each [[custom]] entry adds one pattern.
[[custom]]
name    = "internal-api-key"
pattern = 'mycompany-key-[a-zA-Z0-9]{32}'
# Optional: fully mask, or show N trailing chars (overrides global show_tail)
show_tail   = 0
replacement = ""   # If set, replaces the entire match with this string

[[custom]]
name    = "slack-webhook"
pattern = 'https://hooks\.slack\.com/services/T[A-Z0-9]+/B[A-Z0-9]+/[a-zA-Z0-9]+'

# Display settings apply globally.
[display]
mask_char = "*"    # Character used for masking
show_tail = 4      # Number of trailing chars to preserve (0 = full mask)
color     = true   # Highlight masked regions with ANSI red (only when stdout is a TTY)

# Allowlist: strings matching these regexes are NEVER masked, even if patterns match.
# Useful for known-safe test values, e.g., AWS example keys in docs.
[[allowlist]]
name    = "aws-doc-example"
pattern = 'AKIAIOSFODNN7EXAMPLE'

[[allowlist]]
name    = "jwt-docs-example"
pattern = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\.[^.]+\.[^.\s]+'
```

### Field reference

#### `[builtin]`

A table of `<pattern_id> = <bool>`. Keys correspond to pattern IDs in [002-pattern-library.md](002-pattern-library.md). Default: all `true`.

#### `[[custom]]`

Array of tables defining custom patterns. Each entry:

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | Yes | Unique identifier for the pattern |
| `pattern` | string | Yes | Go RE2-compatible regex |
| `show_tail` | int | No | Per-pattern trailing character count. Overrides `[display]`. |
| `replacement` | string | No | Literal replacement string. If set, overrides the default masking behavior. |

Regex compilation happens at config load. Invalid regexes cause mask-pipe to fail with exit code `1` and a specific error message pointing to the offending `name`.

#### `[display]`

Global display settings:

| Field | Type | Default | Description |
|---|---|---|---|
| `mask_char` | string | `"*"` | Character used in masks. Must be a single Unicode character. |
| `show_tail` | int | `4` | Number of trailing chars preserved. `0` = full mask. |
| `color` | bool | `true` | Emit ANSI color codes when stdout is a TTY. `false` disables colors globally. |

#### `[[allowlist]]`

Array of regex patterns. If an input match also matches ANY allowlist pattern, it is NOT masked. Useful for documented example credentials.

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | Yes | Unique identifier for the allowlist entry |
| `pattern` | string | Yes | Go RE2-compatible regex |

### Loading semantics

- Configs are loaded once at startup. No hot reload.
- Unknown top-level tables or fields are silently ignored (forward compatibility)
- Invalid TOML syntax → exit `1` with a parse error pointing to the line number
- Invalid regex in `[[custom]]` or `[[allowlist]]` → exit `1` with the pattern name and regex error
- Duplicate pattern names (across `[[custom]]` entries) → exit `1` with both line numbers

### Precedence order (at match time)

For each line of input:

1. Find all matches from active built-in patterns and all `[[custom]]` patterns
2. For each match, check against all `[[allowlist]]` entries. If any allowlist entry matches the same span, skip masking.
3. Apply masking to remaining matches, longest match first (to avoid overlapping rewrites)

### Environment variables

| Variable | Effect |
|---|---|
| `MASK_PIPE_CONFIG` | Override config file path |
| `NO_COLOR` | Disable ANSI colors (standard, overrides `[display].color`) |
| `MASK_PIPE_DEBUG` | Enable verbose diagnostic output to stderr |

## Non-goals

- YAML or JSON config formats (TOML only, for consistency)
- Hierarchical config merging (system-wide + user + project) — too much complexity for the value
- Config generation wizard (`mask-pipe init`) — v2+ consideration
- Remote config sources (URL-based configs) — security anti-pattern

## Open questions

- Should `[[custom]]` patterns support named capture groups for partial masking (e.g., preserving the domain but masking the credential)?
- Should there be a `[profile]` section allowing different config sets activated by `--profile` flag (e.g., `work` vs. `personal`)?
- How to handle config file discovery on Windows where `$XDG_CONFIG_HOME` is not standard?
