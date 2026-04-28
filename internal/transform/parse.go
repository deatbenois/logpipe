package transform

import (
	"fmt"
	"strings"
)

// ParseOps parses a slice of raw expression strings into Ops.
// Supported syntax:
//
//	rename:old:new   — rename field "old" to "new"
//	drop:field       — remove field from output
//	redact:field     — replace field value with "[REDACTED]"
//
// Empty strings are silently skipped. An error is returned for any
// unrecognised expression.
func ParseOps(exprs []string) ([]Op, error) {
	ops := make([]Op, 0, len(exprs))
	for _, raw := range exprs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		op, err := parseExpr(raw)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return ops, nil
}

// MustParseOps is like ParseOps but panics on error. Intended for tests.
func MustParseOps(exprs []string) []Op {
	ops, err := ParseOps(exprs)
	if err != nil {
		panic(err)
	}
	return ops
}

func parseExpr(expr string) (Op, error) {
	parts := strings.SplitN(expr, ":", 3)
	switch parts[0] {
	case "rename":
		if len(parts) != 3 || parts[1] == "" || parts[2] == "" {
			return Op{}, fmt.Errorf("transform: invalid rename expression %q: expected rename:old:new", expr)
		}
		return Op{Kind: "rename", Field: parts[1], To: parts[2]}, nil
	case "drop":
		if len(parts) < 2 || parts[1] == "" {
			return Op{}, fmt.Errorf("transform: invalid drop expression %q: expected drop:field", expr)
		}
		return Op{Kind: "drop", Field: parts[1]}, nil
	case "redact":
		if len(parts) < 2 || parts[1] == "" {
			return Op{}, fmt.Errorf("transform: invalid redact expression %q: expected redact:field", expr)
		}
		return Op{Kind: "redact", Field: parts[1]}, nil
	default:
		return Op{}, fmt.Errorf("transform: unknown operation %q in expression %q", parts[0], expr)
	}
}
