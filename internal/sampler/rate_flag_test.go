package sampler

import (
	"testing"
)

func TestRateFlag_SetValid(t *testing.T) {
	f := NewRateFlag(1.0)
	if err := f.Set("0.25"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value() != 0.25 {
		t.Fatalf("expected 0.25, got %v", f.Value())
	}
}

func TestRateFlag_SetInvalid_NotANumber(t *testing.T) {
	f := NewRateFlag(1.0)
	if err := f.Set("abc"); err == nil {
		t.Fatal("expected error for non-numeric input")
	}
}

func TestRateFlag_SetInvalid_Zero(t *testing.T) {
	f := NewRateFlag(1.0)
	if err := f.Set("0"); err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestRateFlag_SetInvalid_AboveOne(t *testing.T) {
	f := NewRateFlag(1.0)
	if err := f.Set("1.5"); err == nil {
		t.Fatal("expected error for rate > 1")
	}
}

func TestRateFlag_SetValid_One(t *testing.T) {
	f := NewRateFlag(0.5)
	if err := f.Set("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value() != 1.0 {
		t.Fatalf("expected 1.0, got %v", f.Value())
	}
}

func TestRateFlag_String(t *testing.T) {
	f := NewRateFlag(0.1)
	if f.String() != "0.1" {
		t.Fatalf("expected \"0.1\", got %q", f.String())
	}
}

func TestRateFlag_Type(t *testing.T) {
	f := NewRateFlag(1.0)
	if f.Type() != "rate" {
		t.Fatalf("expected \"rate\", got %q", f.Type())
	}
}
