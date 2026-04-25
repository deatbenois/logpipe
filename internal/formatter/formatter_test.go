package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/formatter"
)

func TestWrite_RawFormat_NoLabel(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatRaw)
	_ = f.Write("", `{"msg":"hello"}`)
	if !strings.Contains(buf.String(), `{"msg":"hello"}`) {
		t.Errorf("expected raw line, got %q", buf.String())
	}
}

func TestWrite_RawFormat_WithLabel(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatRaw)
	_ = f.Write("app", `{"msg":"hello"}`)
	out := buf.String()
	if !strings.Contains(out, "[app]") {
		t.Errorf("expected label in output, got %q", out)
	}
}

func TestWrite_JSONFormat_InjectsSource(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatJSON)
	_ = f.Write("svc", `{"msg":"ok"}`)
	out := buf.String()
	if !strings.Contains(out, `"_source":"svc"`) {
		t.Errorf("expected _source field, got %q", out)
	}
}

func TestWrite_JSONFormat_InvalidJSON_FallsBackToRaw(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatJSON)
	_ = f.Write("svc", `not json`)
	out := buf.String()
	if !strings.Contains(out, "not json") {
		t.Errorf("expected raw fallback, got %q", out)
	}
}

func TestWrite_PrettyFormat_ContainsMsg(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatPretty)
	_ = f.Write("", `{"level":"info","msg":"started","time":"2024-01-01T12:00:00Z"}`)
	out := buf.String()
	if !strings.Contains(out, "started") {
		t.Errorf("expected msg in pretty output, got %q", out)
	}
}

func TestWrite_PrettyFormat_WithLabel(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatPretty)
	_ = f.Write("worker", `{"level":"error","msg":"boom"}`)
	out := buf.String()
	if !strings.Contains(out, "worker") {
		t.Errorf("expected label in pretty output, got %q", out)
	}
}

func TestWrite_PrettyFormat_ExtraFields(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(&buf, formatter.FormatPretty)
	_ = f.Write("", `{"level":"debug","msg":"trace","req_id":"abc123"}`)
	out := buf.String()
	if !strings.Contains(out, "req_id") {
		t.Errorf("expected extra field in pretty output, got %q", out)
	}
}
