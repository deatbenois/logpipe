package transform_test

import (
	"testing"

	"github.com/user/logpipe/internal/transform"
)

func TestParseOps_Rename(t *testing.T) {
	ops, err := transform.ParseOps([]string{"rename:message:msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 || ops[0].Kind != "rename" || ops[0].Field != "message" || ops[0].To != "msg" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestParseOps_Drop(t *testing.T) {
	ops, err := transform.ParseOps([]string{"drop:level"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 || ops[0].Kind != "drop" || ops[0].Field != "level" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestParseOps_Redact(t *testing.T) {
	ops, err := transform.ParseOps([]string{"redact:password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 || ops[0].Kind != "redact" || ops[0].Field != "password" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestParseOps_SkipsEmpty(t *testing.T) {
	ops, err := transform.ParseOps([]string{"", "  ", "drop:ts"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("expected 1 op, got %d", len(ops))
	}
}

func TestParseOps_InvalidKind(t *testing.T) {
	_, err := transform.ParseOps([]string{"uppercase:field"})
	if err == nil {
		t.Error("expected error for unknown kind")
	}
}

func TestParseOps_InvalidRename_MissingTo(t *testing.T) {
	_, err := transform.ParseOps([]string{"rename:old"})
	if err == nil {
		t.Error("expected error for incomplete rename")
	}
}

func TestParseOps_InvalidDrop_MissingField(t *testing.T) {
	_, err := transform.ParseOps([]string{"drop:"})
	if err == nil {
		t.Error("expected error for empty drop field")
	}
}

func TestOpString_Rename(t *testing.T) {
	op := transform.Op{Kind: "rename", Field: "a", To: "b"}
	if s := op.String(); s != "rename(a->b)" {
		t.Errorf("unexpected string: %q", s)
	}
}

func TestOpString_Drop(t *testing.T) {
	op := transform.Op{Kind: "drop", Field: "ts"}
	if s := op.String(); s != "drop(ts)" {
		t.Errorf("unexpected string: %q", s)
	}
}
