package filter

import (
	"encoding/json"
	"strings"
)

// Rule defines a single filter condition applied to a JSON log entry.
type Rule struct {
	Field    string
	Operator string // "eq", "contains", "exists"
	Value    string
}

// Filter holds a set of rules and applies them to log lines.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from a slice of Rule definitions.
func New(rules []Rule) *Filter {
	return &Filter{Rules: rules}
}

// Match returns true if the JSON log line satisfies all filter rules.
func (f *Filter) Match(line []byte) bool {
	if len(f.Rules) == 0 {
		return true
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return false
	}

	for _, rule := range f.Rules {
		if !applyRule(entry, rule) {
			return false
		}
	}
	return true
}

func applyRule(entry map[string]interface{}, rule Rule) bool {
	val, ok := entry[rule.Field]

	switch rule.Operator {
	case "exists":
		return ok
	case "eq":
		if !ok {
			return false
		}
		return toString(val) == rule.Value
	case "contains":
		if !ok {
			return false
		}
		return strings.Contains(toString(val), rule.Value)
	default:
		return false
	}
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strings.TrimRight(strings.TrimRight(strings.Replace(strings.Replace(string([]byte{}), ".", ".", 1), "", "", 1), "0"), ".")
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
