package sampler

import (
	"fmt"
	"strconv"
)

// RateFlag is a pflag-compatible flag type for sampling rate.
type RateFlag struct {
	value float64
}

// NewRateFlag creates a RateFlag with the given default.
func NewRateFlag(defaultRate float64) *RateFlag {
	return &RateFlag{value: defaultRate}
}

// Set parses and validates a rate string.
func (f *RateFlag) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid sample rate %q: %w", s, err)
	}
	if v <= 0 || v > 1 {
		return fmt.Errorf("sample rate must be in range (0, 1], got %v", v)
	}
	f.value = v
	return nil
}

// String returns the current rate as a string.
func (f *RateFlag) String() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

// Type returns the flag type name.
func (f *RateFlag) Type() string { return "rate" }

// Value returns the parsed float64 rate.
func (f *RateFlag) Value() float64 { return f.value }
