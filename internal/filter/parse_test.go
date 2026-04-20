package filter_test

import (
	"testing"

	"github.com/user/logpipe/internal/filter"
)

func TestParseRules_EqExpression(t *testing.T) {
	rules, err := filter.ParseRules([]string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Field != "level" || r.Operator != "eq" || r.Value != "error" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRules_ContainsExpression(t *testing.T) {
	rules, err := filter.ParseRules([]string{"msg~timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := rules[0]
	if r.Operator != "contains" || r.Value != "timeout" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRules_ExistsExpression(t *testing.T) {
	rules, err := filter.ParseRules([]string{"trace_id?"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := rules[0]
	if r.Field != "trace_id" || r.Operator != "exists" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRules_SkipsEmpty(t *testing.T) {
	rules, err := filter.ParseRules([]string{"", "  ", "level=info"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule after skipping blanks, got %d", len(rules))
	}
}

func TestParseRules_InvalidExpression(t *testing.T) {
	_, err := filter.ParseRules([]string{"justaplainword"})
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestParseRules_EmptyFieldName(t *testing.T) {
	for _, expr := range []string{"=value", "~value", "?"} {
		_, err := filter.ParseRules([]string{expr})
		if err == nil {
			t.Errorf("expected error for expression %q with empty field", expr)
		}
	}
}
