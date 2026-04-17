package filter

import (
	"bufio"
	"io"
	"strings"

	"github.com/c12o-dev/mask-pipe/patterns"
)

// Filter reads lines from a reader, masks secret patterns, and writes to a writer.
type Filter struct {
	Patterns []*patterns.Pattern
	ShowTail int
}

func New(pats []*patterns.Pattern, showTail int) *Filter {
	return &Filter{Patterns: pats, ShowTail: showTail}
}

func (f *Filter) Run(in io.Reader, out io.Writer) error {
	r := bufio.NewReaderSize(in, 1024*1024)
	w := bufio.NewWriter(out)
	defer w.Flush()
	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			hasNewline := strings.HasSuffix(line, "\n")
			content := strings.TrimSuffix(line, "\n")
			masked := f.MaskLine(content)
			w.WriteString(masked)
			if hasNewline {
				w.WriteByte('\n')
			}
			w.Flush()
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (f *Filter) MaskLine(line string) string {
	for _, p := range f.Patterns {
		line = f.applyPattern(line, p)
	}
	return line
}

func (f *Filter) applyPattern(line string, p *patterns.Pattern) string {
	indices := p.Regex.FindAllStringSubmatchIndex(line, -1)
	if len(indices) == 0 {
		return line
	}
	var b strings.Builder
	prev := 0
	idx := p.CaptureIdx * 2
	for _, loc := range indices {
		start := loc[idx]
		end := loc[idx+1]
		if start < 0 {
			continue
		}
		b.WriteString(line[prev:start])
		b.WriteString(f.mask(p, line[start:end]))
		prev = end
	}
	b.WriteString(line[prev:])
	return b.String()
}

func (f *Filter) mask(p *patterns.Pattern, value string) string {
	if p.Replacement != "" {
		return p.Replacement
	}
	return patterns.DefaultMask(value, f.ShowTail)
}
