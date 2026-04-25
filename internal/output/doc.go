// Package output provides thread-safe writers and output controls for logpipe.
//
// The Writer type wraps any io.Writer and ensures concurrent writes do not
// interleave, making it safe to use from multiple goroutines (e.g. when
// fanning in log lines from multiple sources).
//
// The Limiter type wraps a Writer and enforces a maximum number of output
// lines, returning ErrLimitReached once the cap is hit. This is useful for
// the --max-lines / -n flag exposed by the CLI.
//
// Typical usage:
//
//	w := output.NewStdout()
//	lim := output.NewLimiter(w, 100)
//	// pass lim to formatter.New(...)
package output
