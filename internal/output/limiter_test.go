package output

import (
	"bytes"
	"testing"
)

func TestLimiter_AllowsUpToMax(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 3)

	for i := 0; i < 3; i++ {
		_, err := l.Write([]byte("line"))
		if err != nil {
			t.Fatalf("unexpected error on write %d: %v", i, err)
		}
	}
}

func TestLimiter_BlocksAfterMax(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 2)

	l.Write([]byte("a"))
	l.Write([]byte("b"))
	_, err := l.Write([]byte("c"))
	if err != ErrLimitReached {
		t.Errorf("expected ErrLimitReached, got %v", err)
	}
}

func TestLimiter_NoLimit(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 0)

	for i := 0; i < 100; i++ {
		_, err := l.Write([]byte("x"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestLimiter_Remaining(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 5)

	if l.Remaining() != 5 {
		t.Errorf("expected 5 remaining, got %d", l.Remaining())
	}
	l.Write([]byte("a"))
	l.Write([]byte("b"))
	if l.Remaining() != 3 {
		t.Errorf("expected 3 remaining, got %d", l.Remaining())
	}
}

func TestLimiter_Remaining_Unlimited(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 0)
	if l.Remaining() != -1 {
		t.Errorf("expected -1 for unlimited, got %d", l.Remaining())
	}
}

func TestLimiter_WriteLine(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	l := NewLimiter(w, 1)
	l.WriteLine([]byte("hello"))
	if buf.String() != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", buf.String())
	}
}
