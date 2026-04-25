package formatter

import (
	"fmt"
	"strings"
)

// FormatFlag is a flag.Value implementation for the Format type,
// allowing it to be used directly with the standard flag or pflag packages.
type FormatFlag struct {
	Value Format
}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	return string(f.Value)
}

// Set parses and validates the format string provided via CLI flag.
func (f *FormatFlag) Set(s string) error {
	switch Format(strings.ToLower(s)) {
	case FormatPretty, FormatJSON, FormatRaw:
		f.Value = Format(strings.ToLower(s))
		return nil
	default:
		return fmt.Errorf("unknown format %q: must be one of pretty, json, raw", s)
	}
}

// Type returns the type name used in help text.
func (f *FormatFlag) Type() string {
	return "format"
}

// DefaultFormat is the format used when none is specified.
const DefaultFormat = FormatPretty
