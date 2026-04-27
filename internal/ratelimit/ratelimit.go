// Package ratelimit provides a token-bucket rate limiter for controlling
// log line throughput in logpipe pipelines.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate at which log lines are processed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to ratePerSec log lines per second.
// A ratePerSec of 0 means unlimited.
func New(ratePerSec float64) *Limiter {
	if ratePerSec < 0 {
		ratePerSec = 0
	}
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is acquired.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.rate == 0 {
		return nil
	}
	for {
		if l.tryAcquire() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 5):
		}
	}
}

func (l *Limiter) tryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1.0 {
		l.tokens--
		return true
	}
	return false
}

// Unlimited returns true if the limiter imposes no rate limit.
func (l *Limiter) Unlimited() bool {
	return l.rate == 0
}
