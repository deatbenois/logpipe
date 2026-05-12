// Package truncate provides line truncation for long log entries.
package truncate

import "unicode/utf8"

const defaultSuffix = "..."

// Truncator truncates lines that exceed a maximum byte length.
type Truncator struct {
	maxLen int
	suffix string
}

// New returns a Truncator that truncates lines longer than maxLen bytes.
// If maxLen is zero or negative, no truncation is applied.
func New(maxLen int) *Truncator {
	return &Truncator{
		maxLen: maxLen,
		suffix: defaultSuffix,
	}
}

// NewWithSuffix returns a Truncator with a custom suffix appended to truncated lines.
func NewWithSuffix(maxLen int, suffix string) *Truncator {
	return &Truncator{
		maxLen: maxLen,
		suffix: suffix,
	}
}

// Apply returns the line unchanged if it is within the limit, or a truncated
// version with the suffix appended if it exceeds the limit.
func (t *Truncator) Apply(line string) string {
	if t.maxLen <= 0 || len(line) <= t.maxLen {
		return line
	}

	// Trim to a valid UTF-8 boundary within the budget.
	budget := t.maxLen - len(t.suffix)
	if budget <= 0 {
		return t.suffix[:t.maxLen]
	}

	trimmed := trimToValidUTF8(line, budget)
	return trimmed + t.suffix
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool {
	return t.maxLen > 0
}

// MaxLen returns the configured maximum line length.
func (t *Truncator) MaxLen() int {
	return t.maxLen
}

// trimToValidUTF8 trims s to at most maxBytes bytes without splitting a rune.
func trimToValidUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	s = s[:maxBytes]
	// Walk back until we find a valid rune boundary.
	for len(s) > 0 {
		if r, _ := utf8.DecodeLastRuneInString(s); r != utf8.RuneError {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
