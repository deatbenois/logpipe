package truncate_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/truncate"
)

func TestLenFlag_Default(t *testing.T) {
	f := truncate.NewLenFlag(0)
	if f.Value() != 0 {
		t.Fatalf("expected default 0, got %d", f.Value())
	}
}

func TestLenFlag_SetValid(t *testing.T) {
	f := truncate.NewLenFlag(0)
	if err := f.Set("256"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value() != 256 {
		t.Fatalf("expected 256, got %d", f.Value())
	}
}

func TestLenFlag_SetZero_Allowed(t *testing.T) {
	f := truncate.NewLenFlag(100)
	if err := f.Set("0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value() != 0 {
		t.Fatalf("expected 0, got %d", f.Value())
	}
}

func TestLenFlag_SetInvalid_NotANumber(t *testing.T) {
	f := truncate.NewLenFlag(0)
	if err := f.Set("abc"); err == nil {
		t.Fatal("expected error for non-numeric input")
	}
}

func TestLenFlag_SetInvalid_Negative(t *testing.T) {
	f := truncate.NewLenFlag(0)
	if err := f.Set("-1"); err == nil {
		t.Fatal("expected error for negative value")
	}
}

func TestLenFlag_String(t *testing.T) {
	f := truncate.NewLenFlag(512)
	if f.String() != "512" {
		t.Fatalf("expected '512', got %q", f.String())
	}
}

func TestLenFlag_Type(t *testing.T) {
	f := truncate.NewLenFlag(0)
	if f.Type() != "int" {
		t.Fatalf("expected 'int', got %q", f.Type())
	}
}

func TestLenFlag_Truncator_Disabled(t *testing.T) {
	f := truncate.NewLenFlag(0)
	tr := f.Truncator()
	if tr.Enabled() {
		t.Fatal("expected truncator to be disabled")
	}
}

func TestLenFlag_Truncator_Enabled(t *testing.T) {
	f := truncate.NewLenFlag(100)
	tr := f.Truncator()
	if !tr.Enabled() {
		t.Fatal("expected truncator to be enabled")
	}
	if tr.MaxLen() != 100 {
		t.Fatalf("expected MaxLen 100, got %d", tr.MaxLen())
	}
}
