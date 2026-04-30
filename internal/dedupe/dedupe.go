// Package dedupe provides a filter that suppresses consecutive duplicate log lines.
package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// Deduplicator tracks recently seen log lines and reports whether a given
// line is a duplicate of the previous one seen on the same source.
type Deduplicator struct {
	mu      sync.Mutex
	last    map[string]string // source -> last hash
	window  int              // how many unique hashes to remember per source (0 = only last)
	history map[string][]string
}

// New returns a Deduplicator that suppresses consecutive identical lines.
// If window > 1 it remembers the last N distinct hashes per source,
// suppressing any line whose hash appears in that window.
func New(window int) *Deduplicator {
	if window < 1 {
		window = 1
	}
	return &Deduplicator{
		last:    make(map[string]string),
		window:  window,
		history: make(map[string][]string),
	}
}

// IsDuplicate returns true if line has been seen within the configured window
// for the given source, and records the line for future checks.
func (d *Deduplicator) IsDuplicate(source, line string) bool {
	h := hash(line)

	d.mu.Lock()
	defer d.mu.Unlock()

	hist := d.history[source]
	for _, prev := range hist {
		if prev == h {
			return true
		}
	}

	// Append and trim to window size.
	hist = append(hist, h)
	if len(hist) > d.window {
		hist = hist[len(hist)-d.window:]
	}
	d.history[source] = hist
	d.last[source] = h
	return false
}

// Reset clears all tracked state for every source.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = make(map[string]string)
	d.history = make(map[string][]string)
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
