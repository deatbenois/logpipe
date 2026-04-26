package sampler

import (
	"testing"
)

func TestSample_RateOne_AlwaysTrue(t *testing.T) {
	s := NewRandom(1.0)
	for i := 0; i < 100; i++ {
		if !s.Sample() {
			t.Fatal("expected true for rate=1.0")
		}
	}
}

func TestSample_RateZero_Clamped(t *testing.T) {
	s := NewRandom(0)
	if s.Rate() <= 0 {
		t.Fatal("rate should be clamped above 0")
	}
}

func TestSample_RateAboveOne_Clamped(t *testing.T) {
	s := NewRandom(5.0)
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate clamped to 1.0, got %v", s.Rate())
	}
}

func TestSample_Deterministic_EverySecond(t *testing.T) {
	s := NewDeterministic(0.5)
	results := make([]bool, 10)
	for i := range results {
		results[i] = s.Sample()
	}
	// With rate=0.5, every=2, so indices 1,3,5,7,9 (n%2==0 when n=2,4,6,8,10)
	var trueCount int
	for _, r := range results {
		if r {
			trueCount++
		}
	}
	if trueCount != 5 {
		t.Fatalf("expected 5 true results, got %d", trueCount)
	}
}

func TestSample_Deterministic_EveryTenth(t *testing.T) {
	s := NewDeterministic(0.1)
	var trueCount int
	for i := 0; i < 100; i++ {
		if s.Sample() {
			trueCount++
		}
	}
	if trueCount != 10 {
		t.Fatalf("expected 10 true results for rate=0.1, got %d", trueCount)
	}
}

func TestSample_Random_ApproximateRate(t *testing.T) {
	s := NewRandom(0.5)
	var trueCount int
	const n = 10000
	for i := 0; i < n; i++ {
		if s.Sample() {
			trueCount++
		}
	}
	ratio := float64(trueCount) / n
	if ratio < 0.45 || ratio > 0.55 {
		t.Fatalf("random sampling ratio %v out of expected range [0.45, 0.55]", ratio)
	}
}
