package formatter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/formatter"
)

func TestFormatFlag_SetValid(t *testing.T) {
	cases := []struct {
		input    string
		expected formatter.Format
	}{
		{"pretty", formatter.FormatPretty},
		{"json", formatter.FormatJSON},
		{"raw", formatter.FormatRaw},
		{"PRETTY", formatter.FormatPretty},
		{"JSON", formatter.FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			var flag formatter.FormatFlag
			if err := flag.Set(tc.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if flag.Value != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, flag.Value)
			}
		})
	}
}

func TestFormatFlag_SetInvalid(t *testing.T) {
	invalidInputs := []string{"xml", "yaml", "", "  ", "prettyjson"}
	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			var flag formatter.FormatFlag
			if err := flag.Set(input); err == nil {
				t.Errorf("expected error for unknown format %q, got nil", input)
			}
		})
	}
}

func TestFormatFlag_String(t *testing.T) {
	flag := formatter.FormatFlag{Value: formatter.FormatJSON}
	if flag.String() != "json" {
		t.Errorf("expected \"json\", got %q", flag.String())
	}
}

func TestFormatFlag_Type(t *testing.T) {
	var flag formatter.FormatFlag
	if flag.Type() != "format" {
		t.Errorf("expected \"format\", got %q", flag.Type())
	}
}

func TestDefaultFormat(t *testing.T) {
	if formatter.DefaultFormat != formatter.FormatPretty {
		t.Errorf("expected default format to be pretty, got %q", formatter.DefaultFormat)
	}
}
