package filter

import (
	"fmt"
	"strings"
)

// ParseRules parses a slice of raw filter expressions into Rule structs.
// Supported formats:
//   - "field=value"      → eq operator
//   - "field~value"      → contains operator
//   - "field?"           → exists operator
func ParseRules(exprs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(exprs))

	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}

		rule, err := parseExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid filter expression %q: %w", expr, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func parseExpr(expr string) (Rule, error) {
	if strings.HasSuffix(expr, "?") {
		field := strings.TrimSuffix(expr, "?")
		if field == "" {
			return Rule{}, fmt.Errorf("field name cannot be empty")
		}
		return Rule{Field: field, Operator: "exists"}, nil
	}

	if idx := strings.Index(expr, "~"); idx != -1 {
		field := expr[:idx]
		value := expr[idx+1:]
		if field == "" {
			return Rule{}, fmt.Errorf("field name cannot be empty")
		}
		return Rule{Field: field, Operator: "contains", Value: value}, nil
	}

	if idx := strings.Index(expr, "="); idx != -1 {
		field := expr[:idx]
		value := expr[idx+1:]
		if field == "" {
			return Rule{}, fmt.Errorf("field name cannot be empty")
		}
		return Rule{Field: field, Operator: "eq", Value: value}, nil
	}

	return Rule{}, fmt.Errorf("unrecognised operator; use =, ~, or ? suffix")
}

// MustParseRules is like ParseRules but panics on error.
// It is intended for use in tests or program initialisation where
// filter expressions are known to be valid at compile time.
func MustParseRules(exprs []string) []Rule {
	rules, err := ParseRules(exprs)
	if err != nil {
		panic(err)
	}
	return rules
}
