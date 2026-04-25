package output

import (
	"errors"
	"io"
	"sync/atomic"
)

// ErrLimitReached is returned when the line limit has been exceeded.
var ErrLimitReached = errors.New("output: line limit reached")

// Limiter wraps a Writer and stops writing after a maximum number of lines.
type Limiter struct {
	w     *Writer
	max   int64
	count atomic.Int64
}

// NewLimiter creates a Limiter that allows at most max lines through.
// If max <= 0, no limit is applied.
func NewLimiter(w *Writer, max int64) *Limiter {
	return &Limiter{w: w, max: max}
}

// Write writes p if the line limit has not been reached.
// Returns ErrLimitReached once the limit is exceeded.
func (l *Limiter) Write(p []byte) (int, error) {
	if l.max > 0 {
		n := l.count.Add(1)
		if n > l.max {
			return 0, ErrLimitReached
		}
	}
	return l.w.Write(p)
}

// WriteLine writes p followed by a newline if the line limit has not been reached.
func (l *Limiter) WriteLine(p []byte) (int, error) {
	return l.Write(append(p, '\n'))
}

// Remaining returns the number of lines still allowed, or -1 if unlimited.
func (l *Limiter) Remaining() int64 {
	if l.max <= 0 {
		return -1
	}
	r := l.max - l.count.Load()
	if r < 0 {
		return 0
	}
	return r
}

// Underlying returns the wrapped Writer.
func (l *Limiter) Underlying() io.Writer {
	return l.w
}
