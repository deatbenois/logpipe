package truncate

import (
	"fmt"
	"strconv"
)

// LenFlag is a pflag-compatible flag that holds a maximum line length.
type LenFlag struct {
	value int
}

// NewLenFlag returns a LenFlag with the given default value.
func NewLenFlag(defaultVal int) *LenFlag {
	return &LenFlag{value: defaultVal}
}

// Set parses and validates the flag value.
func (f *LenFlag) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid length %q: must be an integer", s)
	}
	if v < 0 {
		return fmt.Errorf("invalid length %d: must be >= 0", v)
	}
	f.value = v
	return nil
}

// String returns the current value as a string.
func (f *LenFlag) String() string {
	return strconv.Itoa(f.value)
}

// Type returns the flag type name for pflag.
func (f *LenFlag) Type() string {
	return "int"
}

// Value returns the parsed integer value.
func (f *LenFlag) Value() int {
	return f.value
}

// Truncator returns a configured Truncator using the flag's value.
func (f *LenFlag) Truncator() *Truncator {
	return New(f.value)
}
