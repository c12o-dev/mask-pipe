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
