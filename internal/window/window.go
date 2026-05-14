// Package window provides a sliding time-window counter for tracking
// log line rates over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Window tracks event counts within a sliding time window.
type Window struct {
	mu       sync.Mutex
	buckets  []int
	size     int
	duration time.Duration
	bucket   time.Duration
	last     time.Time
}

// New creates a Window that tracks events over the given duration,
// divided into the specified number of buckets.
func New(duration time.Duration, buckets int) *Window {
	if buckets < 1 {
		buckets = 1
	}
	return &Window{
		buckets:  make([]int, buckets),
		size:     buckets,
		duration: duration,
		bucket:   duration / time.Duration(buckets),
		last:     time.Now(),
	}
}

// Add records n events at the current time.
func (w *Window) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	w.buckets[0] += n
}

// Count returns the total number of events recorded within the window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	total := 0
	for _, v := range w.buckets {
		total += v
	}
	return total
}

// advance shifts buckets forward based on elapsed time, zeroing stale ones.
func (w *Window) advance(now time.Time) {
	elapsed := now.Sub(w.last)
	if elapsed < w.bucket {
		return
	}
	shift := int(elapsed / w.bucket)
	w.last = now
	if shift >= w.size {
		for i := range w.buckets {
			w.buckets[i] = 0
		}
		return
	}
	copy(w.buckets[shift:], w.buckets[:w.size-shift])
	for i := 0; i < shift; i++ {
		w.buckets[i] = 0
	}
}
