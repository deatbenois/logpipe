// Package redact provides value-level masking for structured JSON log lines.
//
// A Redactor is constructed with a list of field names and an optional mask
// string. When applied to a log line, any JSON object field whose key matches
// a targeted field has its value replaced with the mask.
//
// Non-JSON lines pass through unmodified, making the redactor safe to use in
// pipelines that may receive mixed input.
//
// Example:
//
//	r := redact.New([]string{"password", "token"}, "")
//	masked := r.Apply(`{"user":"alice","password":"s3cr3t"}`)
//	// masked => {"user":"alice","password":"[REDACTED]"}
package redact
