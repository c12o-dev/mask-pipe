package patterns

import (
	"regexp"
	"strings"
)

// Pattern describes a single secret-detection rule.
//
// Replacement: empty string applies DefaultMask (show head + stars + tail);
// non-empty string is used as a literal replacement (e.g. "****").
//
// Multiline patterns use begin/end marker buffering (ADR 0004). The filter
// engine accumulates lines between BeginMarker and EndMarker, then applies
// Regex to the joined block.
type Pattern struct {
	ID          string
	Name        string
	Regex       *regexp.Regexp
	CaptureIdx  int
	Replacement string
	Examples    []string
	NonExamples []string
	Source      string
	Multiline   bool
	BeginMarker *regexp.Regexp
	EndMarker   *regexp.Regexp
}

const DefaultShowTail = 4

// DefaultMask keeps the first 4 and last showTail characters visible.
// showTail <= 0 fully masks the value.
func DefaultMask(value string, showTail int) string {
	if showTail <= 0 {
		return strings.Repeat("*", len(value))
	}
	const head = 4
	if len(value) <= head+showTail {
		return strings.Repeat("*", len(value))
	}
	return value[:head] + strings.Repeat("*", len(value)-head-showTail) + value[len(value)-showTail:]
}
