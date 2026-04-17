package filter

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/c12o-dev/mask-pipe/patterns"
)

func BenchmarkMaskLineClean(b *testing.B) {
	f := New(patterns.Builtins, patterns.DefaultShowTail)
	line := "INFO 2026-04-17T10:32:14Z request processed in 34ms for tenant acme-corp"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.MaskLine(line)
	}
}

func BenchmarkMaskLineWithAWSKey(b *testing.B) {
	f := New(patterns.Builtins, patterns.DefaultShowTail)
	line := "config: aws_key=AKIAIOSFODNN7EXAMPLE region=us-east-1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.MaskLine(line)
	}
}

func BenchmarkMaskLineWithMultipleSecrets(b *testing.B) {
	f := New(patterns.Builtins, patterns.DefaultShowTail)
	line := `{"AccessKeyId":"AKIAIOSFODNN7EXAMPLE","Token":"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.sig"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.MaskLine(line)
	}
}

func BenchmarkRunStream1000Lines(b *testing.B) {
	f := &Filter{
		Patterns: patterns.Builtins,
		ShowTail: patterns.DefaultShowTail,
		MaskChar: "*",
	}

	var lines []string
	for i := 0; i < 990; i++ {
		lines = append(lines, "INFO normal log line without any secrets at all")
	}
	for i := 0; i < 10; i++ {
		lines = append(lines, "WARN leaked key=AKIAIOSFODNN7EXAMPLE in request")
	}
	input := strings.Join(lines, "\n") + "\n"

	b.ResetTimer()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		f.Run(strings.NewReader(input), io.Discard)
	}
}

func BenchmarkRunStream10000LinesClean(b *testing.B) {
	f := &Filter{
		Patterns: patterns.Builtins,
		ShowTail: patterns.DefaultShowTail,
		MaskChar: "*",
	}

	line := "INFO 2026-04-17T10:32:14Z request processed in 34ms for tenant acme-corp\n"
	input := strings.Repeat(line, 10000)

	b.ResetTimer()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		f.Run(strings.NewReader(input), io.Discard)
	}
}

// Prevent compiler from optimizing away results
var benchSink string
var benchBuf bytes.Buffer
