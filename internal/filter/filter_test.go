package filter_test

import (
	"testing"

	"github.com/user/logpipe/internal/filter"
)

func TestMatch_NoRules(t *testing.T) {
	f := filter.New(nil)
	if !f.Match([]byte(`{"level":"info","msg":"hello"}`)) {
		t.Fatal("expected match with no rules")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f := filter.New([]filter.Rule{{Field: "level", Operator: "eq", Value: "info"}})
	if f.Match([]byte(`not json`)) {
		t.Fatal("expected no match for invalid JSON")
	}
}

func TestMatch_EqOperator(t *testing.T) {
	f := filter.New([]filter.Rule{{Field: "level", Operator: "eq", Value: "error"}})

	if f.Match([]byte(`{"level":"info"}`)) {
		t.Fatal("expected no match for level=info when filtering error")
	}
	if !f.Match([]byte(`{"level":"error","msg":"boom"}`)) {
		t.Fatal("expected match for level=error")
	}
}

func TestMatch_ContainsOperator(t *testing.T) {
	f := filter.New([]filter.Rule{{Field: "msg", Operator: "contains", Value: "timeout"}})

	if !f.Match([]byte(`{"msg":"connection timeout reached"}`)) {
		t.Fatal("expected match for msg containing 'timeout'")
	}
	if f.Match([]byte(`{"msg":"all good"}`)) {
		t.Fatal("expected no match for msg not containing 'timeout'")
	}
}

func TestMatch_ExistsOperator(t *testing.T) {
	f := filter.New([]filter.Rule{{Field: "trace_id", Operator: "exists"}})

	if !f.Match([]byte(`{"trace_id":"abc123","msg":"ok"}`)) {
		t.Fatal("expected match when field exists")
	}
	if f.Match([]byte(`{"msg":"no trace"}`)) {
		t.Fatal("expected no match when field absent")
	}
}

func TestMatch_MultipleRules(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Field: "level", Operator: "eq", Value: "error"},
		{Field: "service", Operator: "eq", Value: "auth"},
	})

	if !f.Match([]byte(`{"level":"error","service":"auth","msg":"fail"}`)) {
		t.Fatal("expected match when all rules pass")
	}
	if f.Match([]byte(`{"level":"error","service":"payments"}`)) {
		t.Fatal("expected no match when one rule fails")
	}
}
