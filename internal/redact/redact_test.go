package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/redact"
)

func jsonGet(t *testing.T, line, key string) string {
	t.Helper()
	var obj map[string]string
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return obj[key]
}

func TestApply_NoFields_Unchanged(t *testing.T) {
	r := redact.New(nil, "")
	line := `{"password":"secret","user":"alice"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestApply_NonJSON_Unchanged(t *testing.T) {
	r := redact.New([]string{"password"}, "")
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestApply_RedactsTargetedField(t *testing.T) {
	r := redact.New([]string{"password"}, "")
	line := `{"password":"secret","user":"alice"}`
	got := r.Apply(line)
	if v := jsonGet(t, got, "password"); v != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", v)
	}
	if v := jsonGet(t, got, "user"); v != "alice" {
		t.Errorf("user should be unchanged, got %q", v)
	}
}

func TestApply_CustomMask(t *testing.T) {
	r := redact.New([]string{"token"}, "***")
	line := `{"token":"abc123"}`
	got := r.Apply(line)
	if v := jsonGet(t, got, "token"); v != "***" {
		t.Errorf("expected ***, got %q", v)
	}
}

func TestApply_FieldNotPresent_Unchanged(t *testing.T) {
	r := redact.New([]string{"secret"}, "")
	line := `{"user":"bob","level":"info"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestApply_MultipleFields(t *testing.T) {
	r := redact.New([]string{"password", "token"}, "")
	line := `{"password":"s3cr3t","token":"t0k3n","user":"carol"}`
	got := r.Apply(line)
	if v := jsonGet(t, got, "password"); v != "[REDACTED]" {
		t.Errorf("password: expected [REDACTED], got %q", v)
	}
	if v := jsonGet(t, got, "token"); v != "[REDACTED]" {
		t.Errorf("token: expected [REDACTED], got %q", v)
	}
	if v := jsonGet(t, got, "user"); v != "carol" {
		t.Errorf("user should be unchanged, got %q", v)
	}
}

func TestFields_ReturnsDeclared(t *testing.T) {
	r := redact.New([]string{"password", "token"}, "")
	fields := r.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
}
