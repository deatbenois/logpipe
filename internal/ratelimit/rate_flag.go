package ratelimit

import (
	"fmt"
	"strconv"
)

// RateFlag is a pflag-compatible flag type for specifying lines-per-second.
type RateFlag struct {
	value float64
}

// NewRateFlag returns a RateFlag with the given default lines-per-second.
// A value of 0 means unlimited.
func NewRateFlag(defaultRate float64) *RateFlag {
	return &RateFlag{value: defaultRate}
}

// Set parses a string into a lines-per-second rate.
func (f *RateFlag) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid rate %q: must be a number", s)
	}
	if v < 0 {
		return fmt.Errorf("invalid rate %q: must be >= 0 (0 = unlimited)", s)
	}
	f.value = v
	return nil
}

// String returns the current rate as a string.
func (f *RateFlag) String() string {
	if f.value == 0 {
		return "unlimited"
	}
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

// Type returns the flag type name for pflag.
func (f *RateFlag) Type() string {
	return "rate"
}

// Value returns the parsed float64 rate.
func (f *RateFlag) Value() float64 {
	return f.value
}

// Limiter constructs a Limiter from the flag's current value.
func (f *RateFlag) Limiter() *Limiter {
	return New(f.value)
}
