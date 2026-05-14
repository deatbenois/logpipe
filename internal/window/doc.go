// Package window implements a sliding time-window counter used to track
// the rate of log lines observed over a rolling duration.
//
// The window is divided into a configurable number of equal-sized buckets.
// Each call to Add records events in the current (most-recent) bucket.
// Buckets older than the total window duration are automatically zeroed
// when Count or Add is called, providing an approximate sliding-window
// count without requiring a background goroutine.
//
// Typical usage:
//
//	w := window.New(time.Minute, 60) // 60 one-second buckets
//	w.Add(1)
//	fmt.Println(w.Count()) // events in the last minute
package window
