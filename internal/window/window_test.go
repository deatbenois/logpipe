package window_test

import (
	"testing"
	"time"

	"github.com/user/logpipe/internal/window"
)

func TestCount_EmptyWindow(t *testing.T) {
	w := window.New(time.Second, 10)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_SingleEvent(t *testing.T) {
	w := window.New(time.Second, 10)
	w.Add(1)
	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestAdd_MultipleEvents(t *testing.T) {
	w := window.New(time.Second, 10)
	w.Add(3)
	w.Add(7)
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestCount_EvictsStaleEvents(t *testing.T) {
	// Use a very short window so we can expire it quickly.
	w := window.New(50*time.Millisecond, 5)
	w.Add(5)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5 before expiry, got %d", got)
	}
	time.Sleep(60 * time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_PartialEviction(t *testing.T) {
	// 100ms window, 10 buckets → 10ms per bucket.
	w := window.New(100*time.Millisecond, 10)
	w.Add(10)
	time.Sleep(55 * time.Millisecond)
	w.Add(4)
	// The first batch should be partially evicted; at least the new batch remains.
	got := w.Count()
	if got < 4 {
		t.Fatalf("expected at least 4, got %d", got)
	}
	if got > 14 {
		t.Fatalf("expected at most 14, got %d", got)
	}
}

func TestNew_SingleBucket(t *testing.T) {
	w := window.New(time.Second, 0) // 0 buckets clamped to 1
	w.Add(99)
	if got := w.Count(); got != 99 {
		t.Fatalf("expected 99, got %d", got)
	}
}
