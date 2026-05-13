// Package redact provides field-level value masking for structured JSON log lines.
package redact

import (
	"encoding/json"
	"strings"
)

// Redactor masks specific field values in JSON log lines.
type Redactor struct {
	fields map[string]struct{}
	mask   string
}

// New returns a Redactor that replaces values of the given fields with mask.
// If mask is empty, "[REDACTED]" is used.
func New(fields []string, mask string) *Redactor {
	if mask == "" {
		mask = "[REDACTED]"
	}
	f := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		if field != "" {
			f[strings.TrimSpace(field)] = struct{}{}
		}
	}
	return &Redactor{fields: f, mask: mask}
}

// Apply returns the line with targeted field values replaced by the mask.
// Non-JSON lines are returned unchanged. Lines with no matching fields are
// returned unchanged.
func (r *Redactor) Apply(line string) string {
	if len(r.fields) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	modified := false
	for field := range r.fields {
		if _, ok := obj[field]; ok {
			obj[field] = json.RawMessage(`"` + r.mask + `"`)
			modified = true
		}
	}

	if !modified {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// Fields returns the set of field names this Redactor targets.
func (r *Redactor) Fields() []string {
	result := make([]string, 0, len(r.fields))
	for f := range r.fields {
		result = append(result, f)
	}
	return result
}
