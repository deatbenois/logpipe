// Package transform provides field-level transformations applied to JSON log lines
// before they are passed to the formatter. Transformations include field renaming,
// field dropping, and value redaction.
package transform

import (
	"encoding/json"
	"fmt"
)

// Op represents a single transformation operation.
type Op struct {
	Kind  string // "rename", "drop", "redact"
	Field string
	To    string // used by rename
}

// Transformer applies a sequence of Ops to JSON log lines.
type Transformer struct {
	ops []Op
}

// New returns a Transformer that applies the given ops in order.
func New(ops []Op) *Transformer {
	return &Transformer{ops: ops}
}

// Apply parses line as a JSON object, applies all ops, and returns the
// re-serialised JSON. Non-JSON lines are returned unchanged.
func (t *Transformer) Apply(line string) string {
	if len(t.ops) == 0 {
		return line
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, op := range t.ops {
		switch op.Kind {
		case "rename":
			if v, ok := obj[op.Field]; ok {
				obj[op.To] = v
				delete(obj, op.Field)
			}
		case "drop":
			delete(obj, op.Field)
		case "redact":
			if _, ok := obj[op.Field]; ok {
				obj[op.Field] = "[REDACTED]"
			}
		}
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

// String returns a human-readable description of the op, useful for debugging.
func (o Op) String() string {
	switch o.Kind {
	case "rename":
		return fmt.Sprintf("rename(%s->%s)", o.Field, o.To)
	case "drop":
		return fmt.Sprintf("drop(%s)", o.Field)
	case "redact":
		return fmt.Sprintf("redact(%s)", o.Field)
	default:
		return fmt.Sprintf("unknown(%s)", o.Kind)
	}
}
