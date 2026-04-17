package patterns

import (
	"testing"
)

func TestBuiltinsHaveMinimumExamples(t *testing.T) {
	for _, p := range Builtins {
		if len(p.Examples) < 5 {
			t.Errorf("%s: only %d examples, spec requires >=5", p.ID, len(p.Examples))
		}
		if len(p.NonExamples) < 5 {
			t.Errorf("%s: only %d non-examples, spec requires >=5", p.ID, len(p.NonExamples))
		}
	}
}

func TestBuiltinsMatchExamples(t *testing.T) {
	for _, p := range Builtins {
		t.Run(p.ID, func(t *testing.T) {
			for _, ex := range p.Examples {
				if !p.Regex.MatchString(ex) {
					t.Errorf("expected match on %q", ex)
				}
			}
		})
	}
}

func TestBuiltinsRejectNonExamples(t *testing.T) {
	for _, p := range Builtins {
		t.Run(p.ID, func(t *testing.T) {
			for _, ne := range p.NonExamples {
				if p.Regex.MatchString(ne) {
					t.Errorf("unexpected match on non-example %q (false positive)", ne)
				}
			}
		})
	}
}

const realisticCorpus = `package main

import "fmt"

func main() {
    // Reference: AWS IAM policies can be attached to users, groups, or roles.
    // The AKIA prefix is reserved; STS returns temporary credentials instead.
    fmt.Println("Hello, world!")
}

[INFO] 2026-04-17T10:32:14Z request processed in 34ms
[WARN] 2026-04-17T10:32:15Z retry 1 of 3 for tenant AAAA-BBBB-CCCC-DDDD-EEEE
[ERROR] 2026-04-17T10:32:16Z failed: connection refused to 10.0.0.42:5432
[DEBUG] 2026-04-17T10:32:17Z cache hit ratio 0.98, RSS 12345678 bytes

$ ls -la
-rw-r--r--  1 user group    1234 Apr 17 10:32 README.md
-rw-r--r--  1 user group   56789 Apr 17 10:32 go.mod

commit 1a2b3c4d5e6f7890abcdef1234567890abcdef12
Author: Alice <alice@example.com>
Date:   Thu Apr 17 10:32:00 2026 +0000
    docs: update README with AKIA prefix guidance

PATH=/usr/local/bin:/usr/bin:/bin HOME=/home/user SHELL=/bin/bash
AWS_REGION=us-east-1 AWS_DEFAULT_OUTPUT=json

resource "aws_iam_access_key" "deploy" {
  user = aws_iam_user.deploy.name
}
# Note: the AKIA prefix indicates an AWS access key ID.
# See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
export AWS_ACCESS_KEY_ID=   # placeholder, set by CI

# GitHub Actions CI
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  GH_PAT: ${{ secrets.DEPLOY_PAT }}
run: |
  gh auth status
  curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user

# Stripe integration docs
# Keys look like sk_live_... or pk_test_... (see https://docs.stripe.com/keys)
# In test mode, use the dashboard to rotate keys.
STRIPE_PUBLISHABLE_KEY=   # set in .env

# Database configuration
DATABASE_URL=postgres://localhost:5432/myapp   # no credentials in dev
REDIS_URL=redis://localhost:6379

# JWT documentation
# A JWT has three parts: header.payload.signature
# The header starts with eyJ (base64 for {"...)
# Example structure: {"alg":"HS256","typ":"JWT"}

# Disk operations (should not trigger stripe pattern)
disk_test_partition_cleanup completed in 42s
task_live_migration_batch_12345 started

# Various URLs without credentials
https://example.com:8080/api/v1/@latest
git@github.com:org/repo.git
mailto:admin@example.com
ftp://mirror.example.com/pub/release.tar.gz
`

func TestBuiltinsZeroMatchesOnCorpus(t *testing.T) {
	for _, p := range Builtins {
		t.Run(p.ID, func(t *testing.T) {
			if loc := p.Regex.FindStringIndex(realisticCorpus); loc != nil {
				matched := realisticCorpus[loc[0]:loc[1]]
				t.Errorf("unexpected match on corpus: %q", matched)
			}
		})
	}
}

func TestDefaultMask(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		showTail int
		want     string
	}{
		{"standard", "AKIAIOSFODNN7EXAMPLE", 4, "AKIA************MPLE"},
		{"full-mask", "AKIAIOSFODNN7EXAMPLE", 0, "********************"},
		{"short-value", "AKIA", 4, "****"},
		{"boundary", "ABCDEFGH", 4, "********"},
		{"negative-tail", "AKIAIOSFODNN7EXAMPLE", -1, "********************"},
		{"nine-chars", "ABCDEFGHI", 4, "ABCD*FGHI"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultMask(tt.value, tt.showTail, "*")
			if got != tt.want {
				t.Errorf("DefaultMask(%q, %d) = %q, want %q", tt.value, tt.showTail, got, tt.want)
			}
		})
	}
}
