package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	for _, arg := range []string{"--version", "-V"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{arg}, strings.NewReader(""), &stdout, &stderr)
			if code != 0 {
				t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), "mask-pipe") {
				t.Errorf("stdout = %q, want to contain \"mask-pipe\"", stdout.String())
			}
		})
	}
}

func TestHelpFlag(t *testing.T) {
	for _, arg := range []string{"--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{arg}, strings.NewReader(""), &stdout, &stderr)
			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}
			if !strings.Contains(stdout.String(), "Usage") {
				t.Errorf("stdout missing Usage section: %q", stdout.String())
			}
		})
	}
}

func TestInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--bogus"}, strings.NewReader(""), &stdout, &stderr)
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if stderr.Len() == 0 {
		t.Error("stderr should contain an error message")
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should be empty on flag error, got %q", stdout.String())
	}
}

func TestCleanPassthrough(t *testing.T) {
	input := "hello\nworld\n"
	var stdout, stderr bytes.Buffer
	code := run(nil, strings.NewReader(input), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if stdout.String() != input {
		t.Errorf("stdout = %q, want %q", stdout.String(), input)
	}
}

func TestMaskingEndToEnd(t *testing.T) {
	input := "key is AKIAIOSFODNN7EXAMPLE\nno secret here\n"
	want := "key is AKIA************MPLE\nno secret here\n"
	var stdout, stderr bytes.Buffer
	code := run(nil, strings.NewReader(input), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestListPatternsSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"list-patterns"}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "aws_access_key") {
		t.Errorf("output missing aws_access_key: %q", out)
	}
	if !strings.Contains(out, "builtin") {
		t.Errorf("output missing 'builtin' source: %q", out)
	}
	if !strings.Contains(out, "enabled") {
		t.Errorf("output missing 'enabled' status: %q", out)
	}
}

func TestDoctorSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"doctor"}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "All checks passed") {
		t.Errorf("output missing 'All checks passed': %q", out)
	}
	if !strings.Contains(out, "8/8 enabled") {
		t.Errorf("output missing '8/8 enabled': %q", out)
	}
}

func TestUnknownSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"notreal"}, strings.NewReader(""), &stdout, &stderr)
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
}

func TestDryRunNoColor(t *testing.T) {
	input := "key AKIAIOSFODNN7EXAMPLE here\n"
	var stdout, stderr bytes.Buffer
	code := run([]string{"--dry-run", "--no-color"}, strings.NewReader(input), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "[MATCH:aws_access_key]") {
		t.Errorf("dry-run output missing MATCH tag: %q", out)
	}
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("dry-run should preserve original value: %q", out)
	}
}

func TestDryRunWithColor(t *testing.T) {
	input := "AKIAIOSFODNN7EXAMPLE\n"
	var stdout, stderr bytes.Buffer
	code := run([]string{"--dry-run"}, strings.NewReader(input), &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	out := stdout.String()
	if !strings.Contains(out, "\033[31m") {
		t.Errorf("dry-run with color should contain ANSI red: %q", out)
	}
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("dry-run should preserve original value: %q", out)
	}
}
