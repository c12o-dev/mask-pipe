package filter

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/c12o-dev/mask-pipe/patterns"
)

const (
	maxBufLines = 100
	maxBufBytes = 64 * 1024
)

// Filter reads lines from a reader, masks secret patterns, and writes to a writer.
type Filter struct {
	Patterns  []*patterns.Pattern
	ShowTail  int
	MaskChar  string
	Allowlist []*regexp.Regexp
	Stderr    io.Writer // optional; multi-line safety-limit warnings go here
}

func New(pats []*patterns.Pattern, showTail int) *Filter {
	return &Filter{Patterns: pats, ShowTail: showTail}
}

func (f *Filter) Run(in io.Reader, out io.Writer) error {
	r := bufio.NewReaderSize(in, 1024*1024)
	w := bufio.NewWriter(out)
	defer w.Flush()

	var mlRaw []string // multi-line buffer (raw lines with \n preserved)
	var mlPat *patterns.Pattern
	var mlBytes int

	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			hasNewline := strings.HasSuffix(line, "\n")
			content := strings.TrimSuffix(line, "\n")

			if mlPat != nil {
				// In buffering mode — store raw line
				mlRaw = append(mlRaw, line)
				mlBytes += len(line)

				if mlPat.EndMarker.MatchString(content) {
					// End marker found — join contents, apply regex
					contents := make([]string, len(mlRaw))
					for i, l := range mlRaw {
						contents[i] = strings.TrimSuffix(l, "\n")
					}
					block := strings.Join(contents, "\n")
					masked := f.applyPattern(block, mlPat)
					w.WriteString(masked)
					if hasNewline {
						w.WriteByte('\n')
					}
					w.Flush()
					mlRaw = nil
					mlPat = nil
					mlBytes = 0
				} else if len(mlRaw) >= maxBufLines || mlBytes >= maxBufBytes {
					// Safety limit — flush unmodified
					if f.Stderr != nil {
						fmt.Fprintf(f.Stderr, "mask-pipe: multi-line buffer limit exceeded, flushing %d lines unmasked\n", len(mlRaw))
					}
					for _, rawLine := range mlRaw {
						w.WriteString(rawLine)
					}
					w.Flush()
					mlRaw = nil
					mlPat = nil
					mlBytes = 0
				}
				// Otherwise keep buffering
			} else if p := f.matchBeginMarker(content); p != nil {
				// Start buffering for a multi-line pattern
				mlRaw = []string{line}
				mlPat = p
				mlBytes = len(line)
			} else {
				// Normal single-line processing
				masked := f.MaskLine(content)
				w.WriteString(masked)
				if hasNewline {
					w.WriteByte('\n')
				}
				w.Flush()
			}
		}
		if err == io.EOF {
			// Flush any remaining multi-line buffer unmodified
			if mlPat != nil {
				for _, rawLine := range mlRaw {
					w.WriteString(rawLine)
				}
			}
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (f *Filter) matchBeginMarker(line string) *patterns.Pattern {
	for _, p := range f.Patterns {
		if p.Multiline && p.BeginMarker != nil && p.BeginMarker.MatchString(line) {
			return p
		}
	}
	return nil
}

func (f *Filter) MaskLine(line string) string {
	for _, p := range f.Patterns {
		if p.Multiline {
			continue
		}
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
		matched := line[start:end]
		if f.isAllowlisted(matched) {
			continue
		}
		b.WriteString(line[prev:start])
		b.WriteString(f.mask(p, matched))
		prev = end
	}
	b.WriteString(line[prev:])
	return b.String()
}

func (f *Filter) isAllowlisted(value string) bool {
	for _, re := range f.Allowlist {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}

func (f *Filter) mask(p *patterns.Pattern, value string) string {
	if p.Replacement != "" {
		return p.Replacement
	}
	return patterns.DefaultMask(value, f.ShowTail, f.MaskChar)
}
