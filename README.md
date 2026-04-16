# mask-pipe

**Real-time secret masking for your terminal. Pipe it, mask it, done.**

```
$ env | mask-pipe

AWS_ACCESS_KEY_ID=AKIA****MPLE
AWS_SECRET_ACCESS_KEY=********************************
DATABASE_URL=postgres://admin:********@db.example.com:5432/mydb
GITHUB_TOKEN=ghp_****************************ef01
STRIPE_SK=sk_live_****************************7dc
```

TruffleHog guards your code. secretlint guards your commits. **mask-pipe guards your screen.**

> **Status:** pre-release. Specifications are being drafted. See [docs/specs/](docs/specs/) for the current design, and [CONTRIBUTING.md](CONTRIBUTING.md) for how to propose changes.

---

## Install

```bash
# Homebrew (coming soon)
brew install mask-pipe/tap/mask-pipe

# Go
go install github.com/c12o-dev/mask-pipe@latest

# Binary (Linux amd64)
curl -sSL https://github.com/c12o-dev/mask-pipe/releases/latest/download/mask-pipe_linux_amd64.tar.gz \
  | tar xz -C /usr/local/bin mask-pipe
```

Single binary with zero dependencies. Supports macOS, Linux, and Windows.

---

## Usage

```bash
# Mask environment variables before screen sharing
env | mask-pipe

# Tail Docker logs safely
docker logs -f my-app 2>&1 | mask-pipe

# Review .env files without exposing secrets
cat .env | mask-pipe

# Inspect Kubernetes pod logs
kubectl logs pod-name | mask-pipe

# Preview what would be masked (nothing is replaced)
terraform plan | mask-pipe --dry-run
```

---

## Built-in Patterns

mask-pipe ships with 8 high-precision patterns enabled by default:

| Pattern | Matches | Example |
|---|---|---|
| AWS Access Key | `AKIA[0-9A-Z]{16}` | `AKIAIOSFODNN7EXAMPLE` |
| AWS Secret Key | 40-char key after `AWS_SECRET_ACCESS_KEY=` | `wJalrXUtnFEMI/K7MDENG...` |
| GitHub Token | `ghp_`, `gho_`, `ghs_`, `github_pat_` | `ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZab...` |
| Stripe Key | `sk_live_`, `sk_test_`, `pk_live_`, `pk_test_` | `sk_live_4eC39Hq...7dc` |
| JWT | Three dot-separated base64 segments | `eyJhbGciOi...` |
| Database URL | Password in `://user:pass@host` | `postgres://admin:s3cret@...` |
| Password field | `password=`, `passwd=`, `secret=` values | `DB_PASSWORD=hunter2` |
| Private Key | PEM-encoded private key blocks | `-----BEGIN RSA PRIVATE KEY-----` |

All patterns are tuned for precision over recall — mask-pipe will not shred your normal output.

See [docs/specs/002-pattern-library.md](docs/specs/002-pattern-library.md) for the full specification.

---

## Configuration

Create `~/.mask-pipe.toml` to add custom patterns or tune behavior:

```toml
[builtin]
aws_access_key = true
github_token   = true
jwt            = false

[[custom]]
name    = "internal-api-key"
pattern = 'mycompany-key-[a-zA-Z0-9]{32}'

[display]
mask_char = "*"
show_tail = 4
color     = true
```

Full schema: [docs/specs/003-config-format.md](docs/specs/003-config-format.md).

---

## Shell Integration

The recommended usage is the **explicit pipe**: `command | mask-pipe`. This preserves TTY behavior (colors, TUI apps, interactive prompts).

For non-interactive commands, you can add shell function wrappers:

```zsh
mlogs()  { docker logs "$@"  | mask-pipe; return ${PIPESTATUS[0]}; }
mklogs() { kubectl logs "$@" | mask-pipe; return ${PIPESTATUS[0]}; }
menv()   { env                | mask-pipe; }
```

**Do not wrap interactive or TUI commands** (`docker run -it`, `vim`, `npm install`, etc.) — piping breaks them. See [docs/specs/004-shell-integration.md](docs/specs/004-shell-integration.md) for details.

---

## How It Compares

| | mask-pipe | TruffleHog | secretlint | GitHub `add-mask` |
|---|---|---|---|---|
| **What it protects** | Terminal output | Git repos & logs | Files & commits | CI logs |
| **When it runs** | Real-time (pipe) | Post-hoc scan | Pre-commit hook | CI runtime |
| **Works locally** | Yes | Yes | Yes | No (CI only) |
| **Built-in patterns** | 8 | 800+ | 30+ | Manual only |
| **Designed for** | Screen sharing, log tailing | Secret discovery | Commit prevention | CI log redaction |

These tools are **complementary, not competing.** A solid setup uses all of them.

---

## Contributing

This project uses a **hybrid spec-driven + issue-driven workflow**. Code changes are backed by a specification (in `docs/specs/`) and tracked via GitHub Issues. See [CONTRIBUTING.md](CONTRIBUTING.md) for the full workflow.

---

## License

MIT — see [LICENSE](LICENSE).
