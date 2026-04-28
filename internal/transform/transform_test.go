package transform_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logpipe/internal/transform"
)

func jsonField(t *testing.T, line, field string) any {
	t.Helper()
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("not valid JSON: %v", err)
	}
	return obj[field]
}

func TestApply_NoOps(t *testing.T) {
	tr := transform.New(nil)
	const line = `{"level":"info","msg":"hello"}`
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_PassThrough(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"drop:level"}))
	const line = "not json at all"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_Drop(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"drop:secret"}))
	out := tr.Apply(`{"msg":"hi","secret":"s3cr3t"}`)
	if v := jsonField(t, out, "secret"); v != nil {
		t.Errorf("expected secret to be dropped, got %v", v)
	}
	if v := jsonField(t, out, "msg"); v != "hi" {
		t.Errorf("expected msg=hi, got %v", v)
	}
}

func TestApply_Rename(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"rename:message:msg"}))
	out := tr.Apply(`{"message":"hello"}`)
	if v := jsonField(t, out, "msg"); v != "hello" {
		t.Errorf("expected msg=hello, got %v", v)
	}
	if v := jsonField(t, out, "message"); v != nil {
		t.Errorf("expected message to be removed, got %v", v)
	}
}

func TestApply_Redact(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"redact:token"}))
	out := tr.Apply(`{"token":"abc123","user":"alice"}`)
	if v := jsonField(t, out, "token"); v != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", v)
	}
}

func TestApply_DropMissingField_NoError(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"drop:nonexistent"}))
	const line = `{"msg":"ok"}`
	if got := tr.Apply(line); jsonField(t, got, "msg") != "ok" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestApply_MultipleOps(t *testing.T) {
	tr := transform.New(transform.MustParseOps([]string{"drop:ts", "redact:pass", "rename:message:msg"}))
	out := tr.Apply(`{"ts":"2024","pass":"hunter2","message":"boot"}`)
	if jsonField(t, out, "ts") != nil {
		t.Error("ts should be dropped")
	}
	if jsonField(t, out, "pass") != "[REDACTED]" {
		t.Error("pass should be redacted")
	}
	if jsonField(t, out, "msg") != "boot" {
		t.Error("message should be renamed to msg")
	}
}
