package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNew_Unlimited(t *testing.T) {
	l := New(0)
	if !l.Unlimited() {
		t.Fatal("expected unlimited")
	}
}

func TestNew_Limited(t *testing.T) {
	l := New(10)
	if l.Unlimited() {
		t.Fatal("expected limited")
	}
}

func TestNew_NegativeRate_Clamped(t *testing.T) {
	l := New(-5)
	if !l.Unlimited() {
		t.Fatal("negative rate should be treated as unlimited (0)")
	}
}

func TestWait_Unlimited_NoBlock(t *testing.T) {
	l := New(0)
	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestWait_CancelledContext(t *testing.T) {
	// Very low rate so the first token is consumed and next blocks.
	l := New(0.001)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	// Drain the initial token.
	_ = l.Wait(context.Background())

	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWait_HighRate_AllowsBurst(t *testing.T) {
	l := New(1000)
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error at iteration %d: %v", i, err)
		}
	}
}

func TestRateFlag_SetValid(t *testing.T) {
	f := NewRateFlag(0)
	if err := f.Set("100"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Value() != 100 {
		t.Fatalf("expected 100, got %v", f.Value())
	}
}

func TestRateFlag_SetInvalid_Negative(t *testing.T) {
	f := NewRateFlag(0)
	if err := f.Set("-1"); err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestRateFlag_SetInvalid_NotANumber(t *testing.T) {
	f := NewRateFlag(0)
	if err := f.Set("fast"); err == nil {
		t.Fatal("expected error for non-numeric input")
	}
}

func TestRateFlag_String_Unlimited(t *testing.T) {
	f := NewRateFlag(0)
	if f.String() != "unlimited" {
		t.Fatalf("expected 'unlimited', got %q", f.String())
	}
}

func TestRateFlag_String_Numeric(t *testing.T) {
	f := NewRateFlag(50)
	if f.String() != "50" {
		t.Fatalf("expected '50', got %q", f.String())
	}
}

func TestRateFlag_Limiter(t *testing.T) {
	f := NewRateFlag(0)
	l := f.Limiter()
	if !l.Unlimited() {
		t.Fatal("expected unlimited limiter from zero rate flag")
	}
}
