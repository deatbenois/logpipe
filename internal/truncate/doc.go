// Package truncate provides configurable line truncation for logpipe.
//
// Lines that exceed a maximum byte length are shortened and a suffix (default
// "...") is appended so that downstream consumers always receive output within
// a predictable size budget. Truncation is aware of UTF-8 rune boundaries and
// will never split a multi-byte character.
//
// Usage:
//
//	tr := truncate.New(512)          // truncate at 512 bytes
//	line = tr.Apply(line)            // apply to each log line
//
// A LenFlag is provided for integration with pflag-based CLI flags:
//
//	f := truncate.NewLenFlag(0)      // 0 = disabled by default
//	flagSet.Var(f, "max-line-len", "truncate lines longer than N bytes (0 = off)")
//	tr := f.Truncator()
package truncate
