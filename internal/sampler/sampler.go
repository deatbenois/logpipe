package sampler

import (
	"math/rand"
	"sync/atomic"
)

// Sampler decides whether a log line should be emitted based on a rate.
// A rate of 1.0 means emit everything; 0.1 means emit ~10% of lines.
type Sampler struct {
	rate    float64
	counter atomic.Uint64
	rng     *rand.Rand
	mode    Mode
}

// Mode controls the sampling strategy.
type Mode int

const (
	// ModeRandom uses pseudo-random sampling.
	ModeRandom Mode = iota
	// ModeDeterministic uses round-robin (every N-th line).
	ModeDeterministic
)

// New creates a Sampler with the given rate and mode.
// rate must be in the range (0, 1]. Values outside this range are clamped.
func New(rate float64, mode Mode) *Sampler {
	if rate <= 0 {
		rate = 0.01
	}
	if rate > 1 {
		rate = 1
	}
	return &Sampler{
		rate: rate,
		mode: mode,
		rng:  rand.New(rand.NewSource(42)), //nolint:gosec
	}
}

// NewRandom creates a random-mode sampler.
func NewRandom(rate float64) *Sampler { return New(rate, ModeRandom) }

// NewDeterministic creates a deterministic (round-robin) sampler.
func NewDeterministic(rate float64) *Sampler { return New(rate, ModeDeterministic) }

// Sample returns true if the line should be emitted.
func (s *Sampler) Sample() bool {
	if s.rate >= 1.0 {
		return true
	}
	switch s.mode {
	case ModeDeterministic:
		return s.deterministicSample()
	default:
		return s.randomSample()
	}
}

func (s *Sampler) randomSample() bool {
	return s.rng.Float64() < s.rate //nolint:gosec
}

func (s *Sampler) deterministicSample() bool {
	n := s.counter.Add(1)
	every := uint64(1.0 / s.rate)
	if every < 1 {
		every = 1
	}
	return n%every == 0
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 { return s.rate }
