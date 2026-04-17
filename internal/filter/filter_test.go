package filter

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/c12o-dev/mask-pipe/patterns"
)

func testFilter() *Filter {
	return New(patterns.Builtins, patterns.DefaultShowTail)
}

func TestMaskLineAWSAccessKey(t *testing.T) {
	f := testFilter()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare key", "AKIAIOSFODNN7EXAMPLE", "AKIA************MPLE"},
		{"key in text", "found AKIAIOSFODNN7EXAMPLE in logs", "found AKIA************MPLE in logs"},
		{"multiple keys", "key1=AKIAIOSFODNN7EXAMPLE key2=AKIA2RAY6KGQJM7Q3EYX", "key1=AKIA************MPLE key2=AKIA************3EYX"},
		{"no match", "no secrets here", "no secrets here"},
		{"json context", `{"AccessKeyId":"AKIAIOSFODNN7EXAMPLE"}`, `{"AccessKeyId":"AKIA************MPLE"}`},
		{"empty line", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.MaskLine(tt.input)
			if got != tt.want {
				t.Errorf("MaskLine(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRunPreservesStreamStructure(t *testing.T) {
	f := testFilter()
	input := "line1 AKIAIOSFODNN7EXAMPLE\nline2 clean\nline3 AKIA2RAY6KGQJM7Q3EYX\n"
	want := "line1 AKIA************MPLE\nline2 clean\nline3 AKIA************3EYX\n"
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != want {
		t.Errorf("Run output = %q, want %q", out.String(), want)
	}
}

func TestRunPreservesNoTrailingNewline(t *testing.T) {
	f := testFilter()
	input := "AKIAIOSFODNN7EXAMPLE"
	want := "AKIA************MPLE"
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != want {
		t.Errorf("Run output = %q, want %q", out.String(), want)
	}
}

func TestRunEmptyInput(t *testing.T) {
	f := testFilter()
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(""), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != "" {
		t.Errorf("Run output = %q, want empty", out.String())
	}
}

func TestCaptureGroupMasking(t *testing.T) {
	p := &patterns.Pattern{
		ID:         "test_capture",
		Name:       "Test Capture",
		Regex:      regexp.MustCompile(`secret=([A-Za-z0-9]+)`),
		CaptureIdx: 1,
	}
	f := New([]*patterns.Pattern{p}, 4)
	got := f.MaskLine("config: secret=MyS3cretValue123 done")
	want := "config: secret=MyS3********e123 done"
	if got != want {
		t.Errorf("MaskLine = %q, want %q", got, want)
	}
}

func TestLiteralReplacement(t *testing.T) {
	p := &patterns.Pattern{
		ID:          "test_literal",
		Name:        "Test Literal",
		Regex:       regexp.MustCompile(`SECRET_[A-Z]+`),
		CaptureIdx:  0,
		Replacement: "****",
	}
	f := New([]*patterns.Pattern{p}, 4)
	got := f.MaskLine("value=SECRET_TOKEN")
	want := "value=****"
	if got != want {
		t.Errorf("MaskLine = %q, want %q", got, want)
	}
}

func TestMultilinePEMBlock(t *testing.T) {
	f := testFilter()
	input := "before\n-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJBALRi\nbase64data==\n-----END RSA PRIVATE KEY-----\nafter\n"
	want := "before\n[REDACTED PRIVATE KEY]\nafter\n"
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != want {
		t.Errorf("Run output = %q, want %q", out.String(), want)
	}
}

func TestMultilinePEMWithSurroundingSecrets(t *testing.T) {
	f := testFilter()
	input := "key=AKIAIOSFODNN7EXAMPLE\n-----BEGIN PRIVATE KEY-----\ndata\n-----END PRIVATE KEY-----\nclean\n"
	want := "key=AKIA************MPLE\n[REDACTED PRIVATE KEY]\nclean\n"
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != want {
		t.Errorf("Run output = %q, want %q", out.String(), want)
	}
}

func TestMultilineBeginWithoutEnd(t *testing.T) {
	f := testFilter()
	f.Stderr = io.Discard
	// Begin marker without end — should flush unmodified at EOF
	input := "-----BEGIN RSA PRIVATE KEY-----\norphaned data\n"
	var out bytes.Buffer
	if err := f.Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.String() != input {
		t.Errorf("Run output = %q, want %q (unmodified)", out.String(), input)
	}
}
