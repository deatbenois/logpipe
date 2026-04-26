// Package sampler provides log-line sampling for logpipe.
//
// Two sampling strategies are supported:
//
//   - Random: each line is independently emitted with probability equal to
//     the configured rate. Suitable for high-volume streams where approximate
//     throughput reduction is acceptable.
//
//   - Deterministic: every N-th line is emitted (where N = 1/rate), giving
//     perfectly even distribution. Useful when reproducible output is needed.
//
// Usage:
//
//	s := sampler.NewRandom(0.1)   // emit ~10% of lines
//	if s.Sample() {
//	    // forward the line
//	}
//
// The RateFlag type integrates with pflag/cobra for CLI flag parsing.
package sampler
