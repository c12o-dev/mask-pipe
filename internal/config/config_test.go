package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadDefault(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Display.ShowTail != 4 {
		t.Errorf("ShowTail = %d, want 4", cfg.Display.ShowTail)
	}
	if cfg.Display.MaskChar != "*" {
		t.Errorf("MaskChar = %q, want *", cfg.Display.MaskChar)
	}
}

func TestLoadBuiltinToggle(t *testing.T) {
	path := writeTempConfig(t, `
[builtin]
jwt = false
aws_access_key = true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.IsBuiltinEnabled("jwt") {
		t.Error("jwt should be disabled")
	}
	if !cfg.IsBuiltinEnabled("aws_access_key") {
		t.Error("aws_access_key should be enabled")
	}
	if !cfg.IsBuiltinEnabled("stripe_key") {
		t.Error("stripe_key should default to enabled")
	}
}

func TestLoadDisplayOverrides(t *testing.T) {
	path := writeTempConfig(t, `
[display]
mask_char = "#"
show_tail = 2
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Display.MaskChar != "#" {
		t.Errorf("MaskChar = %q, want #", cfg.Display.MaskChar)
	}
	if cfg.Display.ShowTail != 2 {
		t.Errorf("ShowTail = %d, want 2", cfg.Display.ShowTail)
	}
}

func TestLoadCustomPatterns(t *testing.T) {
	path := writeTempConfig(t, `
[[custom]]
name = "my-key"
pattern = 'myco-[a-z]{32}'
show_tail = 0

[[custom]]
name = "webhook"
pattern = 'https://hooks\.example\.com/[a-z]+'
replacement = "[REDACTED URL]"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Custom) != 2 {
		t.Fatalf("Custom count = %d, want 2", len(cfg.Custom))
	}
	compiled := cfg.CustomPatterns()
	if len(compiled) != 2 {
		t.Fatalf("Compiled count = %d, want 2", len(compiled))
	}
	if compiled[0].ShowTail != 0 {
		t.Errorf("first custom ShowTail = %d, want 0", compiled[0].ShowTail)
	}
	if compiled[1].Replacement != "[REDACTED URL]" {
		t.Errorf("second custom Replacement = %q", compiled[1].Replacement)
	}
}

func TestLoadAllowlist(t *testing.T) {
	path := writeTempConfig(t, `
[[allowlist]]
name = "test-key"
pattern = 'AKIAIOSFODNN7EXAMPLE'
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	regexps := cfg.AllowlistRegexps()
	if len(regexps) != 1 {
		t.Fatalf("Allowlist count = %d, want 1", len(regexps))
	}
	if !regexps[0].MatchString("AKIAIOSFODNN7EXAMPLE") {
		t.Error("allowlist should match the example key")
	}
}

func TestInvalidRegexFails(t *testing.T) {
	path := writeTempConfig(t, `
[[custom]]
name = "bad"
pattern = '[invalid'
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestDuplicateNameFails(t *testing.T) {
	path := writeTempConfig(t, `
[[custom]]
name = "dup"
pattern = 'a'

[[custom]]
name = "dup"
pattern = 'b'
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestInvalidMaskCharFails(t *testing.T) {
	path := writeTempConfig(t, `
[display]
mask_char = "**"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for multi-char mask_char")
	}
}

func TestExplicitConfigNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.toml")
	if err == nil {
		t.Fatal("expected error for nonexistent explicit config")
	}
}
